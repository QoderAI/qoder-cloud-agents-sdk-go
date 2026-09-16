package forward

import (
	"os"
	"slices"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
)

const DefaultBaseURL = "https://api.qoder.com/api/v1/forward/"

// Client exposes Forward resources. Configure options before sharing a client.
type Client struct {
	Options         []option.RequestOption
	Templates       TemplateService
	Identities      IdentityService
	Sessions        SessionService
	Schedules       ScheduleService
	ScheduleRuns    ScheduleRunService
	Batches         BatchService
	Channels        ChannelService
	ChannelPairings ChannelPairingService
	Environments    EnvironmentService
	Files           FileService
	Skills          SkillService
	Vaults          VaultService
	MemoryStores    MemoryStoreService
	Models          ModelService
}

// NewClient reads QODER_PAT and QODER_FORWARD_BASE_URL. Explicit options win.
func NewClient(opts ...option.RequestOption) Client {
	defaults := []option.RequestOption{convention.WithDefaultBaseURL(DefaultBaseURL)}
	if token := os.Getenv("QODER_PAT"); token != "" {
		defaults = append(defaults, option.WithPAT(token))
	}
	if base := os.Getenv("QODER_FORWARD_BASE_URL"); base != "" {
		defaults = append(defaults, option.WithBaseURL(base))
	}
	defaults = append(defaults, convention.RequireCredential())
	opts = slices.Concat(defaults, opts)
	return Client{Options: opts,
		Templates:       NewTemplateService(opts...),
		Identities:      NewIdentityService(opts...),
		Sessions:        NewSessionService(opts...),
		Schedules:       NewScheduleService(opts...),
		ScheduleRuns:    NewScheduleRunService(opts...),
		Batches:         NewBatchService(opts...),
		Channels:        NewChannelService(opts...),
		ChannelPairings: NewChannelPairingService(opts...),
		Environments:    NewEnvironmentService(opts...),
		Files:           NewFileService(opts...),
		Skills:          NewSkillService(opts...),
		Vaults:          NewVaultService(opts...),
		MemoryStores:    NewMemoryStoreService(opts...),
		Models:          NewModelService(opts...),
	}
}
