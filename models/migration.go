package models

import "time"

// ImportStatus values are Checkmarx One's terminal states for a data import.
const (
	ImportStatusQueued    = "queued"
	ImportStatusRunning   = "running"
	ImportStatusCompleted = "completed"
	ImportStatusPartial   = "partial"
	ImportStatusFailed    = "failed"
	ImportStatusBlank     = "blank"
)

// DataImport is one row returned by GET /api/imports and GET /api/imports/{id}.
type DataImport struct {
	MigrationID string             `json:"migrationId,omitempty"`
	Status      string             `json:"status,omitempty"`
	CreatedAt   time.Time          `json:"createdAt,omitempty"`
	Logs        []DataImportStatus `json:"logs,omitempty"`
}

// DataImportStatus is one log entry attached to a [DataImport].
type DataImportStatus struct {
	Level   string `json:"level,omitempty"`
	Message string `json:"msg,omitempty"`
	Error   string `json:"error,omitempty"`
	Worker  string `json:"worker,omitempty"`
	RawLog  string `json:"raw_log,omitempty"`
}

// ImportStartRequest is the body of POST /api/imports. fileName is the
// uploaded data archive filename; projectsMappingFileName is optional and
// points at a previously-uploaded project-id mapping file; encryptionKey is
// the key that was used to encrypt the archive during export.
type ImportStartRequest struct {
	FileName                string `json:"fileName"`
	ProjectsMappingFileName string `json:"projectsMappingFileName,omitempty"`
	EncryptionKey           string `json:"encryptionKey"`
}

// ImportStartResponse is the envelope returned when an import begins —
// migrationId is the handle used by every subsequent [GetImport] call.
type ImportStartResponse struct {
	MigrationID string `json:"migrationId"`
}
