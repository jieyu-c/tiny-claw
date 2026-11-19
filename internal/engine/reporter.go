package engine

import "context"

type Reporter interface {
	// 思考的时候
	OnThinking(ctx context.Context)
	// 调用工具的时候
	OnToolCall(ctx context.Context, toolName, args string)
	// 工具返回结果的时候
	OnToolResult(ctx context.Context, toolName, result string, isError bool)
	// 当模型回答问题向用户输出内容的时候
	OnMessage(ctx context.Context, content string)
}
