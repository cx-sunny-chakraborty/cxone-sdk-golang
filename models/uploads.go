package models

// UploadModel is the response of POST /api/uploads — a presigned URL pointing
// to the object-store bucket where the source archive should be PUT.
//
// The URL embeds short-lived credentials and is one-shot; the caller is
// expected to PUT the file body to the URL within minutes of receiving it
// and then reference the same URL when starting the scan.
type UploadModel struct {
	URL string `json:"url"`
}
