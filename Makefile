NAME:=$(shell basename `git rev-parse --show-toplevel`)
HASH:=$(shell git rev-parse --verify --short HEAD)

all: setbin

setbin: build
	cp $(NAME) /usr/local/bin

build: dep
	go build -ldflags "-X main.version=$(HASH)" -o $(NAME)

dep:
	go get -v
