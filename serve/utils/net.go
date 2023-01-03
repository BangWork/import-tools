package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

func PostJSON(url string, body interface{}) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func PostJSONWithHeader(url string, body interface{}, header map[string]string) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	request.Header.Set("content-type", "application/json")
	for k, v := range header {
		request.Header.Set(k, v)
	}
	resp, err := client.Do(request)
	return resp, err
}

func PostByteWithHeader(url string, body []byte, header map[string]string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		request.Header.Set(k, v)
	}
	resp, err := client.Do(request)
	return resp, err
}

func GetWithHeader(url string, header map[string]string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		request.Header.Set(k, v)
	}
	resp, err := client.Do(request)
	return resp, err
}

func PostFileUpload(url string, token string, file *os.File, realFileName string) (*http.Response, error) {
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("token", token)
	part, err := writer.CreateFormFile("file", realFileName)
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	request, err := http.NewRequest("POST", url, payload)

	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(request)
	return resp, err
}
