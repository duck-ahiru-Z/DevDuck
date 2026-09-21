package model

type ErrorInfo struct {
	Source  string
	Kind    string
	Message string
	File    string
	Line    int
	Raw     string
}
