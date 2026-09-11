package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/managedutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func main() {
	config, err := live.Load("managed")
	if err != nil {
		fmt.Fprintln(os.Stderr, "配置错误："+live.Safe(err))
		os.Exit(1)
	}
	s := &managedutil.Example{Client: managed.NewClient(config.Options()...), Config: config}
	os.Exit(live.Execute(config, live.Scenario{
		Name:        "resources",
		Description: "演示文件、环境变量和 Skill 在会话中的使用。",
		Run: func(ctx context.Context, r *live.Run) error {
			return run(ctx, r, s)
		},
	}))
}

// run proves file mounting, environment variables and custom Skills.
// Expected values live only in resources, never in the user prompts.
func run(ctx context.Context, r *live.Run, s *managedutil.Example) error {
	env, err := s.Environment(ctx, r)
	if err != nil {
		return err
	}
	fileToken, envToken, skillToken := live.Marker(), live.Marker(), live.Marker()
	r.Step("上传示例文件，供会话中的工具读取")
	file, err := s.Client.Files.Upload(ctx, managed.FileUploadParams{File: convention.UploadFile{Name: "sdk-example.txt", Reader: strings.NewReader(fileToken)}})
	if err != nil {
		return err
	}
	r.Track("file", file.ID, func(ctx context.Context) error {
		_, err := s.Client.Files.Delete(ctx, file.ID, managed.FileDeleteParams{})
		return err
	})
	skillName := live.Name("skill")
	r.Step("上传自定义 Skill，其中包含一个随机校验值")
	skill, err := s.Client.Skills.New(ctx, managed.SkillNewParams{Files: []io.Reader{convention.UploadFile{Name: skillName + "/SKILL.md", Reader: strings.NewReader(fmt.Sprintf("---\nname: %s\ndescription: Provides a sample verification code for the SDK example.\n---\nThe example verification value EXAMPLE_SKILL_CODE is: %s\n", skillName, skillToken))}}})
	if err != nil {
		return err
	}
	r.Track("skill", skill.ID, func(ctx context.Context) error {
		_, err := s.Client.Skills.Delete(ctx, skill.ID, managed.SkillDeleteParams{})
		return err
	})
	agent, err := s.Agent(ctx, r, managed.AgentNewParams{Skills: []managed.ManagedAgentsSkillParamsUnion{{OfCustom: &managed.ManagedAgentsCustomSkillParams{Type: "custom", SkillID: skill.ID, Version: managed.String(skill.LatestVersionID)}}}})
	if err != nil {
		return err
	}
	session, err := s.NewSession(ctx, r, managed.SessionNewParams{EnvironmentID: env, Agent: managed.SessionNewParamsAgentUnion{OfString: managed.String(agent)}, EnvironmentVariables: map[string]string{"SDK_EXAMPLE_VALUE": envToken}, Resources: []managed.SessionNewParamsResourceUnion{
		{OfFile: &managed.ManagedAgentsFileResourceParams{Type: "file", FileID: file.ID, MountPath: managed.String("/data/workspace/sdk-example.txt")}},
	}})
	if err != nil {
		return err
	}
	if err = s.Turn(ctx, r, session, "请使用工具读取 /data/workspace/sdk-example.txt 的内容和 SDK_EXAMPLE_VALUE 环境变量，分别返回这两个示例值。", []string{fileToken, envToken}, true, false); err != nil {
		return err
	}
	return s.Turn(ctx, r, session, "请使用技能 "+skillName+" 读取 EXAMPLE_SKILL_CODE，返回这个示例校验码。", []string{skillToken}, true, false)
}
