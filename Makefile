build:
	docker build -t d4r:latest .

compile:
	go build -o d4r main.go
