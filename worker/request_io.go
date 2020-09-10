package worker

import (
	"io"
	"net/http"

	"nano-run/services/meta"
)

func openResponse(writer io.Writer) *responseStream {
	return &responseStream{stream: writer, meta: meta.AttemptHeader{Headers: make(http.Header)}}
}

type responseStream struct {
	meta       meta.AttemptHeader
	statusSent bool
	stream     io.Writer
}

func (mo *responseStream) Header() http.Header {
	return mo.meta.Headers
}

func (mo *responseStream) Write(bytes []byte) (int, error) {
	if !mo.statusSent {
		mo.WriteHeader(http.StatusOK)
	}
	return mo.stream.Write(bytes)
}

func (mo *responseStream) WriteHeader(statusCode int) {
	if mo.statusSent {
		return
	}
	mo.statusSent = true
	mo.meta.Code = statusCode
}
