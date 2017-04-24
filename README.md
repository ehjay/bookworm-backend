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

Install mongodb:

  https://docs.mongodb.com/master/tutorial/install-mongodb-on-ubuntu/

Setup mongo config:

  sudo vi /etc/systemd/system/mongodb.service

with file:

```
[Unit]
Description=High-performance, schema-free document-oriented database
After=network.target

[Service]
User=mongodb
ExecStart=/usr/bin/mongod --quiet --config /etc/mongod.conf

[Install]
WantedBy=multi-user.target
```

Enable db on boot:

  sudo systemctl enable mongodb

Seed the db:

Start the server:

  go run main.go

# Examples

Add a book:

Get the books:

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

# Using mongodb

CLI:

  mongo

Start:

  sudo systemctl start mongodb

Stop:

  sudo systemctl stop mongodb

Status:

  sudo systemctl status mongodb

Default logs:

  /var/log/mongodb/mongod.log
