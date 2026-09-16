// Package managed implements the Qoder Cloud Agents managed API.
// Service methods and entity fields follow the managed API contract.
package managed

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"os"
	"slices"
)

const DefaultBaseURL = "https://api.qoder.com/api/v1/cloud/"

// Client exposes concrete services; nested services follow the upstream hierarchy.
// A client can be shared between goroutines. Configure options before first use.
type Client struct {
	Options        []option.RequestOption
	Agents         AgentService
	Sessions       SessionService
	Deployments    DeploymentService
	DeploymentRuns DeploymentRunService
	Dreams         DreamService
	Environments   EnvironmentService
	Skills         SkillService
	Vaults         VaultService
	Files          FileService
	MemoryStores   MemoryStoreService
	Models         ModelService
}

// NewClient reads QODER_PAT and QODER_BASE_URL. Explicit options win.
func NewClient(opts ...option.RequestOption) Client {
	defaults := []option.RequestOption{convention.WithDefaultBaseURL(DefaultBaseURL)}
	if base := os.Getenv("QODER_BASE_URL"); base != "" {
		defaults = append(defaults, option.WithBaseURL(base))
	}
	if token := os.Getenv("QODER_PAT"); token != "" {
		defaults = append(defaults, option.WithPAT(token))
	}
	defaults = append(defaults, convention.RequireCredential())
	opts = slices.Concat(defaults, opts)
	return Client{Options: opts, Agents: NewAgentService(opts...), Sessions: NewSessionService(opts...), Deployments: NewDeploymentService(opts...), DeploymentRuns: NewDeploymentRunService(opts...), Dreams: NewDreamService(opts...), Environments: NewEnvironmentService(opts...), Skills: NewSkillService(opts...), Vaults: NewVaultService(opts...), Files: NewFileService(opts...), MemoryStores: NewMemoryStoreService(opts...), Models: NewModelService(opts...)}
}

func String(v string) param.Opt[string]  { return param.NewOpt(v) }
func Int(v int64) param.Opt[int64]       { return param.NewOpt(v) }
func Bool(v bool) param.Opt[bool]        { return param.NewOpt(v) }
func Float(v float64) param.Opt[float64] { return param.NewOpt(v) }
