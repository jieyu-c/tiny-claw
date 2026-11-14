package schema

import "encoding/json"

type Role string

const (
	RoleUser      Role = "user"
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
)

type Message struct {
	Role       Role        `json:"role"`
	Content    string      `json:"content"`
	ToolCalls  []*ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string      `json:"tool_call_id,omitempty"`
}

func UserMessage(input string) *Message {
	return &Message{
		Role:    RoleUser,
		Content: input,
	}
}

func MessageWithTools(content string, ToolCallID string) *Message {
	return &Message{
		Role:       RoleUser,
		Content:    content,
		ToolCallID: ToolCallID,
	}
}

func SystemMessage(input string) *Message {
	return &Message{
		Role:    RoleSystem,
		Content: input,
	}
}

func AssistantMessage(input string, tools []*ToolCall) *Message {
	return &Message{
		Role:      RoleAssistant,
		Content:   input,
		ToolCalls: tools,
	}
}

type ToolCall struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

type ToolResult struct {
	ToolCallID string `json:"tool_call_id"`
	Output     string `json:"output"`
	IsError    bool   `json:"is_error"`
}

type ToolDefinition struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	InputSchema any    `json:"input_schema"`
}
