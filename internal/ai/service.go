package ai

import (
	"context"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
	"github.com/duck-ahiru-Z/DevDuck/internal/redact"
	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
)

type Service struct {
	provider Provider
	redactor *redact.Redactor
}

func NewService(
	provider Provider,
	redactor *redact.Redactor,
) *Service {
	return &Service{
		provider: provider,
		redactor: redactor,
	}
}

func (s *Service) Explain(
	ctx context.Context,
	err model.ErrorInfo,
	level teaching.Level,
	contextText string,
) (model.Explanation, error) {

	safeRaw := s.redactor.Sanitize(
		err.Raw,
	)

	safeContext := s.redactor.Sanitize(
		contextText,
	)

	safeError := err
	safeError.Raw = safeRaw

	request := Request{
		Error:   safeError,
		Level:   level,
		Context: safeContext,
	}

	return s.provider.Explain(
		ctx,
		request,
	)
}
