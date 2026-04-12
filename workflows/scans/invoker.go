package scans

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/checkmarx-open-labs/cxone-sdk-golang/cxerrors"
	endpoints "github.com/checkmarx-open-labs/cxone-sdk-golang/internal/endpoints/scans"
	"github.com/checkmarx-open-labs/cxone-sdk-golang/internal/endpoints/uploads"
	"github.com/checkmarx-open-labs/cxone-sdk-golang/models"
	"github.com/checkmarx-open-labs/cxone-sdk-golang/workflows/projects"
)

// Invoker is the fluent builder for starting a new Checkmarx One scan.
//
// Construct one with [NewInvoker]; chain configuration setters; call
// [Invoker.Start] to send the scan request.
//
// CLAUDE.md §9.2 mandates fluent builders with up-front validation for
// operations with many optional inputs. The Invoker validates engine names,
// branch presence, and source-mode consistency before issuing any HTTP
// request — invalid combinations surface as *cxerrors.ConfigurationError
// from Start, not from a remote 400.
//
// Two source modes are supported:
//
//   - Git mode: ForBranch sets a branch on the project's primary repo URL.
//     The platform checks out the source via its own SCM integration.
//   - Upload mode: WithUploadReader uploads a zip via the two-step
//     presigned-URL dance and references the URL in the scan request.
//
// Use exactly one source mode per Invoker. Mixing them surfaces a
// ConfigurationError on Start.
type Invoker struct {
	backend *Backend
	repo    *projects.ProjectRepoConfig

	branch       string
	tag          string
	commit       string
	engines      []string
	tags         map[string]string
	sastConfig   *models.SastConfig
	scaConfig    *models.ScaConfig
	kicsConfig   *models.KicsConfig
	containerCfg *models.ContainerConfig
	apisecCfg    *models.APISecConfig
	scsCfg       *models.SCSConfig

	// upload-mode state
	uploadBody io.Reader
	uploadSize int64

	// httpClient is used for the upload-mode PUT (which bypasses the funnel).
	httpClient *http.Client

	// validation errors collected during the chain — surfaced from Start.
	configErr error
}

// NewInvoker constructs a fluent scan builder bound to the supplied
// project. The httpClient is used for the upload-mode presigned-URL PUT
// path; pass the same client the cxone client is using for shared connection
// pools.
func NewInvoker(backend *Backend, repo *projects.ProjectRepoConfig, httpClient *http.Client) *Invoker {
	if backend == nil || backend.Executor == nil {
		return &Invoker{configErr: &cxerrors.ConfigurationError{Field: "backend", Reason: "is required"}}
	}
	if repo == nil {
		return &Invoker{configErr: &cxerrors.ConfigurationError{Field: "repo", Reason: "is required"}}
	}
	return &Invoker{
		backend:    backend,
		repo:       repo,
		httpClient: httpClient,
		tags:       map[string]string{},
	}
}

// ForBranch selects the git branch to scan in git mode. Mutually exclusive
// with [Invoker.WithUploadReader].
func (i *Invoker) ForBranch(branch string) *Invoker {
	i.branch = branch
	return i
}

// ForCommit narrows the git-mode scan to a specific commit. Optional.
func (i *Invoker) ForCommit(commit string) *Invoker {
	i.commit = commit
	return i
}

// ForTag narrows the git-mode scan to a specific git tag. Optional.
func (i *Invoker) ForTag(tag string) *Invoker {
	i.tag = tag
	return i
}

// WithEngines selects which Checkmarx One engines to run.
//
// Valid names: "sast", "sca", "kics", "containers", "apisec", "scs".
// Unknown engine names cause Start to return a *cxerrors.ConfigurationError.
func (i *Invoker) WithEngines(engines ...string) *Invoker {
	i.engines = append([]string(nil), engines...)
	return i
}

// WithTags adds key/value tags to the scan. Multiple calls accumulate.
func (i *Invoker) WithTags(tags map[string]string) *Invoker {
	for k, v := range tags {
		i.tags[k] = v
	}
	return i
}

// WithSastConfig overrides the SAST engine configuration for this scan
// (preset name, fast-scan mode, light queries, etc.). The platform's
// project-level configuration is used for any field left unset.
func (i *Invoker) WithSastConfig(c models.SastConfig) *Invoker {
	i.sastConfig = &c
	return i
}

// WithScaConfig overrides the SCA engine configuration.
func (i *Invoker) WithScaConfig(c models.ScaConfig) *Invoker {
	i.scaConfig = &c
	return i
}

// WithKicsConfig overrides the KICS (IaC) engine configuration.
func (i *Invoker) WithKicsConfig(c models.KicsConfig) *Invoker {
	i.kicsConfig = &c
	return i
}

// WithContainerConfig overrides the container-scanning engine configuration.
func (i *Invoker) WithContainerConfig(c models.ContainerConfig) *Invoker {
	i.containerCfg = &c
	return i
}

// WithAPISecConfig overrides the API security engine configuration.
func (i *Invoker) WithAPISecConfig(c models.APISecConfig) *Invoker {
	i.apisecCfg = &c
	return i
}

// WithSCSConfig overrides the source-code security engine configuration.
func (i *Invoker) WithSCSConfig(c models.SCSConfig) *Invoker {
	i.scsCfg = &c
	return i
}

// WithUploadReader switches the invoker to upload mode. body is streamed
// to the platform's object store via the two-step presigned-URL dance;
// size is the byte length of the body and MUST be accurate (object stores
// require Content-Length).
//
// Mutually exclusive with [Invoker.ForBranch] / [Invoker.ForCommit] /
// [Invoker.ForTag].
func (i *Invoker) WithUploadReader(body io.Reader, size int64) *Invoker {
	i.uploadBody = body
	i.uploadSize = size
	return i
}

// Start validates the builder, optionally uploads the source archive, and
// issues POST /api/scans. Returns a [*ScanInspector] populated with the
// platform's response so the caller can poll for completion or pass it to
// [WaitUntilComplete] directly.
//
// Validation errors surface as *cxerrors.ConfigurationError. Upload errors
// (when in upload mode) and HTTP errors are returned verbatim from the
// underlying endpoint helpers.
func (i *Invoker) Start(ctx context.Context) (*ScanInspector, error) {
	if i.configErr != nil {
		return nil, i.configErr
	}
	if err := i.validate(); err != nil {
		return nil, err
	}

	// Build handler payload depending on source mode.
	var (
		handlerJSON json.RawMessage
		err         error
	)
	if i.uploadBody != nil {
		handlerJSON, err = i.buildUploadHandler(ctx)
	} else {
		handlerJSON, err = i.buildGitHandler()
	}
	if err != nil {
		return nil, err
	}

	scanReq := &models.Scan{
		Type:    i.sourceType(),
		Handler: handlerJSON,
		Project: models.ScanProject{ID: i.repo.ID(), Tags: i.tags},
		Tags:    i.tags,
		Config:  i.buildEngineConfigs(),
	}

	resp, err := endpoints.Create(ctx, i.backend.Executor, i.backend.BaseURL, scanReq)
	if err != nil {
		return nil, err
	}
	return FromScanModel(i.backend, resp)
}

// validate runs up-front sanity checks. CLAUDE.md §9.2 — invalid combinations
// must surface as *cxerrors.ConfigurationError BEFORE any HTTP call.
func (i *Invoker) validate() error {
	if len(i.engines) == 0 {
		return &cxerrors.ConfigurationError{Field: "Engines", Reason: "at least one engine must be selected"}
	}
	if !validEngines(i.engines) {
		return &cxerrors.ConfigurationError{Field: "Engines", Reason: "unknown engine name; valid: sast, sca, kics, containers, apisec, scs"}
	}
	gitMode := i.branch != "" || i.commit != "" || i.tag != ""
	uploadMode := i.uploadBody != nil
	if gitMode && uploadMode {
		return &cxerrors.ConfigurationError{Field: "Source", Reason: "git mode (ForBranch/ForCommit/ForTag) and upload mode (WithUploadReader) are mutually exclusive"}
	}
	if !gitMode && !uploadMode {
		// Fall back to the project's primary branch when nothing is set.
		if i.repo.MainBranch() == "" {
			return &cxerrors.ConfigurationError{Field: "Branch", Reason: "no branch supplied and project has no MainBranch set"}
		}
		i.branch = i.repo.MainBranch()
	}
	if uploadMode && i.uploadSize <= 0 {
		return &cxerrors.ConfigurationError{Field: "UploadSize", Reason: "must be > 0 (object stores require Content-Length)"}
	}
	if uploadMode && i.httpClient == nil {
		return &cxerrors.ConfigurationError{Field: "httpClient", Reason: "is required for upload mode"}
	}
	return nil
}

func (i *Invoker) sourceType() string {
	if i.uploadBody != nil {
		return "upload"
	}
	return "git"
}

func (i *Invoker) buildGitHandler() (json.RawMessage, error) {
	repoURL := i.repo.RepoURL()
	if repoURL == "" {
		return nil, &cxerrors.ConfigurationError{Field: "repoUrl", Reason: "project has no RepoURL set"}
	}
	h := models.GitProjectHandler{
		RepoURL: repoURL,
		Branch:  i.branch,
		Commit:  i.commit,
		Tag:     i.tag,
	}
	return json.Marshal(h)
}

func (i *Invoker) buildUploadHandler(ctx context.Context) (json.RawMessage, error) {
	// Step 1: get a presigned URL.
	uploadURL, err := uploads.GetPresignedURL(ctx, i.backend.Executor, i.backend.BaseURL)
	if err != nil {
		return nil, err
	}
	// Step 2: PUT the body. uploads.PutFile bypasses the funnel because the
	// presigned URL points at the object store, not the API host.
	if err := uploads.PutFile(ctx, i.httpClient, uploadURL, i.uploadBody, i.uploadSize); err != nil {
		return nil, err
	}
	h := models.ScanHandler{
		RepoURL:   i.repo.RepoURL(),
		UploadURL: uploadURL,
	}
	return json.Marshal(h)
}

// buildEngineConfigs assembles the [models.Config] entries for the engines
// the caller selected. Each config is wrapped as { "type": "<engine>",
// "value": { ... } }.
//
// Engines whose user-supplied config struct is nil get an empty
// {"value": {}} entry — this lets the platform fall back to its
// project-level defaults rather than disabling the engine entirely.
func (i *Invoker) buildEngineConfigs() []models.Config {
	out := make([]models.Config, 0, len(i.engines))
	for _, eng := range i.engines {
		var src any
		switch eng {
		case "sast":
			src = i.sastConfig
		case "sca":
			src = i.scaConfig
		case "kics":
			src = i.kicsConfig
		case "containers":
			src = i.containerCfg
		case "apisec":
			src = i.apisecCfg
		case "scs":
			src = i.scsCfg
		}
		out = append(out, models.Config{
			Type:  eng,
			Value: structToMap(src),
		})
	}
	return out
}

// structToMap marshals a typed config struct into a map for the wire's
// untyped {string: interface{}} value. nil sources collapse to nil so the
// platform applies project-level defaults.
func structToMap(src any) map[string]any {
	if src == nil {
		return nil
	}
	// Defensive: if the typed pointer is the typed-nil case (e.g.
	// `(*SastConfig)(nil)`), reflect would still see a non-nil interface.
	// Marshal/Unmarshal handles both shapes uniformly.
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(src); err != nil || buf.Len() <= len("null\n") {
		return nil
	}
	var out map[string]any
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		return nil
	}
	return out
}

// validEngines returns true if every name is one of the platform's known
// engine identifiers.
func validEngines(names []string) bool {
	for _, n := range names {
		switch n {
		case "sast", "sca", "kics", "containers", "apisec", "scs":
			// ok
		default:
			return false
		}
	}
	return true
}
