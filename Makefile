serve:
	DB_SOURCE="foo.db" go run main.go

test-unit:
	go test -v controller/*.go