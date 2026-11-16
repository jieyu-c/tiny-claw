package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"
	"tiny-claw/internal/schema"
)

type BashTool struct {
	workDir string
}

func NewBashTool(workDir string) *BashTool {
	return &BashTool{
		workDir: workDir,
	}
}

var _ BaseTool = new(BashTool)

// Description implements [BaseTool].
func (b *BashTool) Description() *schema.ToolDefinition {
	return &schema.ToolDefinition{
		Name:        b.Name(),
		Description: "在当前工作区执行任意的 bash 命令。支持链式命令(如 &&)。返回标准输出(stdout)和标准错误(stderr)。",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"command": map[string]any{
					"type":        "string",
					"description": "要执行的 bash 命令，例如: ls -la 或 go test ./...",
				},
			},
			"required": []string{"command"},
		},
	}
}

type bashArgs struct {
	Command string
}

// Exec implements [BaseTool].
func (b *BashTool) Exec(ctx context.Context, args json.RawMessage) (string, error) {
	var input bashArgs
	err := json.Unmarshal(args, &input)
	if err != nil {
		return "", fmt.Errorf("解析参数错误，%w", err)
	}
	log.Printf("执行命令, %s\n", input.Command)

	timeoutCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	commandName := "bash"
	commandList := []string{
		"-c",
		input.Command,
	}
	cmd := exec.CommandContext(timeoutCtx, commandName, commandList...)
	// 设置工作区
	cmd.Dir = b.workDir

	out, err := cmd.CombinedOutput()

	outputStr := string(out)

	if timeoutCtx.Err() == context.DeadlineExceeded {
		return fmt.Sprintf("%s\n[警告: 命令执行超时(30s)，已被系统强制终止。如果是启动常驻服务，请尝试将其转入后台。]", outputStr), nil
	}

	if err != nil {
		return fmt.Sprintf("执行报错: %v\n输出:\n%s", err, outputStr), nil
	}

	// 如果没有终端输出（比如仅仅执行了 mkdir），给模型一个明确的执行成功的反馈
	if outputStr == "" {
		return "命令执行成功，无终端输出。", nil
	}
	// 【驾驭底线 4】：长度截断保护 (防 OOM)
	const maxLen = 8000
	if len(outputStr) > maxLen {
		return fmt.Sprintf("%s\n\n...[终端输出过长，已截断至前 %d 字节]...", outputStr[:maxLen], maxLen), nil
	}
	return outputStr, nil
}

// Name implements [BaseTool].
func (b *BashTool) Name() string {
	return "bash_tool"
}
