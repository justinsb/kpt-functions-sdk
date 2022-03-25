package fn

import "context"

type Context struct {
	context.Context

	Results []*Result
}
