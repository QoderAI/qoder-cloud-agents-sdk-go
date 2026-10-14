// Drive Beta: upload, list, download and delete with a new isolated Identity.
package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/forwardutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	config, err := live.Load("forward")
	if err != nil {
		fmt.Fprintln(os.Stderr, live.Safe(err))
		os.Exit(1)
	}
	s := &forwardutil.Example{Client: forward.NewClient(config.Options()...), Config: config}
	os.Exit(live.Execute(config, live.Scenario{Name: "drive", Description: "Drive 文件上传、分页、下载和清理。", Run: func(ctx context.Context, r *live.Run) error { return run(ctx, r, s) }}))
}

func run(ctx context.Context, r *live.Run, s *forwardutil.Example) error {
	id, err := s.Identity(ctx, r)
	if err != nil {
		return err
	}
	path := "sdk-example/说明.txt"
	r.Track("drive directory", "sdk-example", func(ctx context.Context) error {
		_, err := s.Client.Drives.DeleteEntry(ctx, forward.DriveDeleteEntryParams{IdentityID: id, Path: "sdk-example", Recursive: forward.Bool(true)})
		return err
	})
	// Independent storage client: do not copy QCA authentication or middleware.
	storage := &http.Client{Timeout: 30 * time.Second, CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	defer storage.CloseIdleConnections()
	transfer := func(method, url string, headers map[string]string, body io.Reader) ([]byte, error) {
		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, fmt.Errorf("invalid storage request")
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		res, err := storage.Do(req)
		if err != nil {
			return nil, fmt.Errorf("storage transfer failed")
		}
		defer res.Body.Close()
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return nil, fmt.Errorf("storage HTTP %d", res.StatusCode)
		}
		return io.ReadAll(io.LimitReader(res.Body, 1024))
	}
	r.Step("申请签名 URL，上传示例文件")
	upload, err := s.Client.Drives.CreateUploadURL(ctx, forward.DriveCreateUploadURLParams{IdentityID: id, Path: path, ContentType: forward.String("text/plain")})
	if err != nil {
		return err
	}
	if _, err = transfer(upload.Method, upload.URL, upload.Headers, strings.NewReader("Hello from QCA Drive")); err != nil {
		return err
	}
	r.Step("使用 Drive 游标列出目录")
	params := forward.DriveListEntriesParams{IdentityID: id, Path: forward.String("sdk-example"), Limit: forward.Int(20)}
	for {
		page, err := s.Client.Drives.ListEntries(ctx, params)
		if err != nil {
			return err
		}
		for _, entry := range page.Entries {
			r.Log("drive_entry", entry.Type+": "+entry.Path)
		}
		if page.NextPageToken == "" {
			break
		}
		params.PageToken = forward.String(page.NextPageToken)
	}
	r.Step("下载文件，然后删除示例条目")
	download, err := s.Client.Drives.CreateDownloadURL(ctx, forward.DriveCreateDownloadURLParams{IdentityID: id, Path: path})
	if err != nil {
		return err
	}
	data, err := transfer(download.Method, download.URL, download.Headers, nil)
	if err != nil {
		return err
	}
	r.Log("download", string(data))
	deleted, err := s.Client.Drives.DeleteEntry(ctx, forward.DriveDeleteEntryParams{IdentityID: id, Path: path})
	if err != nil {
		return err
	}
	r.Log("deleted_count", fmt.Sprint(deleted.DeletedCount))
	return nil
}
