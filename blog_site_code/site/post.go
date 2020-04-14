package site

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/russross/blackfriday"
)

//Meta holds the post metadata
type Meta struct {
	Title   string   `json:"title"`
	Short   string   `json:"short"`
	Date    string   `json:"date"`
	Tags    []string `json:"tags"`
	Created time.Time
}

//Post holds post data
type Post struct {
	Name    string
	Meta    *Meta
	Content []byte
	Etag    string
}

//PostConfig holds post config
type PostConfig struct {
	OrderedPosts []*Post
	Paths        []string
	Template     *template.Template
	Cache        *Cache
}

//NewPostConfig generate new post configuration
func NewPostConfig(paths []string, template *template.Template) *PostConfig {
	return &PostConfig{
		Paths:    paths,
		Template: template,
		Cache:    NewCache(),
	}
}

//Generate the posts from config
func (p *PostConfig) Generate() error {
	templ := p.Template

	var posts []*Post
	for _, path := range p.Paths {
		fmt.Printf("\tGenerating Post : %s...\n", path)
		post, err := newPost(path, templ)
		if err != nil {
			return err
		}
		p.Cache.Set(post.Name, post)
		fmt.Printf("\tFinished generating Post: %s...\n", post.Meta.Title)
	}
	values := p.Cache.GetValues()

	for _, post := range values {
		posts = append(posts, post.(*Post))
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Meta.Created.After(posts[j].Meta.Created)
	})

	p.OrderedPosts = posts

	return nil
}

// Get the post content from cache
func (p *PostConfig) Get(key string) *Post {
	item := p.Cache.Get(key)
	if item == nil {
		return nil
	}

	return item.(*Post)
}

func newPost(path string, templ *template.Template) (*Post, error) {
	meta, err := getMeta(path)
	if err != nil {
		return nil, err
	}

	postbodyhtml, err := createHTML(path)
	if err != nil {
		return nil, err
	}

	name := filepath.Base(path)

	posthtml, err := newIndexhtml(meta.Title, postbodyhtml, templ)
	if err != nil {
		return nil, err
	}

	return &Post{
		Name:    name,
		Meta:    meta,
		Content: posthtml,
		Etag:    getEtag(&posthtml),
	}, nil
}

func createHTML(path string) ([]byte, error) {
	filePath := filepath.Join(path, "post.md")
	input, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %v", filePath, err)
	}
	mdtohtml := blackfriday.MarkdownCommon(input)
	replaced, err := replaceCodeParts(mdtohtml)
	if err != nil {
		return nil, fmt.Errorf("error during syntax highlighting of %s: %v", filePath, err)
	}
	return []byte(replaced), nil
}

func replaceCodeParts(htmlFile []byte) (string, error) {
	byteReader := bytes.NewReader(htmlFile)
	doc, err := goquery.NewDocumentFromReader(byteReader)
	if err != nil {
		return "", fmt.Errorf("error while parsing html: %v", err)
	}
	// find code-parts via css selector and replace them with highlighted versions
	doc.Find("code[class*=\"language-\"]").Each(func(i int, s *goquery.Selection) {
		class, _ := s.Attr("class")
		lang := strings.TrimPrefix(class, "language-")
		oldCode := s.Text()
		lexer := lexers.Get(lang)
		formatter := html.New(html.WithClasses())
		iterator, err := lexer.Tokenise(nil, string(oldCode))
		if err != nil {
			fmt.Printf("ERROR during syntax highlighting, %v", err)
		}
		b := bytes.Buffer{}
		buf := bufio.NewWriter(&b)
		err = formatter.Format(buf, styles.GitHub, iterator)
		if err != nil {
			fmt.Printf("ERROR during syntax highlighting, %v", err)
		}
		buf.Flush()
		s.SetHtml(b.String())
	})
	new, err := doc.Html()
	if err != nil {
		return "", fmt.Errorf("error while generating html: %v", err)
	}
	// replace unnecessarily added html tags
	new = strings.Replace(new, "<html><head></head><body>", "", 1)
	new = strings.Replace(new, "</body></html>", "", 1)
	return new, nil
}
