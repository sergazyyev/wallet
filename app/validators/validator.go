package validators

//JSONRequestValidator validator interface
type JSONRequestValidator interface {
	Rules() map[string][]string
	Messages() map[string][]string
}
