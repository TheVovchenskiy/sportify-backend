package api

import (
	"net/http"
	"net/http/httptest"
)

func WriteHeaderToWriter(headers http.Header, writer http.ResponseWriter) {
	for key, values := range headers {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}
}

func WriteFromDummyWriterToWriter(dummyWriter *httptest.ResponseRecorder, writer http.ResponseWriter) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(dummyWriter.Code)
	WriteHeaderToWriter(dummyWriter.Header(), writer)
	writer.Write(dummyWriter.Body.Bytes())
	return
}
