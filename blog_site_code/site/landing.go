package site

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
)

const indexKey = "index"

//LandingPage holds the content of the main/index page
type LandingPage struct {
	Content      *[]byte
	Mime         string
	Etag         string
	CacheControl string
}

//PostsListing holds the data for post listed in index page
type PostsListing struct {
	Title string
	Date  string
	Short string
	Link  string
	Tags  []string
}

//LandingConfig holds Landing page config
type LandingConfig struct {
	posts     *PostConfig
	templates *template.Template
	pageTitle string
	cache     *Cache
}

//NewLandingConfig initiate the Landing page config
func NewLandingConfig(posts *PostConfig, templates *template.Template) *LandingConfig {
	return &LandingConfig{
		posts:     posts,
		templates: templates,
		cache:     NewCache(),
	}
}

// Generate the landing page
func (l *LandingConfig) Generate() error {
	shortTemplatePath := filepath.Join("templates", "short.html")
	posts := l.posts.OrderedPosts
	templ := l.templates
	pageTitle := l.pageTitle
	short, err := GetTemplate(shortTemplatePath)
	if err != nil {
		return err
	}

	var postBlocks []string
	for _, post := range posts {
		meta := post.Meta
		link := fmt.Sprintf("/%s/", post.Name)
		ld := PostsListing{
			Title: meta.Title,
			Date:  meta.Date,
			Short: meta.Short,
			Link:  link,
			Tags:  meta.Tags,
		}
		block := bytes.Buffer{}
		if err := short.Execute(&block, ld); err != nil {
			//changed from return Errorf to Printf
			fmt.Printf("error executing template %s: %v\n", shortTemplatePath, err)
		}
		postBlocks = append(postBlocks, block.String())
	}
	htmlBlocks := []byte(strings.Join(postBlocks, "<br />"))

	landinghtml, err := newIndexhtml(pageTitle, htmlBlocks, templ)
	if err != nil {
		return fmt.Errorf("could not initiate newIndexhtml for landing page: %v", err)
	}

	l.cache.Set(indexKey, &LandingPage{
		Content:      &landinghtml,
		Etag:         getEtag(&landinghtml),
		Mime:         "text/html; charset=utf-8",
		CacheControl: "public, must-revalidate",
	})

	return nil
}

//Get the page from cache
func (l *LandingConfig) Get(key string) *LandingPage {
	item := l.cache.Get(key)
	if item == nil {
		return nil
	}

	return item.(*LandingPage)
}
