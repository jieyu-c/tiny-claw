package main

import (
	"context"
	"log"
	"os"
	engine2 "tiny-claw/internal/engine"
	"tiny-claw/internal/provider"
	"tiny-claw/internal/tools"
)

func main() {
	log.Println("Welcome to claw!")

	// TODO 1. 初始化模型
	p := provider.NewOpenAI("glm-5.3")

	// TODO 2. 注册工具
	r := tools.NewRegistry()
	workDir, _ := os.Getwd()

	r.Registry(tools.NewWriteFileTool(workDir))
	r.Registry(tools.NewReadFileTool(workDir))
	r.Registry(tools.NewBashTool(workDir))
	r.Registry(tools.NewEditFileTool(workDir))

	// TODO 3. 初始化上下文
	// ctx := context.NewManager(...)§
	// TODO 4. 初始化执行引擎
	engine := engine2.NewAgentEngine(p, r, workDir, true)
	prompt := `
 我需要在当前目录下新建一个 ping.go，提供一个简单的 http ping 接口。
 写完之后，帮我把代码用 git 提交一下。
 `

	reporter := engine2.NewTerminalReporter()
	err := engine.Run(context.Background(), prompt, reporter)
	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}
}
