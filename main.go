package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/jessevdk/go-flags"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const URL = "http://www.dagbladet.no/tegneserie/lunch/"

func GetImageUrl() (string, error) {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return "", err
	}
	selection := doc.Find("img.image_comic")
	if len(selection.Nodes) == 0 {
		return "", fmt.Errorf("node not found")
	}
	src, exists := selection.First().Attr("src")
	if !exists {
		return "", fmt.Errorf("src attribute not found")
	}
	return src, nil
}

func GetImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func GetFilename() string {
	return time.Now().Format("2006-01-02") + ".gif"
}

func IsExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func WriteFile(path string, filename string, data []byte, force bool) (string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", err
	}
	dst := filepath.Join(path, filepath.Base(filename))
	if IsExists(dst) && !force {
		return "", fmt.Errorf("file exists: %s", dst)
	}
	if err := ioutil.WriteFile(dst, data, 0644); err != nil {
		return "", err
	}
	return dst, nil
}

func main() {
	var opts struct {
		Quiet bool   `short:"q" long:"quiet" description:"Be quiet"`
		Force bool   `short:"f" long:"force" description:"Overwrite existing file"`
		Dir   string `short:"d" long:"dir" description:"Directory where file should be saved" value-name:"DIR" required:"true"`
	}
	_, err := flags.ParseArgs(&opts, os.Args)
	if err != nil {
		os.Exit(1)
	}
	url, err := GetImageUrl()
	if err != nil {
		log.Fatal(err)
	}
	data, err := GetImage(url)
	if err != nil {
		log.Fatal(err)
	}
	dst, err := WriteFile(opts.Dir, GetFilename(), data, opts.Force)
	if err != nil {
		log.Fatal(err)
	}
	if !opts.Quiet {
		log.Printf("Wrote file: %s", dst)
	}
}
