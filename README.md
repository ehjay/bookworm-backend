# bookworm-backend

# Getting started

Install go:

  wget https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz
  tar -xvzf go1.8.1.linux-amd64.tar.gz
  sudo mv go /usr/local

Follow the standard go workspace and set GOPATH:

  https://golang.org/doc/code.html

Install godep:

  go get github.com/tools/godep

Install dependencies:

  godep go install

Start the server:

  go run main.go

Example request:

  curl http://localhost:8080/books

# Workflow

  // write a test
  // write some code
  godep go test

# Add a new dependency

  go get foo/bar
  godep save

# Update a dependency

  go get -u foo/bar
  godep update foo/bar
