package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"tiny-claw/internal/schema"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/openai/openai-go/v3/shared"
)

type OpenAI struct {
	client *openai.Client
	model  string
}

func NewOpenAI(model string) *OpenAI {
	apiKey := os.Getenv("OPENAI_API_KEY")
	baseUrl := os.Getenv("OPENAI_BASE_URL")
	client := openai.NewClient(option.WithBaseURL(baseUrl), option.WithAPIKey(apiKey))
	return &OpenAI{
		client: &client,
		model:  model,
	}
}

var _ LLMProvider = new(OpenAI)

func (o *OpenAI) Generate(ctx context.Context, messages []*schema.Message, tools []*schema.ToolDefinition) (*schema.Message, error) {
	chatMessages := o.buildMessages(messages)
	openaiTools := o.buildTools(tools)

	// 构建open ai参数
	params := openai.ChatCompletionNewParams{
		Model:    o.model,
		Messages: chatMessages,
	}

	// 存在可用工具
	if len(tools) != 0 {
		params.Tools = openaiTools
	}

	resp, err := o.client.Chat.Completions.New(ctx, params)
	if err != nil {
		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("响应为空")
	}

	choice := resp.Choices[0].Message

	var resultTools []*schema.ToolCall

	openaiToolCalls := choice.ToolCalls
	for _, tool := range openaiToolCalls {
		if tool.Type == "function" {
			resultTools = append(resultTools, &schema.ToolCall{
				ID:        tool.ID,
				Name:      tool.Function.Name,
				Arguments: json.RawMessage(tool.Function.Arguments),
			})
		}
	}

	resultMsg := schema.AssistantMessage(choice.Content, resultTools)

	return resultMsg, nil
}

// buildMessages 构建open ai格式的messages
func (o *OpenAI) buildMessages(messages []*schema.Message) (chatMessages []openai.ChatCompletionMessageParamUnion) {
	for _, msg := range messages {
		switch msg.Role {
		case schema.RoleUser:
			if msg.ToolCallID != "" {
				chatMessages = append(chatMessages, openai.ToolMessage(msg.Content, msg.ToolCallID))
			} else {
				chatMessages = append(chatMessages, openai.UserMessage(msg.Content))
			}
		case schema.RoleAssistant:
			astParam := openai.ChatCompletionAssistantMessageParam{}
			if msg.Content != "" {
				astParam.Content = openai.ChatCompletionAssistantMessageParamContentUnion{
					OfString: openai.String(msg.Content),
				}
			}
			// 存在可调用的工具
			if len(msg.ToolCalls) != 0 {
				var toolCalls []openai.ChatCompletionMessageToolCallUnionParam
				for _, tool := range msg.ToolCalls {
					toolCalls = append(toolCalls, openai.ChatCompletionMessageToolCallUnionParam{
						OfFunction: &openai.ChatCompletionMessageFunctionToolCallParam{
							ID:   tool.ID,
							Type: "function",
							Function: openai.ChatCompletionMessageFunctionToolCallFunctionParam{
								Name:      tool.Name,
								Arguments: string(tool.Arguments),
							},
						},
					})
				}
				astParam.ToolCalls = toolCalls
			}
			chatMessages = append(chatMessages, openai.ChatCompletionMessageParamUnion{
				OfAssistant: &astParam,
			})
		case schema.RoleSystem:
			chatMessages = append(chatMessages, openai.SystemMessage(msg.Content))
		}
	}
	return
}

// buildTools 构建工具
func (o *OpenAI) buildTools(tools []*schema.ToolDefinition) (openaiTools []openai.ChatCompletionToolUnionParam) {
	for _, tool := range tools {
		var params shared.FunctionParameters
		if m, ok := tool.InputSchema.(map[string]any); ok {
			params = m
		} else {
			// fallback：JSON 往返序列化
			b, _ := json.Marshal(tool.InputSchema)
			_ = json.Unmarshal(b, &params)
		}
		openaiTools = append(openaiTools, openai.ChatCompletionFunctionTool(
			shared.FunctionDefinitionParam{
				Name:        tool.Name,
				Description: openai.String(tool.Description),
				Parameters:  params,
			},
		))
	}
	return
}

func (o *OpenAI) buildResponse() *schema.Message {
	return nil
}
