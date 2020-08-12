package grifts

import (
	"github.com/danrusei/blog_site_buf/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
