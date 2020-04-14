package site

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/styles"
)

//IndexData holds Index page data used by both posts & landingpage
type IndexData struct {
	HTMLTitle       string
	PageTitle       string
	Content         template.HTML
	Year            int
	Name            string
	CanonicalLink   string
	MetaDescription string
	HighlightCSS    template.CSS
}

func newIndexhtml(title string, postbodyhtml []byte, templ *template.Template) ([]byte, error) {
	var err error

	hlbuf := bytes.Buffer{}
	hlw := bufio.NewWriter(&hlbuf)
	formatter := html.New(html.WithClasses())
	formatter.WriteCSS(hlw, styles.MonokaiLight)
	hlw.Flush()

	buf := &bytes.Buffer{}
	err = templ.ExecuteTemplate(buf, "template.html", &IndexData{
		HTMLTitle:       "Dev State",
		PageTitle:       title,
		Content:         template.HTML(postbodyhtml),
		Year:            time.Now().Year(),
		Name:            "Dan Rusei",
		CanonicalLink:   "http://www.dev_state.com",
		MetaDescription: "None",
		HighlightCSS:    template.CSS(hlbuf.String()),
	})
	if err != nil {
		return nil, fmt.Errorf("error executing template %v", err)
	}

	content := buf.Bytes()

	return content, nil
}
