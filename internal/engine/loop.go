package engine

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	ctxpkg "tiny-claw/internal/context"
	"tiny-claw/internal/provider"
	"tiny-claw/internal/schema"
	"tiny-claw/internal/tools"
)

type AgentEngine struct {
	provider       provider.LLMProvider
	registry       tools.Registry
	EnableThinking bool
	compactor      *ctxpkg.Compactor
	PlanMode       bool
}

func NewAgentEngine(p provider.LLMProvider, registry tools.Registry, enableThinking bool, planMode bool) *AgentEngine {
	return &AgentEngine{
		provider:       p,
		registry:       registry,
		EnableThinking: enableThinking,
		PlanMode:       planMode,
		compactor:      ctxpkg.NewCompactor(3000, 10),
	}
}

type indexedResult struct {
	index  int
	result schema.ToolResult
}

func (engine *AgentEngine) Run(ctx context.Context, session *Session, reporter Reporter) error {
	log.Println("starting agent engine")
	composer := ctxpkg.NewPromptComposer(session.WorkDir, engine.PlanMode)
	turnCount := 0

	for {
		memory := session.GetWorkingMemory(10)
		history := []*schema.Message{
			composer.Build(),
		}
		history = append(history, memory...)

		compactedContext := engine.compactor.Compact(history)
		if engine.EnableThinking {
			if reporter != nil {
				reporter.OnThinking(ctx)
			}
			thinkingResult, err := engine.provider.Generate(ctx, compactedContext, nil)
			if err != nil {
				log.Println(err)
				return fmt.Errorf("engine thinking failed: %w", err)
			}
			if thinkingResult.Content != "" {
				compactedContext = append(compactedContext, thinkingResult)
			}
		}
		turnCount++
		availableTools := engine.registry.GetAvailableTools()
		result, err := engine.provider.Generate(ctx, compactedContext, availableTools)
		if err != nil {
			return fmt.Errorf("模型生成失败: %w", err)
		}

		compactedContext = append(compactedContext, result)
		session.Append(result)

		if result.Content != "" && reporter != nil {
			reporter.OnMessage(ctx, result.Content)
		}

		if len(result.ToolCalls) == 0 {
			log.Println("[Engine] 任务完成，退出循环。")
			break
		}

		ch := make(chan indexedResult, 2)
		timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		engine.callTools(timeoutCtx, ch, result.ToolCalls, reporter)
		toolResults := make([]*schema.ToolResult, len(result.ToolCalls))
	loop:
		for {
			select {
			case r, ok := <-ch:
				if !ok {
					break loop
				}
				toolResults[r.index] = &r.result
			case <-timeoutCtx.Done():
				break loop
			}
		}
		for _, tool := range toolResults {
			withTools := schema.MessageWithTools(tool.Output, tool.ToolCallID)
			history = append(history, withTools)
			session.Append(withTools)
		}
		cancel()
	}
	return nil
}

func (e *AgentEngine) callTools(ctx context.Context, ch chan<- indexedResult, tools []*schema.ToolCall, reporter Reporter) {
	var wg sync.WaitGroup
	defer func() {
		wg.Wait()
	}()
	// 信号量， 最大并发量2
	sem := make(chan struct{}, 2)
	// 调用工具
	for i, call := range tools {
		select {
		case sem <- struct{}{}:
		case <-ctx.Done():
			return
		}
		wg.Add(1)
		// 这里不传call参数，在新版本的go中没有问题， 但是在1.22之前， 可能会出现执行的工具都是最后一次调用。
		go func(index int, call *schema.ToolCall) {
			defer wg.Done()
			defer func() {
				<-sem
			}()
			if reporter != nil {
				reporter.OnToolCall(ctx, call.Name, string(call.Arguments))
			}
			toolResult := e.registry.Exec(ctx, call)
			// 是否成功
			if toolResult.IsError {
				log.Printf(" -> ❌ 工具执行报错: %s\n", toolResult.Output)
			} else {
				log.Printf(" -> ✅ 工具执行成功 (返回 %d 字节)\n", len(toolResult.Output))
			}
			if reporter != nil {
				reporter.OnToolResult(ctx, call.Name, toolResult.Output, toolResult.IsError)
			}
			select {
			case ch <- indexedResult{
				index:  index,
				result: toolResult,
			}:
			case <-ctx.Done():
			}

		}(i, call)
	}
}

var wg sync.WaitGroup

func ClazzToSchema() {
	wg.Go(func() {
		fmt.Println("1111")
	})

	wg.Wait()
}
