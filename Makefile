
gladder: *.go bindata.go
	go install
	go build # so we have an executable for docker

bindata.go: resources/templates/* resources/css/* resources/js/*
	go-bindata resources/...

clean:
	rm gladder

docker-image: gladder Dockerfile
	docker build -t asokoloski/gladder .

