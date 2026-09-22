package teaching

import "github.com/duck-ahiru-Z/DevDuck/internal/model"

type Session struct {
	explanation model.Explanation
	hintIndex   int
}

func NewSession(explanation model.Explanation) *Session {
	return &Session{
		explanation: explanation,
		hintIndex:   0,
	}
}
func (s *Session) HasNextHint() bool {
	return s.hintIndex < len(s.explanation.Hints)
}

func (s *Session) NextHint() (string, bool) {
	if !s.HasNextHint() {
		return "", false
	}

	hint := s.explanation.Hints[s.hintIndex]
	s.hintIndex++

	return hint, true
}
