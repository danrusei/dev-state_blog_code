package actions

import (
	"github.com/danrusei/blog_site_buf/models"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop"
	"github.com/pkg/errors"
)

// PostsIndex default implementation.
func PostsIndex(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	posts := &models.Posts{}
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())
	// Retrieve all Posts from the DB
	if err := q.All(posts); err != nil {
		return errors.WithStack(err)
	}
	// Make posts available inside the html template
	c.Set("posts", posts)
	// Add the paginator to the context so it can be used in the template.
	c.Set("pagination", q.Paginator)
	return c.Render(200, r.HTML("posts/index.html"))
}

// PostsDetail default implementation.
func PostsDetail(c buffalo.Context) error {
	tx := c.Value("tx").(*pop.Connection)
	post := &models.Post{}
	if err := tx.Find(post, c.Param("pid")); err != nil {
		return c.Error(404, err)
	}
	c.Set("post", post)
	comment := &models.Comment{}
	c.Set("comment", comment)
	comments := models.Comments{}
	if err := tx.BelongsTo(post).All(&comments); err != nil {
		return errors.WithStack(err)
	}
	for i := 0; i < len(comments); i++ {
		u := models.User{}
		if err := tx.Find(&u, comments[i].AuthorID); err != nil {
			return c.Error(404, err)
		}
		comments[i].Author = u
	}
	c.Set("comments", comments)
	return c.Render(200, r.HTML("posts/detail"))
}
