package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/docopt/docopt-go"
	"io/ioutil"
	_log "log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

const URL = "http://www.dagbladet.no/tegneserie/lunch/"

var log *_log.Logger

func GetImageUrl() (string, error) {
	doc, err := goquery.NewDocument(URL)
	if err != nil {
		return "", err
	}
	selection := doc.Find("img#lunch-stripe")
	if len(selection.Nodes) == 0 {
		return "", errors.New("Node was not found")
	}
	src, exists := selection.First().Attr("src")
	if !exists {
		return "", errors.New("Node is missing the 'src' attribute")
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

func WriteFile(path string, filename string, data []byte, force bool,
	quiet bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("Path does not exist: %s", path)
	}
	dst := filepath.Join(path, filepath.Base(filename))
	if IsExists(dst) && !force {
		if !quiet {
			log.Printf("File already exists: %s", dst)
		}
		return
	}
	if err := ioutil.WriteFile(dst, data, 0644); err != nil {
		log.Fatalf("Failed to write file: %s", dst)
	}
	if !quiet {
		log.Printf("Wrote file: %s", dst)
	}
}

func parseDirOption(arguments map[string]interface{}) string {
	arg := arguments["--dir"]
	if arg == nil {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatalln(err)
		}
		dir, err := filepath.Rel(wd, wd)
		if err != nil {
			log.Fatalln(err)
		}
		return dir
	}
	return arg.(string)
}

func main() {
	log = _log.New(os.Stderr, "", 0)

	usage := `Tool for scraping Lunch comic strips

    Usage:
      lunch-scraper [-q] [-f] [-d DIR]
      lunch-scraper -h | --help

    Options:
      -h --help             Show help.
      -q --quiet            Be quiet.
      -f --force            Overwrite if file already exists.
      -d --dir=DIR          Downloaded file to DIR instead of working directory.`

	arguments, _ := docopt.Parse(usage, nil, true, "", false)
	quiet := arguments["--quiet"].(bool)
	force := arguments["--force"].(bool)
	path := parseDirOption(arguments)

	url, err := GetImageUrl()
	if err != nil {
		log.Fatalln(err)
	}
	data, err := GetImage(url)
	if err != nil {
		log.Fatalln(err)
	}
	WriteFile(path, GetFilename(), data, force, quiet)
}
