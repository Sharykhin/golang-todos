Simple Todo App:
============

Requirements:
------------

Uses [Golang](https://golang.org/) 1.9

Usage:
------

1. Make docker build:
```bash
docker build -t golang-todoapp .
```

2. Run container:
```bash
docker run --publish 8082:8082 --name todoapp_service --rm golang-todoapp
```

3. Go to http://localhost:8082