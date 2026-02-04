package errors

import (
	"errors"
	"testing"
)

func TestErrorString(t *testing.T) {
	err := New(CodeNotFound, "project not found")
	expected := "[NOT_FOUND] project not found"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestErrorWrap(t *testing.T) {
	cause := errors.New("connection refused")
	err := Wrap(CodeDBError, "database error", cause)

	if !errors.Is(err, cause) {
		t.Error("Wrapped error should contain cause")
	}

	expected := "[DB_ERROR] database error: connection refused"
	if err.Error() != expected {
		t.Errorf("Error() = %q, want %q", err.Error(), expected)
	}
}

func TestConvenienceConstructors(t *testing.T) {
	tests := []struct {
		name string
		err  *Error
		code Code
	}{
		{"NotFound", NotFound("project"), CodeNotFound},
		{"AlreadyExists", AlreadyExists("project"), CodeAlreadyExists},
		{"InvalidInput", InvalidInput("bad input"), CodeInvalidInput},
		{"Timeout", Timeout("timed out"), CodeTimeout},
		{"Cancelled", Cancelled("cancelled"), CodeCancelled},
		{"ConfigError", ConfigError("bad config"), CodeConfig},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.code {
				t.Errorf("Code = %v, want %v", tt.err.Code, tt.code)
			}
		})
	}
}

func TestIsCode(t *testing.T) {
	err := NotFound("project")

	if !IsCode(err, CodeNotFound) {
		t.Error("IsCode should return true for matching code")
	}

	if IsCode(err, CodeDBError) {
		t.Error("IsCode should return false for non-matching code")
	}

	if IsCode(errors.New("plain error"), CodeNotFound) {
		t.Error("IsCode should return false for non-Error type")
	}
}

func TestGetCode(t *testing.T) {
	err := DBError(errors.New("db failed"))
	if GetCode(err) != CodeDBError {
		t.Errorf("GetCode = %v, want %v", GetCode(err), CodeDBError)
	}

	if GetCode(errors.New("plain error")) != CodeUnknown {
		t.Error("GetCode should return CodeUnknown for non-Error type")
	}
}
