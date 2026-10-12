//go:build live

package integration_test

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func saEnvironment(path string, environment map[string]string, execution bool) (map[string]string, error) {
	values := map[string]string{}
	file, err := os.Open(path)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("cannot read SA env file")
	}
	if err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		assignment := regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_]*$`)
		for line := 1; scanner.Scan(); line++ {
			text := strings.TrimSpace(scanner.Text())
			if text == "" || strings.HasPrefix(text, "#") {
				continue
			}
			key, value, ok := strings.Cut(strings.TrimPrefix(text, "export "), "=")
			key, value = strings.TrimSpace(key), strings.TrimSpace(value)
			if !ok || !assignment.MatchString(key) {
				return nil, fmt.Errorf("invalid env assignment at line %d", line)
			}
			if strings.HasPrefix(value, `"`) || strings.HasPrefix(value, "'") {
				end := strings.IndexByte(value[1:], value[0])
				if end < 0 {
					return nil, fmt.Errorf("unterminated env value at line %d", line)
				}
				value = value[1 : end+1]
			} else if end := strings.Index(value, " #"); end >= 0 {
				value = strings.TrimSpace(value[:end])
			}
			values[key] = value
		}
		if scanner.Err() != nil {
			return nil, fmt.Errorf("cannot read SA env file")
		}
	}
	for key, value := range environment {
		values[key] = value
	}
	for _, name := range []string{"QODER_SA_CN_KEY", "QODER_SA_GLOBAL_KEY"} {
		if strings.TrimSpace(values[name]) == "" {
			return nil, fmt.Errorf("%s is required", name)
		}
	}
	if execution {
		for _, name := range []string{"QODER_SA_LIVE_ALLOW_WRITE", "QODER_SA_LIVE_ALLOW_EXECUTION"} {
			if values[name] != "true" {
				return nil, fmt.Errorf("%s=true is required", name)
			}
		}
	}
	return values, nil
}

func saRequireConfiguration(t *testing.T, execution bool) {
	t.Helper()
	path := os.Getenv("LIVE_ENV_FILE")
	if path == "" {
		path = filepath.Join("..", ".env.live")
	}
	environment := map[string]string{}
	for _, entry := range os.Environ() {
		key, value, _ := strings.Cut(entry, "=")
		environment[key] = value
	}
	values, err := saEnvironment(path, environment, execution)
	if err != nil {
		t.Fatal(err)
	}
	for key, value := range values {
		if strings.HasPrefix(key, "QODER_SA_") {
			t.Setenv(key, value)
		}
	}
}

func TestSAConfiguration(t *testing.T) {
	path := filepath.Join(t.TempDir(), "sa.env")
	content := "QODER_SA_CN_KEY='file-cn'\nexport QODER_SA_GLOBAL_KEY=file-global # comment\nQODER_SA_MODEL=$(do-not-execute)\n"
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	values, err := saEnvironment(path, map[string]string{"QODER_SA_CN_KEY": "process-cn"}, false)
	if err != nil || values["QODER_SA_CN_KEY"] != "process-cn" || values["QODER_SA_MODEL"] != "$(do-not-execute)" {
		t.Fatal("env-file loading must preserve process precedence and literal values")
	}
	for _, execution := range []bool{false, true} {
		_, err = saEnvironment(path, map[string]string{"QODER_SA_GLOBAL_KEY": ""}, execution)
		if err == nil || err.Error() != "QODER_SA_GLOBAL_KEY is required" {
			t.Fatal("a missing Global key must fail before either region can run")
		}
	}
	for _, flag := range []string{"QODER_SA_LIVE_ALLOW_WRITE", "QODER_SA_LIVE_ALLOW_EXECUTION"} {
		env := map[string]string{"QODER_SA_LIVE_ALLOW_WRITE": "true", "QODER_SA_LIVE_ALLOW_EXECUTION": "true", flag: "false"}
		_, err = saEnvironment(path, env, true)
		if err == nil || err.Error() != flag+"=true is required" {
			t.Fatal("execution must require both flags")
		}
	}
	_, err = saEnvironment(path, map[string]string{"QODER_SA_LIVE_ALLOW_WRITE": "true", "QODER_SA_LIVE_ALLOW_EXECUTION": "true"}, true)
	if err != nil {
		t.Fatal(err)
	}
}
