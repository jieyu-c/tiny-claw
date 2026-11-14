package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"tiny-claw/internal/schema"
)

type Registry interface {
	Registry(tool BaseTool)
	GetAvailableTools() []*schema.ToolDefinition
	Exec(ctx context.Context, call *schema.ToolCall) schema.ToolResult
}

type BaseTool interface {
	Name() string
	Description() *schema.ToolDefinition
	Exec(ctx context.Context, args json.RawMessage) (string, error)
}

type registryImpl struct {
	tools map[string]BaseTool
}

func NewRegistry() Registry {
	return &registryImpl{
		tools: make(map[string]BaseTool),
	}
}

var _ Registry = new(registryImpl)

// Exec implements [Registry].
func (r *registryImpl) Exec(ctx context.Context, call *schema.ToolCall) (result schema.ToolResult) {
	tool, exixts := r.tools[call.Name]
	if !exixts {
		errMsg := fmt.Sprintf("Error: 系统中不存在名为 '%s' 的工具。", call.Name)
		return schema.ToolResult{
			ToolCallID: call.ID,
			Output:     errMsg,
			IsError:    true, // 标记为错误，模型看到后会尝试纠正
		}
	}
	toolResult, err := tool.Exec(ctx, call.Arguments)
	if err != nil {
		return schema.ToolResult{
			ToolCallID: call.ID,
			Output: err.Error(),
			IsError: true,
		}
	}
	result.IsError = false
	result.Output = toolResult
	result.ToolCallID = call.ID
	return
}

// GetAvailableTools implements [Registry].
func (r *registryImpl) GetAvailableTools() (result []*schema.ToolDefinition) {
	for _, tool := range r.tools {
		result = append(result, tool.Description())
	}
	return
}

// Registry implements [Registry].
func (r *registryImpl) Registry(tool BaseTool) {
	name := tool.Name()
	if _, exixts := r.tools[name]; exixts {
		log.Printf("[Warning] 工具 '%s' 已经被注册，将被覆盖。\n", name)
	}
	r.tools[name] = tool
	log.Printf("[Registry] 成功挂载工具: %s\n", name)

}
