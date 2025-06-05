package health

import "github.com/harrisonju123/mcp-agent-poc/router"

func Wrap(t router.Tool) *Breaker {
	return &Breaker{
		callable: t,
		r:        &Recorder{},
	}
}
