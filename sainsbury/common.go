package sainsbury

import (
	"bytes"
	"io"
	"net/http"
)

func getRawHTML(url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte{}))

	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
