package factory

import (
	"fmt"
	"os"

	"github.com/duck-ahiru-Z/DevDuck/internal/ai"
	"github.com/duck-ahiru-Z/DevDuck/internal/ai/providers/gemini"
	"github.com/duck-ahiru-Z/DevDuck/internal/credential"
)

// NewProviderFromEnv creates the configured provider without exposing provider
// construction details to the command layer.
const credentialService = "DevDuck"

func NewProviderFromEnv(stores ...credential.Store) (ai.Provider, error) {
	name := os.Getenv("DEVDUCK_AI_PROVIDER")
	if name == "" {
		name = "gemini"
	}

	switch name {
	case "gemini":
		var apiKey string
		if len(stores) > 0 {
			stored, err := stores[0].Get(credentialService, "gemini")
			if err == nil {
				apiKey = stored
			}
		}
		if apiKey == "" {
			apiKey = os.Getenv("GEMINI_API_KEY")
		}
		if apiKey == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY が設定されていません")
		}
		modelName := os.Getenv("DEVDUCK_GEMINI_MODEL")
		if modelName == "" {
			modelName = "gemini-3.5-flash-lite"
		}
		return gemini.New(apiKey, modelName), nil
	default:
		return nil, fmt.Errorf("未対応のAI providerです: %s", name)
	}
}
