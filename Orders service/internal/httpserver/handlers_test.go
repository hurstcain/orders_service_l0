package httpserver

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestStartServer(t *testing.T) {
	respStartPage, err := http.Get("http://127.0.0.1:8000")

	assert.NoError(t, err)
	assert.Equal(t, respStartPage.StatusCode, 200)

	respInvalidPage, err := http.Get("http://127.0.0.1:8000/p/s/a/q")

	assert.NoError(t, err)
	assert.Equal(t, respInvalidPage.StatusCode, 404)

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	writer.WriteField("order_uid", "test")
	writer.Close()
	req, _ := http.NewRequest("POST", "http://127.0.0.1:8000/order", payload)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	respOrderInfo, err := client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, respOrderInfo.StatusCode, 200)
}
