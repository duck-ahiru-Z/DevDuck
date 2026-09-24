package ai

import (
	"context"
	"strings"
	"testing"

	"github.com/duck-ahiru-Z/DevDuck/internal/model"
	"github.com/duck-ahiru-Z/DevDuck/internal/redact"
	"github.com/duck-ahiru-Z/DevDuck/internal/teaching"
)

type fakeProvider struct {
	lastRequest Request
}

func (f *fakeProvider) Name() string {
	return "fake"
}

func (f *fakeProvider) Explain(
	ctx context.Context,
	request Request,
) (model.Explanation, error) {

	f.lastRequest = request

	return model.Explanation{
		Summary: "fake explanation",
		Hints: []string{
			"fake hint",
		},
	}, nil
}

func TestServiceRedactsSecrets(t *testing.T) {
	provider := &fakeProvider{}

	service := NewService(
		provider,
		redact.New(),
	)

	errInfo := model.ErrorInfo{
		Source: "python",
		Kind:   "TestError",
		Raw:    "API_KEY=super-secret-key",
	}

	_, err := service.Explain(
		context.Background(),
		errInfo,
		teaching.Beginner,
		"password=my-password",
	)

	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(
		provider.lastRequest.Error.Raw,
		"super-secret-key",
	) {
		t.Fatal(
			"secret remained in Error.Raw",
		)
	}

	if strings.Contains(
		provider.lastRequest.Context,
		"my-password",
	) {
		t.Fatal(
			"secret remained in Context",
		)
	}
}
