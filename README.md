Simple Todo App:
============

Uses [Golang](https://golang.org/) 1.9

Requirements:
------------

[Docker](https://www.docker.com/)

Usage:
------

1. Make docker build:
```bash
docker build -t golang-todoapp .
```

2. Run container:
```bash
docker run -p 8082:8082 golang-todoapp --rm golang-todoapp
```

3. Go to http://localhost:8082

API:
----

#### Get list of todos

```bash
GET /?limit={num}&offset={num}
Status: 200 OK
```
Filter params:
- *limit*  - number | default:10. Number of todos that should be return in response
- *offset* - number | default:0. The offset of todos in response

JSON-out:
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "title": "Make todo",
            "description": "Create simple tood application",
            "completed": false,
            "created": "Mon, 05 Feb 2018 08:03:32 UTC"
        }
    ],
    "error": null,
    "meta": {
        "count": 1,
        "total": 17
    }
}
```

#### Create a new todo
```bash
POST /create
Status: 201 Created
```
JSON-in:
```json
{
	"title": "fix tood",
	"description": "fix todo lost"
}
```

JSON-out:
```json
{
    "success": true,
    "data": {
        "id": 18,
        "title": "fix tood",
        "description": "fix todo lost",
        "completed": false,
        "created": "Wed, 07 Feb 2018 07:35:43 UTC"
    },
    "error": null,
    "meta": {}
}
```
