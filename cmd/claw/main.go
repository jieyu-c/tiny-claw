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
	r.Registry(tools.NewReadFileTool(workDir))

	// TODO 3. 初始化上下文
	// ctx := context.NewManager(...)

	// TODO 4. 初始化执行引擎
	engine := engine2.NewAgentEngine(p, r, "~", false)

	prompt := "请调用工具读取一下当前工作区目录下 hello.txt 文件的内容，并用一句话向我总结它说了什么。"

	// TODO 5. 开始执行
	err := engine.Run(context.TODO(), prompt)

	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}

	log.Println("Welcome to claw!")
}
