package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
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
		Name:        "resources",
		Description: "演示文件挂载、Identity 环境变量覆盖和自定义 Skill 的读取。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s)
		},
	}))
}

// run proves file mounting, Identity overrides and custom Skill access. Expected values live only in resources, never in the user prompts.
func run(ctx context.Context, r *live.Run, s *forwardutil.Example) error {
	env, err := s.Environment(ctx, r)
	if err != nil {
		return err
	}
	identity, err := s.Identity(ctx, r)
	if err != nil {
		return err
	}
	fileToken, envToken, skillToken := live.Marker(), live.Marker(), live.Marker()
	r.Step("上传示例文件，供会话中的工具读取")
	file, err := s.Client.Files.Upload(ctx, forward.FileUploadParams{File: convention.UploadFile{Name: "sdk-example.txt", Reader: strings.NewReader(fileToken)}, Purpose: forward.String("session_resource")})
	if err != nil {
		return err
	}
	r.Track("file", file.ID, func(ctx context.Context) error { return s.Client.Files.Delete(ctx, file.ID) })
	skillName := live.Name("skill")
	r.Step("上传自定义 Skill，其中包含一个随机校验值")
	skill, err := s.Client.Skills.New(ctx, forward.SkillNewParams{Files: []io.Reader{convention.UploadFile{Name: skillName + "/SKILL.md", Reader: strings.NewReader(fmt.Sprintf("---\nname: %s\ndescription: Provides a sample verification code for the SDK example.\n---\nThe example verification value EXAMPLE_SKILL_CODE is: %s\n", skillName, skillToken))}}})
	if err != nil {
		return err
	}
	r.Track("skill", skill.ID, func(ctx context.Context) error { return s.Client.Skills.Delete(ctx, skill.ID) })
	template, err := s.Template(ctx, r, forward.TemplateNewParams{EnvironmentID: env, Skills: []forward.SkillBindingParam{{Type: "custom", SkillID: skill.ID, Version: forward.String(skill.LatestVersion)}}, EnvironmentVariables: forward.EnvironmentVariablesUnionParam{OfMap: map[string]any{"SDK_EXAMPLE_VALUE": "template-default"}}})
	if err != nil {
		return err
	}
	r.Step("设置 Identity 的环境变量，覆盖模板中的默认值")
	_, err = s.Client.Identities.Configs.Upsert(ctx, identity, template, forward.IdentityConfigUpsertParams{IdentityConfig: forward.IdentityConfigSpecParam{EnvironmentVariables: map[string]forward.EnvironmentVariableOverrideParam{"SDK_EXAMPLE_VALUE": {Op: "set", Value: forward.String(envToken)}}}})
	if err != nil {
		return err
	}
	session, err := s.NewSession(ctx, r, forward.SessionNewParams{IdentityID: identity, TemplateID: template, Resources: []forward.SessionResourceSpecParam{{Type: "file", FileID: file.ID, MountPath: forward.String("/data/workspace/sdk-example.txt")}}})
	if err != nil {
		return err
	}
	if err = s.Turn(ctx, r, session, "请使用工具读取 /data/workspace/sdk-example.txt 的内容和 SDK_EXAMPLE_VALUE 环境变量，分别返回这两个示例值。", []string{fileToken, envToken}, true, false); err != nil {
		return err
	}
	return s.Turn(ctx, r, session, "请使用技能 "+skillName+"，读取并返回其中的示例校验码 EXAMPLE_SKILL_CODE。", []string{skillToken}, true, false)
}
