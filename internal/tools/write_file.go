package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"tiny-claw/internal/schema"
)

type WriteFileTool struct {
	workDir string
}

func NewWriteFileTool(workDir string) *WriteFileTool {
	return &WriteFileTool{
		workDir: workDir,
	}
}

var _ BaseTool = new(WriteFileTool)

// Description implements [BaseTool].
func (w *WriteFileTool) Description() *schema.ToolDefinition {
	return &schema.ToolDefinition{
		Name:        w.Name(),
		Description: "创建或覆盖写入一个文件。如果目录不存在会自动创建。请提供相对于工作区的相对路径。",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "要写入的文件路径，如 src/main.go",
				},
				"content": map[string]any{
					"type":        "string",
					"description": "要写入的完整文件内容",
				},
			},
			"required": []string{"path", "content"},
		},
	}
}

type writeFileParams struct {
	Path    string
	Content string
}

// Exec implements [BaseTool].
func (w *WriteFileTool) Exec(ctx context.Context, args json.RawMessage) (string, error) {
	var input writeFileParams
	err := json.Unmarshal(args, &input)
	if err != nil {
		return "", fmt.Errorf("参数解析失败, %w", err)
	}

	fullPath := filepath.Join(w.workDir, input.Path)

	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("创建父目录失败, %w", err)
	}

	err = os.WriteFile(fullPath, []byte(input.Content), 0644)
	if err != nil {
		return "", fmt.Errorf("写入文件失败，%w", err)
	}
	return fmt.Sprintf("成功将内容写入文件, %s\n", fullPath), nil
}

// Name implements [BaseTool].
func (w *WriteFileTool) Name() string {
	return "write_file"
}
