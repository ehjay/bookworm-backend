# bookworm-backend

# About

This is the server for Bookworm, it is a RESTful service written in Go that stores a book catalogue in MongoDB

# Getting started

Install go:

```
wget https://storage.googleapis.com/golang/go1.8.1.linux-amd64.tar.gz
tar -xvzf go1.8.1.linux-amd64.tar.gz
sudo mv go /usr/local
```

Follow the standard go workspace and set GOPATH https://golang.org/doc/code.html

Install godep:

```
go get github.com/tools/godep
```

Install dependencies:

```
godep go install
```

Install mongodb https://docs.mongodb.com/master/tutorial/install-mongodb-on-ubuntu/

Setup mongo config:

```
sudo vi /etc/systemd/system/mongodb.service
```

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

```
sudo systemctl enable mongodb
```

Seed the db:

Start the server:

```
go run main.go
```

# Examples

Add a book:

```
curl -X POST http://localhost:8080/book -H "application/json" -d '{"title":"Hamlet","author":"Shakespeare"}'
```

Get the books:

```
curl http://localhost:8080/books
```

Delete a book:

```
curl -X DELETE http://localhost:8080/book/:id
```

# Workflow

```
// write a test
// write some code
godep go test
```

# Add a new dependency

```
go get foo/bar
godep save
```

# Update a dependency

```
go get -u foo/bar
godep update foo/bar
```

# Using mongodb

CLI:

```
mongo
```

Start:

```
sudo systemctl start mongodb
```

Stop:

```
sudo systemctl stop mongodb
```

Status:

```
sudo systemctl status mongodb
```

Default logs:

```
/var/log/mongodb/mongod.log
```
