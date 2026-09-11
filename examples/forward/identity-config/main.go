package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/forwardutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	config, err := live.Load("forward")
	if err != nil {
		fmt.Fprintln(os.Stderr, "配置错误："+live.Safe(err))
		os.Exit(1)
	}
	s := &forwardutil.Example{Client: forward.NewClient(config.Options()...), Config: config}
	os.Exit(live.Execute(config, live.Scenario{
		Name:        "identity-config",
		Description: "两个 Identity 共用一个 Template，读取各自生效配置，并在会话中验证个性化配置。",
		Run:         func(ctx context.Context, r *live.Run) error { return run(ctx, r, s) },
	}))
}

func run(ctx context.Context, r *live.Run, s *forwardutil.Example) error {
	env, err := s.Environment(ctx, r)
	if err != nil {
		return err
	}
	shared := "shared-" + live.Marker()
	baseline := "default-" + live.Marker()
	template, err := s.Template(ctx, r, forward.TemplateNewParams{
		EnvironmentID: env,
		EnvironmentVariables: forward.EnvironmentVariablesUnionParam{OfMap: map[string]any{
			"SDK_SHARED_VALUE": shared, "SDK_PERSONAL_VALUE": baseline,
		}},
	})
	if err != nil {
		return err
	}
	// Create both identities before either session so reverse cleanup removes all
	// sessions before their parent identities and the shared template.
	identities := make([]string, 2)
	values := []string{"alice-" + live.Marker(), "bob-" + live.Marker()}
	for i := range identities {
		identities[i], err = s.Identity(ctx, r)
		if err != nil {
			return err
		}
		r.Step(fmt.Sprintf("为 Identity %d 设置个性化环境变量", i+1))
		_, err = s.Client.Identities.Configs.Upsert(ctx, identities[i], template, forward.IdentityConfigUpsertParams{
			IdentityConfig: forward.IdentityConfigSpecParam{EnvironmentVariables: map[string]forward.EnvironmentVariableOverrideParam{
				"SDK_PERSONAL_VALUE": {Op: "set", Value: forward.String(values[i])},
			}},
		})
		if err != nil {
			return err
		}
		r.Step("查询最终生效配置，验证默认值继承与用户覆盖")
		effective, err := s.Client.Identities.Configs.GetEffective(ctx, identities[i], template)
		if err != nil {
			return err
		}
		vars := effective.Session.EnvironmentVariables
		if vars["SDK_SHARED_VALUE"] != shared || vars["SDK_PERSONAL_VALUE"] != values[i] {
			return fmt.Errorf("identity %d effective configuration did not preserve inheritance and override", i+1)
		}
		r.Log("checked", fmt.Sprintf("Identity %d：公共值=%s，个性化值=%s", i+1, shared, values[i]))
	}
	for i, identity := range identities {
		session, err := s.NewSession(ctx, r, forward.SessionNewParams{IdentityID: identity, TemplateID: template})
		if err != nil {
			return err
		}
		prompt := "请使用工具读取 SDK_SHARED_VALUE 和 SDK_PERSONAL_VALUE 两个环境变量，只返回这两个变量的实际值。"
		r.Step(fmt.Sprintf("在 Identity %d 的会话中读取实际配置", i+1))
		r.Message("user", prompt)
		after, err := s.Send(ctx, session, prompt)
		if err != nil {
			return err
		}
		reply, err := s.WaitReply(ctx, r, session, after, false)
		if err != nil {
			return err
		}
		if err := reply.Verify([]string{shared, values[i]}, true); err != nil {
			return err
		}
		if strings.Contains(reply.Text, values[1-i]) || strings.Contains(reply.Text, baseline) {
			return fmt.Errorf("identity %d reply contains another identity's value or the overridden default", i+1)
		}
		r.Log("checked", fmt.Sprintf("Identity %d 的会话读到了自己的覆盖值和模板公共值。", i+1))
	}
	return nil
}
