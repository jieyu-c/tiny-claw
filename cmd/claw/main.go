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
	p := provider.NewOpenAI("deepseek-v4-pro")

	// TODO 2. 注册工具
	r := tools.NewRegistry()
	workDir, _ := os.Getwd()

	r.Registry(tools.NewWriteFileTool(workDir))
	r.Registry(tools.NewReadFileTool(workDir))
	r.Registry(tools.NewBashTool(workDir))

	// TODO 3. 初始化上下文
	// ctx := context.NewManager(...)

	// TODO 4. 初始化执行引擎
	engine := engine2.NewAgentEngine(p, r, "~", false)
	prompt := `
 请帮我执行以下操作：
 1. 用 bash 查看一下我当前电脑的 Go 版本。
 2. 帮我写一个简单的 helloworld.go 文件，输出 "Hello, go-tiny-claw!"。
 3. 用 bash 编译并运行这个 go 文件，确认它能正常工作。
 `

	// TODO 5. 开始执行
	err := engine.Run(context.TODO(), prompt)

	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}

	log.Println("Welcome to claw!")
}
