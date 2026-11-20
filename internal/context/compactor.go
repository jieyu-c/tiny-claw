package ctxpkg

import (
	"fmt"
	"log"
	"tiny-claw/internal/schema"
)

type Compactor struct {
	MaxChars       int
	RetainLastMsgs int
}

func NewCompactor(maxChars int, retainLastMsgs int) *Compactor {
	return &Compactor{maxChars, retainLastMsgs}
}

func (c *Compactor) Compact(msgs []*schema.Message) []*schema.Message {
	currentLength := c.estimateLength(msgs)
	if currentLength < c.MaxChars {
		return msgs
	}

	log.Printf("[Compactor] ⚠️ 内存告警：当前上下文长度 (%d 字符) 超过阈值 (%d)，触发压缩清理...\n", currentLength, c.MaxChars)
	var compacted []*schema.Message
	msgCount := len(msgs)

	protectStartIndex := max(msgCount-c.RetainLastMsgs, 0)

	for index, msg := range msgs {
		if msg.Role == schema.RoleSystem {
			compacted = append(compacted, msg)
			continue
		}

		newMsg := *msg
		isInWorkingMemory := index >= protectStartIndex
		if msg.Role == schema.RoleUser && msg.ToolCallID != "" {
			if !isInWorkingMemory {
				// 【第一道防线：远期历史】如果是早期对话，执行无情替换 (Full Masking)
				if len(msg.Content) > 200 {
					newMsg.Content = fmt.Sprintf("...[为了节省内存，早期的工具输出已被系统强制清理。原始长度: %d 字节]...", len(msg.Content))
				}
			} else {
				// 【第二道防线：短期记忆】即使处于近期保护区，只要单条内容过大，也必须截断防 OOM (Head-Tail Truncation)
				// 我们保留前 500 字符和后 500 字符（掐头去尾法，大模型通常只需要看开头报错和结尾总结）
				const maxKeep = 1000
				if len(msg.Content) > maxKeep {
					head := msg.Content[:500]
					tail := msg.Content[len(msg.Content)-500:]
					newMsg.Content = fmt.Sprintf("%s\n\n...[内容过长，中间 %d 字节已被系统截断]...\n\n%s", head, len(msg.Content)-maxKeep, tail)
				}
			}
		} else if msg.Role == schema.RoleAssistant && msg.Content != "" {
			if !isInWorkingMemory && len(msg.Content) > 200 {
				newMsg.Content = "...[早期推理思考过程已折叠]..."
			}
		}
		compacted = append(compacted, &newMsg)
	}
	newLength := c.estimateLength(compacted)
	log.Printf("[Compactor] ✅ 压缩完成。上下文长度从 %d 降至 %d 字符。\n", currentLength, newLength)
	return compacted
}

// estimateLength 粗略计算当前上下文的总字符长度
func (c *Compactor) estimateLength(msgs []*schema.Message) int {
	length := 0
	for _, msg := range msgs {
		length += len(msg.Content)
		for _, tc := range msg.ToolCalls {
			length += len(tc.Name) + len(tc.Arguments)
		}
	}
	return length
}
