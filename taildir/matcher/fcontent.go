package matcher

import (
	"regexp"
	"strings"
)

// ContentMatcher 内容匹配
type ContentMatcher struct {
	keyword string
}

func NewContentMatcher() *ContentMatcher {
	return &ContentMatcher{}
}
func (r *ContentMatcher) Match(line string) bool {
	return strings.Contains(line, r.keyword)
}

// ContentRegMatcher 内容正则匹配
type ContentRegMatcher struct {
	regx *regexp.Regexp
}

func NewRegexpMatcher(regexP string) (*ContentRegMatcher, error) {
	regx, err := regexp.Compile(regexP)
	if err != nil {
		return nil, err
	}
	return &ContentRegMatcher{regx: regx}, nil
}
func (r *ContentRegMatcher) Match(line string) bool {
	return r.regx.FindString(line) != ""
}
