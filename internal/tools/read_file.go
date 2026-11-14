package tools

import (
	"context"
	"encoding/json"
	"tiny-claw/internal/schema"
)

type ReadFileTool struct {
	workDir string
}

var _ BaseTool = new(ReadFileTool)

// Description implements [BaseTool].
func (r *ReadFileTool) Description() *schema.ToolDefinition {
	return &schema.ToolDefinition{
		Name:        "read_file",
		Description: "读取指定路径的文件内容。请提供相对工作区的路径。",
		// 可以用反射优化一下， 这里这么写性能可能比较好， 但是开发成本较高。
		// 用反射性能不会下滑很多。因为是初始化时使用反射注入一次， 不会在运行期持续使用反射
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "要读取的文件路径，如 cmd/claw/main.go",
				},
			},
			"required": []string{"path"},
		},
	}
}

// readFileArgs 内部定义用于反序列化的结构体
type readFileArgs struct {
	Path string `json:"path"`
}

// Exec implements [BaseTool].
func (r *ReadFileTool) Exec(ctx context.Context, args json.RawMessage) (string, error) {
	panic("unimplemented")
}

// Name implements [BaseTool].
func (r *ReadFileTool) Name() string {
	return "read_file"
}
