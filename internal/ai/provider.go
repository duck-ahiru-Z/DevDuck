package ai

import (
	"context"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
)

// Request はAIに渡す共通リクエスト。
type Request struct {
	Error   model.ErrorInfo
	Level   teaching.Level
	Context string
}

// Provider はAIサービスが実装する共通インターフェース。
type Provider interface {
	Name() string

	Explain(
		ctx context.Context,
		request Request,
	) (model.Explanation, error)
}
