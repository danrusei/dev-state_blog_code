package site

import (
	"io/ioutil"
	"mime"
	"path/filepath"
)

//Image holds image file
type Image struct {
	Name    string
	Imgfile []byte
	Mime    string
	Etag    string
}

//ImgConfig holds the configuration for images
type ImgConfig struct {
	ImgPath string
	Cache   *Cache
}

//NewImgConfig generate new post configuration
func NewImgConfig(imgpath string) *ImgConfig {
	return &ImgConfig{
		ImgPath: imgpath,
		Cache:   NewCache(),
	}
}

func (i *ImgConfig) Read() error {
	imgpath := i.ImgPath

	files, err := ioutil.ReadDir(imgpath)
	if err != nil {
		return err
	}

	for _, file := range files {
		image, err := getpostsImage(imgpath, file.Name())
		if err != nil {
			return err
		}
		i.Cache.Set(image.Name, image)
	}

	return nil
}

//GetImage get the image content from cache
func (i *ImgConfig) GetImage(key string) *Image {
	item := i.Cache.Get(key)
	if item == nil {
		return nil
	}

	return item.(*Image)
}

func getpostsImage(imgpath string, imgfile string) (*Image, error) {
	contents, err := ioutil.ReadFile(imgpath + imgfile)
	if err != nil {
		return nil, err
	}

	ext := filepath.Ext(imgfile)
	mimeType := mime.TypeByExtension(ext)

	return &Image{
		Name:    imgfile,
		Imgfile: contents,
		Mime:    mimeType,
		Etag:    getEtag(&contents),
	}, nil

}
