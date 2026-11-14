package provider

import (
	"context"
	"tiny-claw/internal/schema"
)

type LLMProvider interface {
	Generate(ctx context.Context, messages []*schema.Message, tools []*schema.ToolDefinition) (*schema.Message, error)
}
