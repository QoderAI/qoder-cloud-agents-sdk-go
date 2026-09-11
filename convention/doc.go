// Package convention implements Qoder authentication, transport, request options,
// errors and download conventions. Its param, pagination, respjson and ssestream
// subpackages provide the shared data structures for Managed and Forward modes.
package convention

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
	"io"
)

type Error = apierror.Error
type Metadata = map[string]string

// UploadFile supplies a filename (including relative paths for skill file trees)
// and media type to the multipart serializer. Reader is consumed by the call.
type UploadFile struct {
	io.Reader
	Name      string
	MediaType string
}

func (f UploadFile) Filename() string { return f.Name }
func (f UploadFile) ContentType() string {
	if f.MediaType == "" {
		return "application/octet-stream"
	}
	return f.MediaType
}
