package aiverr

import (
	"errors"
	"strings"
	"testing"
)

func TestNewPlatformError(t *testing.T) {
	err := New(CodeAPIError, "something went wrong")
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if err.Code != CodeAPIError {
		t.Errorf("expected code %q, got %q", CodeAPIError, err.Code)
	}
	if err.Message != "something went wrong" {
		t.Errorf("expected message %q, got %q", "something went wrong", err.Message)
	}
	if err.Details != nil {
		t.Error("expected nil Details for plain New")
	}
}

func TestErrorCodes(t *testing.T) {
	codes := map[string]string{
		"not_authenticated": CodeNotAuthenticated,
		"rate_limited":      CodeRateLimited,
		"api_error":         CodeAPIError,
		"network_error":     CodeNetworkError,
		"parse_error":       CodeParseError,
		"not_found":         CodeNotFound,
		"forbidden":         CodeForbidden,
		"invalid_input":     CodeInvalidInput,
	}
	for expected, actual := range codes {
		if actual != expected {
			t.Errorf("expected code %q, got %q", expected, actual)
		}
	}
}

func TestWithPlatform(t *testing.T) {
	err := New(CodeAPIError, "msg").WithPlatform("test-platform")
	if err.Details == nil {
		t.Fatal("expected non-nil Details after WithPlatform")
	}
	if err.Details.Platform != "test-platform" {
		t.Errorf("expected platform %q, got %q", "test-platform", err.Details.Platform)
	}
}

func TestNotAuthenticated(t *testing.T) {
	err := NotAuthenticated("test-platform", "need login")
	if err.Code != CodeNotAuthenticated {
		t.Errorf("expected code %q, got %q", CodeNotAuthenticated, err.Code)
	}
	if err.Details.Platform != "test-platform" {
		t.Errorf("expected platform %q, got %q", "test-platform", err.Details.Platform)
	}
	if !strings.Contains(err.Details.Suggestion, "aiview") {
		t.Errorf("expected suggestion to contain 'aiview', got %q", err.Details.Suggestion)
	}
}

func TestAPIError(t *testing.T) {
	err := APIError("test-platform", "api failed")
	if err.Code != CodeAPIError {
		t.Errorf("expected code %q, got %q", CodeAPIError, err.Code)
	}
	if err.Details.Platform != "test-platform" {
		t.Errorf("expected platform %q, got %q", "test-platform", err.Details.Platform)
	}
}

func TestNetworkError(t *testing.T) {
	err := NetworkError("test-platform", "network timeout")
	if err.Code != CodeNetworkError {
		t.Errorf("expected code %q, got %q", CodeNetworkError, err.Code)
	}
	if err.Details.Platform != "test-platform" {
		t.Errorf("expected platform %q, got %q", "test-platform", err.Details.Platform)
	}
}

func TestParseError(t *testing.T) {
	err := ParseError("test-platform", "invalid json")
	if err.Code != CodeParseError {
		t.Errorf("expected code %q, got %q", CodeParseError, err.Code)
	}
	if err.Details.Platform != "test-platform" {
		t.Errorf("expected platform %q, got %q", "test-platform", err.Details.Platform)
	}
}

func TestNotFound(t *testing.T) {
	err := NotFound("test-platform", "resource not found")
	if err.Code != CodeNotFound {
		t.Errorf("expected code %q, got %q", CodeNotFound, err.Code)
	}
	if err.Details.Platform != "test-platform" {
		t.Errorf("expected platform %q, got %q", "test-platform", err.Details.Platform)
	}
}

func TestInvalidInput(t *testing.T) {
	err := InvalidInput("test-platform", "bad input")
	if err.Code != CodeInvalidInput {
		t.Errorf("expected code %q, got %q", CodeInvalidInput, err.Code)
	}
	if err.Details.Platform != "test-platform" {
		t.Errorf("expected platform %q, got %q", "test-platform", err.Details.Platform)
	}
}

func TestWrap(t *testing.T) {
	// wrap a standard error
	stdErr := errors.New("standard error")
	pe := Wrap(stdErr, CodeNetworkError, "test-platform")
	if pe == nil {
		t.Fatal("expected non-nil result for non-nil error")
	}
	if pe.Code != CodeNetworkError {
		t.Errorf("expected code %q, got %q", CodeNetworkError, pe.Code)
	}
	if pe.Message != "standard error" {
		t.Errorf("expected message %q, got %q", "standard error", pe.Message)
	}

	// wrap nil should return nil
	nilResult := Wrap(nil, CodeAPIError, "p")
	if nilResult != nil {
		t.Error("expected nil when wrapping nil error")
	}

	// wrapping an already PlatformError should return it as-is
	existing := New(CodeAPIError, "already wrapped").WithPlatform("x")
	result := Wrap(existing, CodeNetworkError, "y")
	if result != existing {
		t.Error("expected same PlatformError pointer")
	}
	if result.Code != CodeAPIError {
		t.Errorf("expected original code %q, got %q", CodeAPIError, result.Code)
	}
}

func TestRateLimited(t *testing.T) {
	err := RateLimited("test-platform", "too many requests")
	if err.Code != CodeRateLimited {
		t.Errorf("expected code %q, got %q", CodeRateLimited, err.Code)
	}
	if err.Details.Platform != "test-platform" {
		t.Errorf("expected platform %q, got %q", "test-platform", err.Details.Platform)
	}
	if err.Details.Suggestion != "Please wait a moment and try again" {
		t.Errorf("unexpected suggestion: %q", err.Details.Suggestion)
	}
}

func TestError(t *testing.T) {
	// without suggestion
	err1 := New(CodeAPIError, "basic error")
	errStr1 := err1.Error()
	expected1 := "api_error: basic error"
	if errStr1 != expected1 {
		t.Errorf("expected %q, got %q", expected1, errStr1)
	}

	// with suggestion
	err2 := New(CodeAPIError, "basic error").WithSuggestion("try again")
	errStr2 := err2.Error()
	if !strings.Contains(errStr2, "try again") {
		t.Errorf("expected output to contain suggestion, got %q", errStr2)
	}

	// without details
	err3 := &PlatformError{Code: CodeNotFound, Message: "gone"}
	if err3.Error() != "not_found: gone" {
		t.Errorf("unexpected error string: %q", err3.Error())
	}
}