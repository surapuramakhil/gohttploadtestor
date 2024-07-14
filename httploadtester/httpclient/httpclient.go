package httpclient

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func NewFileUploadRequest(uri, method, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, uri, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}

func SendRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	

	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("received status code %d", resp.StatusCode)
	}

	return resp, nil
}

