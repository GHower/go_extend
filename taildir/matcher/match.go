package matcher

type Matcher interface {
	// Match 匹配
	Match(str string) bool
}
