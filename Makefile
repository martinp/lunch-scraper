PREFIX ?= /usr/local
NAME=lunch-scraper

all: test build

fmt:
	gofmt -tabs=false -tabwidth=4 -w=true *.go

deps:
	go get -u github.com/PuerkitoBio/goquery
	go get -u github.com/docopt/docopt-go

build:
	go build -o $(NAME) main.go

test:
	go test

install:
	cp $(NAME) $(PREFIX)/bin/$(NAME)
