package python

import "testing"

func TestParseZeroDivisionError(t *testing.T) {
	stderr := `Traceback (most recent call last):
  File "C:\Users\iwaku\pro\devduck\test.py", line 2, in <module>
    x = 10 / 0
        ~~~^~~
ZeroDivisionError: division by zero
`
	info, ok := Parse(stderr)

	if !ok {
		t.Fatal("expected traceback to be parsed")
	}

	if info.Source != "python" {
		t.Errorf("expected Source python, got %q", info.Source)
	}

	if info.Kind != "ZeroDivisionError" {
		t.Errorf(
			"expected Kind ZeroDivisionError, got %q",
			info.Kind,
		)
	}

	if info.Message != "division by zero" {
		t.Errorf(
			"expected Message division by zero, got %q",
			info.Message,
		)
	}

	expectedFile := `C:\Users\iwaku\pro\devduck\test.py`

	if info.File != expectedFile {
		t.Errorf(
			"expected File %q, got %q",
			expectedFile,
			info.File,
		)
	}

	if info.Line != 2 {
		t.Errorf(
			"expected Line 2, got %d",
			info.Line,
		)
	}

}
