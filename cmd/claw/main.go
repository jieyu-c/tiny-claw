package main

import (
	"log"
	"net/http"
	"os"
	engine2 "tiny-claw/internal/engine"
	"tiny-claw/internal/feishu"
	"tiny-claw/internal/provider"
	"tiny-claw/internal/tools"

	"github.com/larksuite/oapi-sdk-go/v3/core/httpserverext"
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
	// 2. 初始化飞书 Bot 调度器
	bot := feishu.NewFeishuBot(engine)
	handler := httpserverext.NewEventHandlerFunc(bot.GetEventDispatcher())
	// 3. 注册路由并启动 HTTP 服务
	http.HandleFunc("/webhook/event", handler)
	port := ":48080"
	log.Printf("🚀 go-tiny-claw 飞书服务端已启动，正在监听 %s 端口\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
