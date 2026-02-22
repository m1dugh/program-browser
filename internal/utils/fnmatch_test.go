package utils_test

import (
	"testing"

	"github.com/m1dugh/program-browser/internal/utils"
)

func TestFnMatchBasic(t *testing.T) {
	str := "hello"
	pattern := "hello"

	if !utils.FnMatch(pattern, str) {
		t.Errorf("Expected pattern '%s', to match string '%s'", pattern, str)
	}
}

func TestFnMatchBasicStar(t *testing.T) {
	str := "hello"
	pattern := "*"

	if !utils.FnMatch(pattern, str) {
		t.Errorf("Expected pattern '%s', to match string '%s'", pattern, str)
	}
}

func TestFnMatchTrailingStar(t *testing.T) {
	str := "hello"
	pattern := "hello*"

	if !utils.FnMatch(pattern, str) {
		t.Errorf("Expected pattern '%s', to match string '%s'", pattern, str)
	}
}

func TestFnMatchMiddleStar(t *testing.T) {
	str := "hello"
	pattern := "h*o"

	if !utils.FnMatch(pattern, str) {
		t.Errorf("Expected pattern '%s', to match string '%s'", pattern, str)
	}
}

func TestFnMatchMiddleStarInvalid(t *testing.T) {
	str := "hello"
	pattern := "h*e"

	if utils.FnMatch(pattern, str) {
		t.Errorf("Expected pattern '%s', to match string '%s'", pattern, str)
	}
}

func TestFnMatchSingle(t *testing.T) {
	str := "hello"
	pattern := "he??o"

	if !utils.FnMatch(pattern, str) {
		t.Errorf("Expected pattern '%s', to match string '%s'", pattern, str)
	}
}

func TestFnMatchSingleInvalidEnd(t *testing.T) {
	str := "hello"
	pattern := "hello?"

	if utils.FnMatch(pattern, str) {
		t.Errorf("Expected pattern '%s', to match string '%s'", pattern, str)
	}
}

func TestFnMatchDiffLenPattern(t *testing.T) {
	str := "hello"
	pattern := "hel"

	if utils.FnMatch(pattern, str) {
		t.Errorf("Expected pattern '%s', to match string '%s'", pattern, str)
	}
}

func TestFnMatchDiffLenStr(t *testing.T) {
	str := "hel"
	pattern := "hello"

	if utils.FnMatch(pattern, str) {
		t.Errorf("Expected pattern '%s', to match string '%s'", pattern, str)
	}
}
