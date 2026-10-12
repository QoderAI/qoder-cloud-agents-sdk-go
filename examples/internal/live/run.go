// Package live contains the small amount of plumbing shared by the runnable examples.
package live

import (
	"bufio"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
)

type Config struct {
	BaseURL, PAT, Model   string
	Mode, EnvFile, Output string
	Timeout               time.Duration
	AuthMode              string
	credential            convention.Credential
	accessToken, secret   string
}

// ReadEnv reads assignments as data, never as shell commands. Existing environment
// variables take precedence; callers can use either separate PATs or one shared PAT.
func ReadEnv(path string) (map[string]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	values := map[string]string{}
	s := bufio.NewScanner(f)
	for line := 1; s.Scan(); line++ {
		text := strings.TrimSpace(s.Text())
		if text == "" || strings.HasPrefix(text, "#") {
			continue
		}
		text = strings.TrimPrefix(text, "export ")
		key, value, ok := strings.Cut(text, "=")
		if !ok {
			return nil, fmt.Errorf("invalid env assignment at line %d", line)
		}
		key, value = strings.TrimSpace(key), strings.TrimSpace(value)
		if strings.HasPrefix(value, "\"") || strings.HasPrefix(value, "'") {
			quote := value[0]
			end := strings.IndexByte(value[1:], quote)
			if end < 0 {
				return nil, fmt.Errorf("unterminated env value at line %d", line)
			}
			value = value[1 : 1+end]
		} else if end := strings.Index(value, " #"); end >= 0 {
			value = strings.TrimSpace(value[:end])
		}
		values[key] = value
	}
	return values, s.Err()
}

func Load(mode string) (Config, error) {
	file := flag.String("env", ".env.live", "configuration file, relative to the working directory")
	region := flag.String("region", "", "API region: cn or global (international is an alias)")
	timeout := flag.Duration("timeout", 5*time.Minute, "deadline per scenario, excluding cleanup")
	format := flag.String("output", "text", "output format: text or json")
	flag.Parse()
	if *format != "text" && *format != "json" {
		return Config{}, errors.New("output must be text or json")
	}
	values, err := ReadEnv(*file)
	if err != nil {
		return Config{}, err
	}
	get := func(key string) string {
		if v, ok := os.LookupEnv(key); ok {
			return v
		}
		return values[key]
	}
	prefix := "QODER_" + strings.ToUpper(mode) + "_"
	base := get(prefix + "BASE_URL")
	if base == "" {
		suffix := mode
		if mode == "managed" {
			suffix = "cloud"
		}
		base = "https://api.qoder.com.cn/api/v1/" + suffix
	}
	u, err := url.Parse(base)
	if err != nil || u.Scheme != "https" || u.User != nil || u.RawQuery != "" || (u.Host != "api.qoder.com.cn" && u.Host != "api.qoder.com") {
		return Config{}, errors.New("examples require an HTTPS api.qoder.com.cn or api.qoder.com base URL without credentials or query parameters")
	}
	selectedRegion := *region
	if selectedRegion == "" && get("QODER_AUTH_MODE") == "sa" {
		selectedRegion = get("QODER_SA_REGION")
	}
	if selectedRegion == "" {
		selectedRegion = "cn"
	}
	switch selectedRegion {
	case "cn":
		u.Host = "api.qoder.com.cn"
	case "global", "international":
		u.Host = "api.qoder.com"
	default:
		return Config{}, errors.New("region must be cn or global (international is an alias)")
	}
	if *timeout <= 0 {
		return Config{}, errors.New("timeout must be positive")
	}
	config := Config{Mode: mode, EnvFile: *file, Output: *format, BaseURL: u.String(), Model: get(prefix + "MODEL"), Timeout: *timeout}
	if err := config.configureAuthentication(get); err != nil {
		return Config{}, err
	}
	return config, nil
}

func (c *Config) configureAuthentication(get func(string) string) error {
	c.AuthMode = get("QODER_AUTH_MODE")
	if c.AuthMode == "" {
		c.AuthMode = "pat"
	}
	if c.AuthMode == "pat" {
		c.PAT = get("QODER_" + strings.ToUpper(c.Mode) + "_PAT")
		if c.PAT == "" {
			c.PAT = get("QODER_PAT")
		}
		if strings.TrimSpace(c.PAT) == "" {
			return fmt.Errorf("configure QODER_%s_PAT or QODER_PAT", strings.ToUpper(c.Mode))
		}
		c.secret = c.PAT
		return nil
	}
	if c.AuthMode != "sa" {
		return errors.New("QODER_AUTH_MODE must be pat or sa")
	}
	key, token := get("QODER_SA_KEY"), get("QODER_ACCESS_TOKEN")
	if (strings.TrimSpace(key) == "") == (strings.TrimSpace(token) == "") {
		return errors.New("SA examples require exactly one of QODER_SA_KEY and QODER_ACCESS_TOKEN")
	}
	if strings.TrimSpace(token) != "" {
		c.accessToken, c.secret = token, token
		return nil
	}
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return errors.New("invalid example base URL")
	}
	region := "global"
	if u.Hostname() == "api.qoder.com.cn" {
		region = "cn"
	}
	scopes := []string{"qca.access"}
	if c.Mode == "forward" {
		scopes = append(scopes, "forward.access")
	}
	c.credential, err = convention.NewServiceAccountCredential(convention.ServiceAccountCredentialConfig{
		SAKey: key, Region: region, Scopes: scopes,
	})
	c.secret = key
	return err
}

func (c Config) Options() []option.RequestOption {
	options := []option.RequestOption{option.WithBaseURL(c.BaseURL), option.WithMaxRetries(0), option.WithRequestTimeout(30 * time.Second)}
	if c.AuthMode == "sa" {
		options = append(options, option.WithAuthMode("sa"))
		if c.credential != nil {
			return append(options, option.WithCredential(c.credential))
		}
		return append(options, option.WithAuthToken(c.accessToken))
	}
	return append(options, option.WithPAT(c.PAT))
}

type cleanup struct {
	kind, id string
	fn       func(context.Context) error
}
type Run struct {
	Name     string
	cleanups []cleanup
	output   *output
	step     int
	action   string
}
type Scenario struct {
	Name, Description string
	Run               func(context.Context, *Run) error
}

func (r *Run) Track(kind, id string, fn func(context.Context) error) {
	r.cleanups = append(r.cleanups, cleanup{kind: kind, id: id, fn: fn})
	r.Log("created", resourceName(kind)+"："+id)
}

// RetainResources transfers cleanup responsibility to the caller after a
// successful conversation that explicitly requested -keep-session.
func (r *Run) RetainResources() {
	r.Log("info", fmt.Sprintf("按 -keep-session 保留本次创建的 %d 项资源；后续需手动清理上面列出的资源。", len(r.cleanups)))
	r.cleanups = nil
}

// Safe hides response bodies, tokens and signed URLs, even through wrapped errors.
func Safe(err error) string {
	var api *apierror.Error
	if errors.As(err, &api) {
		return fmt.Sprintf("HTTP %d code=%s type=%s request_id=%s", api.StatusCode, api.Code, api.Type(), api.RequestID)
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return "HTTP transport failed: " + Safe(urlErr.Err)
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return "deadline exceeded"
	}
	if errors.Is(err, context.Canceled) {
		return "canceled"
	}
	return err.Error() // Other errors are locally authored assertions, never raw bodies.
}

var webURL = regexp.MustCompile(`https?://[^\s"'<>]+`)

func Explain(err error, pat string) string {
	message := Safe(err)
	var api *apierror.Error
	if errors.As(err, &api) {
		if api.Request != nil {
			message += " " + api.Request.Method + " " + api.Request.URL.EscapedPath()
		}
		message += " " + api.Message
	}
	if pat != "" {
		message = strings.ReplaceAll(message, pat, "[REDACTED]")
	}
	message = webURL.ReplaceAllString(message, "[URL REDACTED]")
	if len(message) > 1600 {
		message = message[:1600]
	}
	return message
}

func Execute(c Config, scenario Scenario) int {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return execute(ctx, c, scenario, os.Stdout)
}

func execute(ctx context.Context, c Config, s Scenario, writer io.Writer) (code int) {
	secret := c.secret
	if secret == "" {
		secret = c.PAT
	}
	out := &output{writer: writer, format: c.Output, pat: secret}
	mode := map[string]string{"forward": "Forward", "managed": "Managed"}[c.Mode]
	if mode == "" {
		mode = "SDK"
	}
	out.write("", "header", fmt.Sprintf("%s 示例\n配置：%s\n地址：%s\n每个场景最长等待：%s", mode, c.EnvFile, c.BaseURL, c.Timeout))
	started := time.Now()
	r := &Run{Name: s.Name, output: out}
	sceneFailed := false
	defer func() {
		if recovered := recover(); recovered != nil {
			sceneFailed = true
			r.Log("failed", fmt.Sprintf("步骤「%s」异常中止（%T）", r.action, recovered))
		}
		cleanFailed := 0
		if len(r.cleanups) > 0 {
			r.Log("cleanup_start", fmt.Sprintf("清理本次创建的资源（%d 项）", len(r.cleanups)))
		}
		for i := len(r.cleanups) - 1; i >= 0; i-- {
			step := r.cleanups[i]
			cleanCtx, done := context.WithTimeout(context.Background(), 90*time.Second)
			err := step.fn(cleanCtx)
			done()
			label := cleanupLabel(c.Mode, step.kind, step.id)
			if err != nil {
				cleanFailed++
				r.Log("cleanup_failed", label+"\n      "+Explain(err, secret))
			} else {
				r.Log("cleaned", label)
			}
		}
		status, passed, failed := "通过", 1, 0
		if sceneFailed || cleanFailed > 0 {
			status, passed, failed = "失败", 0, 1
			code = 1
		}
		r.Log("scenario_result", fmt.Sprintf("场景 %s：%s · 用时 %.1f 秒 · 清理 %d/%d 项", s.Name, status, time.Since(started).Seconds(), len(r.cleanups)-cleanFailed, len(r.cleanups)))
		out.write("summary", "summary", fmt.Sprintf("运行结束：1 个场景，%d 个通过，%d 个失败；清理失败 %d 项；总用时 %.1f 秒。", passed, failed, cleanFailed, time.Since(started).Seconds()))
	}()
	if ctx.Err() != nil {
		sceneFailed = true
		r.Log("failed", "运行已中断")
		return 1
	}
	title := "── " + s.Name + " ──"
	if s.Description != "" {
		title += "\n" + s.Description
	}
	r.Log("start", title)
	sceneCtx, cancel := context.WithTimeout(ctx, c.Timeout)
	defer cancel()
	if err := s.Run(sceneCtx, r); err != nil {
		sceneFailed = true
		r.Log("failed", "步骤「"+r.action+"」失败："+Explain(err, secret))
	} else {
		r.Log("verified", "场景验证完成")
	}
	return 0
}

func Name(prefix string) string { return "sdk-example-" + prefix + "-" + Marker() }
func Marker() string {
	var b [8]byte
	if _, err := rand.Read(b[:]); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b[:])
}
func Pause(ctx context.Context) error {
	timer := time.NewTimer(2 * time.Second)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func ChooseModel(configured string, enabled []string) (string, error) {
	if configured != "" {
		for _, id := range enabled {
			if id == configured {
				return id, nil
			}
		}
		return "", errors.New("configured model is not enabled for this account")
	}
	for _, preferred := range []string{"qoder-lite", "lite", "qoder-plus", "plus"} {
		for _, id := range enabled {
			if id == preferred {
				return id, nil
			}
		}
	}
	if len(enabled) > 0 {
		return enabled[0], nil
	}
	return "", errors.New("no enabled models returned")
}
