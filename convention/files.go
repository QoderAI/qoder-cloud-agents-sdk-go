package convention

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func DownloadTo(ctx context.Context, client *http.Client, url string, destination io.Writer) error {
	if client == nil {
		client = http.DefaultClient
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		body, readErr := io.ReadAll(response.Body)
		if readErr != nil {
			return readErr
		}
		apiError := &Error{StatusCode: response.StatusCode, Request: request, Response: response, RequestID: response.Header.Get("X-Request-ID")}
		_ = json.Unmarshal(body, apiError)
		return apiError
	}
	_, err = io.Copy(destination, response.Body)
	return err
}

func DownloadBytes(ctx context.Context, client *http.Client, url string) ([]byte, error) {
	var output bytes.Buffer
	if err := DownloadTo(ctx, client, url, &output); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func DownloadAPITo(
	ctx context.Context,
	client *http.Client,
	baseURL string,
	path string,
	destination io.Writer,
) error {
	return DownloadTo(
		ctx,
		client,
		strings.TrimRight(baseURL, "/")+"/"+strings.TrimLeft(path, "/"),
		destination,
	)
}
