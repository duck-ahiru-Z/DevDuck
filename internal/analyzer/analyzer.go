package analyzer

import "github.com/duck-ahiru-Z/DevDuck/internal/model"

type Analyzer interface {
	Explain(err model.ErrorInfo) (model.Explanation, bool)
}
