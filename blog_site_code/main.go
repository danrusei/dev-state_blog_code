package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/danrusei/blog_site_code/site"
	"github.com/gorilla/mux"
)

const (
	//ContentDir is the location of the post folders
	ContentDir = "./content/"
	//ImageDir is the location of the post images
	ImageDir = "./images/"
)

//SiteConfig holds the configuration of the site
type SiteConfig struct {
	port    string
	sources []string

	landing  *site.LandingConfig
	posts    *site.PostConfig
	images   *site.ImgConfig
	template *template.Template
}

//NewSite pass config parameters from main function
func NewSite(port string, sources []string) *SiteConfig {
	return &SiteConfig{
		port:    port,
		sources: sources,
	}
}

func (s *SiteConfig) run() error {
	var err error
	templatePath := filepath.Join("templates", "template.html")
	paths := s.sources
	imgpath := ImageDir

	s.template, err = site.GetTemplate(templatePath)
	if err != nil {
		return err
	}

	s.images = site.NewImgConfig(imgpath)
	if err := s.images.Read(); err != nil {
		return err
	}

	s.posts = site.NewPostConfig(paths, s.template)
	if err := s.posts.Generate(); err != nil {
		return err
	}

	s.landing = site.NewLandingConfig(s.posts, s.template)
	if err := s.landing.Generate(); err != nil {
		return err
	}

	router := mux.NewRouter()
	router.HandleFunc("/", s.landingPage).Methods("GET")
	router.HandleFunc("", s.landingPage).Methods("GET")
	router.HandleFunc("/{key}", s.postHandler).Methods("GET")
	router.HandleFunc("/{key}/", s.postHandler).Methods("GET")
	router.HandleFunc("/images/{key}", s.imageHandler).Methods("GET")

	err = http.ListenAndServe(":8080", router)
	if err != nil {
		return err
	}

	return nil
}

func (s *SiteConfig) landingPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// The root page uses the "index" key
	if key == "" {
		key = "index"
	}

	// Try to get cache page
	page := s.landing.Get(key)
	if page == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Header.Get("If-None-Match") == page.Etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", page.Mime)
	w.Header().Set("Cache-Control", page.CacheControl)
	w.Header().Set("Etag", page.Etag)
	w.WriteHeader(http.StatusOK)
	w.Write(*page.Content)
}

func (s *SiteConfig) postHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// Try to get cache entry for post
	post := s.posts.Get(key)
	if post == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Header.Get("If-None-Match") == post.Etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "public, must-revalidate")
	w.Header().Set("Etag", post.Etag)
	w.WriteHeader(http.StatusOK)
	w.Write(post.Content)
}

func (s *SiteConfig) imageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	// Try to get cached entry for asset

	image := s.images.GetImage(key)
	if image == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if r.Header.Get("If-None-Match") == image.Etag {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	w.Header().Set("Content-Type", image.Mime)
	w.Header().Set("Cache-Control", "public, max-age=2419200")
	w.Header().Set("Etag", image.Etag)
	w.WriteHeader(http.StatusOK)
	w.Write(image.Imgfile)
}

func main() {
	var err error
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	dirs, err := site.GetContentFolders(ContentDir)

	newsite := NewSite(port, dirs)
	err = newsite.run()
	if err != nil {
		log.Fatal(err)
	}
}
