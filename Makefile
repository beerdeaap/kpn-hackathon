install:
	sudo apt-get install -y --force-yes golang-go
go:
	go build queens.go
	./queens

