package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/aarzilli/toml"
	"github.com/speedata/go-epub"
)

type ebpubconf struct {
	Author   string
	Title    string
	Filename string
	Imagedir string
	Cover    string
	CSS      string
	Fonts    []string
	Sections [][]string
}

func isImage(fn string) bool {
	return strings.HasSuffix(fn, ".png") || strings.HasSuffix(fn, ".svg") || strings.HasSuffix(fn, ".jpg")
}

func dothings() error {
	configdata, err := ioutil.ReadFile("config.toml")
	if err != nil {
		return err
	}

	var conf ebpubconf
	if _, err := toml.Decode(string(configdata), &conf); err != nil {
		return err
	}

	ep := epub.NewEpub(conf.Title)
	ep.SetAuthor(conf.Author)
	ep.SetCover(conf.Cover, "")
	for _, fnt := range conf.Fonts {
		fn, err := ep.AddFont(fnt, "")
		if err != nil {
			return err
		}
		fmt.Println("Adding font", fnt, " - resulting filename", fn)
	}
	var cssfilename string
	if cssfilename, err = ep.AddCSS(conf.CSS, ""); err != nil {
		return err
	}
	fmt.Println("Adding CSS", conf.CSS, " - resulting filename", cssfilename)

	for _, sec := range conf.Sections {
		filename := sec[0]
		destfilename := strings.TrimPrefix(filename, "out/")
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return err
		}
		title := sec[1]
		children := sec[2:]
		for {
			if len(children) > 1 {
				children[1] = "xhtml/" + destfilename + "#" + children[1]
				children = children[2:]
			} else {
				break
			}
		}

		fn, err := ep.AddSection(string(b), title, destfilename, cssfilename, sec[2:]...)
		if err != nil {
			return err
		}
		fmt.Println("resulting filename", fn)
	}
	imgs, err := filepath.Glob(filepath.Join(conf.Imagedir, "*"))
	if err != nil {
		return err
	}

	for _, img := range imgs {
		if isImage(img) {
			_, err := ep.AddImage(img, strings.TrimPrefix(img, conf.Imagedir+"/"))
			if err != nil {
				return err
			}
			// fmt.Println("resulting filename", fn)
		}
	}

	return ep.Write(conf.Filename)
}

func main() {
	if err := dothings(); err != nil {
		fmt.Println(err)
	}

}
