package site

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

//GetContentFolders retrive the dir names from content folder
func GetContentFolders(path string) ([]string, error) {
	var result []string
	dir, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error accessing directory %s: %v", path, err)
	}
	defer dir.Close()
	files, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("error reading contents of directory %s: %v", path, err)
	}
	for _, file := range files {
		if file.IsDir() && file.Name()[0] != '.' {
			result = append(result, filepath.Join(path, file.Name()))
		}
	}
	return result, nil
}

//GetTemplate read the template file
func GetTemplate(path string) (*template.Template, error) {
	t, err := template.ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("error reading template %s: %v", path, err)
	}
	return t, nil
}

func getMeta(path string) (*Meta, error) {
	filePath := filepath.Join(path, "meta.json")
	metaraw, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error while reading file %s: %v", filePath, err)
	}
	meta := Meta{}
	if err := json.Unmarshal(metaraw, &meta); err != nil {
		return nil, fmt.Errorf("error unmarshaling meta file %s, %v", filePath, err)

	}

	parsed, err := time.Parse("02.01.2006", meta.Date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date %s: %v", filePath, err)
	}
	meta.Created = parsed

	return &meta, nil
}

func getImages(path string) (string, []string, error) {
	dirPath := filepath.Join(path, "images")
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil, nil
		}
		return "", nil, fmt.Errorf("error while reading folder %s: %v", dirPath, err)
	}
	images := []string{}
	for _, file := range files {
		images = append(images, file.Name())
	}
	return dirPath, images, nil
}

func getEtag(buffer *[]byte) string {
	hash := md5.Sum(*buffer)
	return fmt.Sprintf("%x", hash)
}
