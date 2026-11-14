package main

import (
	"context"
	"log"
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
	//registry.Register(tools.NewBashTool())

	// TODO 3. 初始化上下文
	// ctx := context.NewManager(...)

	// TODO 4. 初始化执行引擎
	engine := engine2.NewAgentEngine(p, r, "~", false)

	// TODO 5. 开始执行
	err := engine.Run(context.TODO(), "我在北京， 帮我查查北京今天天气如何？")

	if err != nil {
		log.Fatalf("引擎运行崩溃: %v", err)
	}

	log.Println("Welcome to claw!")
}
