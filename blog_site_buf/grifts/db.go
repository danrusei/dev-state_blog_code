package grifts

import (
	"io/ioutil"

	"github.com/danrusei/blog_site_buf/models"
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		path := "content/"

		files, err := ioutil.ReadDir(path)
		if err != nil {
			return err
		}

		for _, file := range files {
			name := file.Name()
			content, err := ioutil.ReadFile(path + name)
			if err != nil {
				return err
			}
			u := &models.Post{
				Title:   name,
				Content: string(content),
			}
			err = models.DB.Create(u)
		}
		return err
	})

})
