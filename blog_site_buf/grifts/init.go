package grifts

import (
	"github.com/Danr17/blog_site_buf/actions"
	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
