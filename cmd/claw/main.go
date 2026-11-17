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
	engine := engine2.NewAgentEngine(p, r, workDir, false)
	prompt := `
 我当前目录下有 a.txt, b.txt, c.txt 三个文件。不要用命令读文件， 用读文件的工具
 为了节省时间，请你同时一次性读取这三个文件，并将它们的内容综合起来，告诉我它们分别记录了什么领域的信息。
 `

	// TODO 5. 开始执行
	err := engine.Run(context.Background(), prompt, nil)

	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}


	log.Println("Welcome to claw!")
}
