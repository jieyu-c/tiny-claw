package engine

import (
	"context"
	"fmt"
	"log"
	"tiny-claw/internal/provider"
	"tiny-claw/internal/schema"
	"tiny-claw/internal/tools"
)

type AgentEngine struct {
	provider       provider.LLMProvider
	registry       tools.Registry
	WorkDir        string
	EnableThinking bool
}

func NewAgentEngine(p provider.LLMProvider, registry tools.Registry, workDir string, enableThinking bool) *AgentEngine {
	return &AgentEngine{
		provider:       p,
		registry:       registry,
		WorkDir:        workDir,
		EnableThinking: enableThinking,
	}
}

func (engine *AgentEngine) Run(ctx context.Context, input string) error {
	log.Println("starting agent engine")
	history := []*schema.Message{
		schema.SystemMessage("You are Shaosheng LIU, an expert coding assistant. You have full access to tools in the workspace."),
		schema.UserMessage(input),
	}

	turnCount := 0

	for {
		if engine.EnableThinking {
			thinkingResult, err := engine.provider.Generate(ctx, history, nil)
			if err != nil {
				log.Println(err)
				return fmt.Errorf("engine thinking failed: %w", err)
			}
			if thinkingResult.Content != "" {
				history = append(history, thinkingResult)
			}
		}
		turnCount++
		tools := engine.registry.GetAvailableTools()
		result, err := engine.provider.Generate(ctx, history, tools)
		if err != nil {
			return fmt.Errorf("模型生成失败: %w", err)
		}

		history = append(history, result)

		if result.Content != "" {
			fmt.Println(result.Content)
		}

		if len(result.ToolCalls) == 0 {
			log.Println("[Engine] 任务完成，退出循环。")
			break
		}
		// 调用工具
		for _, call := range result.ToolCalls {
			toolResult := engine.registry.Exec(ctx, call)
			// 是否成功
			if toolResult.IsError {
				log.Printf(" -> ❌ 工具执行报错: %s\n", toolResult.Output)
			} else {
				log.Printf(" -> ✅ 工具执行成功 (返回 %d 字节)\n", len(toolResult.Output))
			}
			history = append(history, schema.MessageWithTools(toolResult.Output, toolResult.ToolCallID))
		}
	}
	return nil
}
