package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func CallAPI(url string, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add headers to the request
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func DecodeJSONResponse(body []byte, v interface{}) error {
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&v); err != nil {
		return err
	}
	return nil
}
