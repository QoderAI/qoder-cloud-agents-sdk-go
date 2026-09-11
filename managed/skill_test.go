package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestSkillNew(t *testing.T) {
	client := contractClient(t, "SkillService", "New")
	response, err := client.Skills.New(context.Background(), managed.SkillNewParams{Files: []io.Reader{strings.NewReader("file")}})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Skill.New")
}

func TestSkillGet(t *testing.T) {
	client := contractClient(t, "SkillService", "Get")
	response, err := client.Skills.Get(context.Background(), "segment /?%#", managed.SkillGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Skill.Get")
}

func TestSkillList(t *testing.T) {
	client := contractClient(t, "SkillService", "List")
	response, err := client.Skills.List(context.Background(), managed.SkillListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Skill.List")
}

func TestSkillDelete(t *testing.T) {
	client := contractClient(t, "SkillService", "Delete")
	response, err := client.Skills.Delete(context.Background(), "segment /?%#", managed.SkillDeleteParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Skill.Delete")
}

func TestMultipartUploads(t *testing.T) {
	for _, kind := range []string{"file", "skill", "version"} {
		t.Run(kind, func(t *testing.T) {
			c := testClient(func(r *http.Request) (*http.Response, error) {
				typ, p, e := mime.ParseMediaType(r.Header.Get("Content-Type"))
				if e != nil || typ != "multipart/form-data" {
					t.Fatalf("content-type %s", r.Header.Get("Content-Type"))
				}
				reader := multipart.NewReader(r.Body, p["boundary"])
				fields := map[string][]string{}
				files := 0
				for {
					part, e := reader.NextPart()
					if e == io.EOF {
						break
					}
					if e != nil {
						t.Fatal(e)
					}
					b, e := io.ReadAll(part)
					if e != nil {
						t.Fatal(e)
					}
					if part.FileName() != "" {
						files++
						want := "files"
						if kind == "file" {
							want = "file"
						}
						if part.FormName() != want {
							t.Fatalf("wrong upload field %s", part.FormName())
						}
						_, header, _ := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
						if header["filename"] != "demo/SKILL.md" {
							t.Fatal("relative upload path lost")
						}
						if string(b) != "content" {
							t.Fatal("file bytes lost")
						}
					}
					fields[part.FormName()] = append(fields[part.FormName()], string(b))
				}
				if files != 1 {
					t.Fatalf("files %d", files)
				}
				if kind != "version" {
					if len(fields["metadata"]) != 1 || fields["metadata"][0] != `{"project":"cloud"}` {
						t.Fatalf("metadata must be JSON: %v", fields)
					}
				}
				if kind == "skill" && (len(fields["display_title"]) != 1 || fields["display_title"][0] != "Demo") {
					t.Fatal("skill display title mapping lost")
				}
				return reply(r, 200, `{}`), nil
			})
			f := convention.UploadFile{Reader: strings.NewReader("content"), Name: "demo/SKILL.md", MediaType: "text/markdown"}
			var e error
			switch kind {
			case "file":
				_, e = c.Files.Upload(context.Background(), managed.FileUploadParams{File: f, Metadata: map[string]string{"project": "cloud"}})
			case "skill":
				_, e = c.Skills.New(context.Background(), managed.SkillNewParams{Files: []io.Reader{f}, DisplayName: managed.String("Demo"), Metadata: map[string]string{"project": "cloud"}})
			case "version":
				_, e = c.Skills.Versions.New(context.Background(), "skill", managed.SkillVersionNewParams{Files: []io.Reader{f}})
			}
			if e != nil {
				t.Fatal(e)
			}
		})
	}
}
