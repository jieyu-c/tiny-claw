package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"tiny-claw/internal/schema"
)

const maxLen = 5000;

type ReadFileTool struct {
	workDir string
}

func NewReadFileTool(workDir string) *ReadFileTool {
	return &ReadFileTool{
		workDir: workDir,
	}
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
	var input readFileArgs
	if err := json.Unmarshal(args, &input); err != nil {
		return "", fmt.Errorf("参数解析失败， %w", err)
	}
	// 正式环境需要校验 ../../这种地址， 防止目录穿越
	fullPath := filepath.Join(r.workDir, input.Path)

	file, err := os.Open(fullPath)
	if err != nil {
		return "", fmt.Errorf("打开文件失败， %w", err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("读取文件失败, %w", err)
	}
	if len(content) > maxLen {
		content = content[:maxLen]
	}

	return string(content), err

}

// Name implements [BaseTool].
func (r *ReadFileTool) Name() string {
	return "read_file"
}
