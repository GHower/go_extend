package matcher

import (
	"regexp"
	"strings"
)

// FileMatcher 后缀匹配
type FileMatcher struct {
}

func (fn *FileMatcher) Match(filename string) bool {
	return strings.HasSuffix(filename, ".log")
}

func NewFileSufMatcher() *FileMatcher {
	return &FileMatcher{}
}

// FileRegMatcher 正则匹配
type FileRegMatcher struct {
	regx *regexp.Regexp
}

func (fn *FileRegMatcher) Match(filename string) bool {
	return fn.regx.FindString(filename) != ""
}

func NewFileRegMatcher(regex string) (*FileRegMatcher, error) {
	regx, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}
	return &FileRegMatcher{regx}, nil
}
