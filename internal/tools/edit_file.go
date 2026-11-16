package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"tiny-claw/internal/schema"
)

type EditFileTool struct {
	workDir string
}

var _ BaseTool = new(EditFileTool)

// Description implements [BaseTool].
func (e *EditFileTool) Description() *schema.ToolDefinition {
	return &schema.ToolDefinition{
		Name:        e.Name(),
		Description: "对现有文件进行局部的字符串替换。这比重写整个文件更安全、更快速。请提供足够的 old_text 上下文以确保匹配的唯一性。",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"path": map[string]any{
					"type":        "string",
					"description": "要修改的文件路径",
				},
				"old_text": map[string]any{
					"type":        "string",
					"description": "文件中原有的文本。必须包含足够的上下文（建议上下各多包含几行），以确保在文件中的唯一性。",
				},
				"new_text": map[string]any{
					"type":        "string",
					"description": "要替换成的新文本",
				},
			},
			"required": []string{"path", "old_text", "new_text"},
		},
	}

}

type editFileArgs struct {
	Path    string `json:"path"`
	OldText string `json:"old_text"`
	NewText string `json:"new_text"`
}

// Exec implements [BaseTool].
func (e *EditFileTool) Exec(ctx context.Context, args json.RawMessage) (string, error) {
	return "", nil
}

// Name implements [BaseTool].
func (e *EditFileTool) Name() string {
	return "edit_file"
}

func fuzzyReplace(originalContent, oldText, newText string) (string, error) {
	count := strings.Count(originalContent, oldText)

	
	if count == 1 {
		return strings.Replace(originalContent, oldText, newText, 1), nil
	}

	if count > 1 {
		return "", fmt.Errorf("old_text 匹配到了 %d 处，请提供更多的上下文代码以确保唯一性", count)
	}


}
