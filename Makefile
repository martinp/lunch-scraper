PREFIX ?= /usr/local
NAME=lunch-scraper

all: deps build

fmt:
	gofmt -w=true *.go

deps:
	go get -d -v

build:
	@mkdir -p bin
	go build -o bin/$(NAME)

install:
	cp $(NAME) $(PREFIX)/bin/$(NAME)
