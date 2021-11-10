BINARY_NAME={{.Appname}}
USERNAME={{.Username}}

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin ./cmd/main.go
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux ./cmd/main.go

run:
	./${BINARY_NAME}-darwin

build_and_run: 
	build
	run

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux

build_docker:
	docker build -t ${USERNAME}/${BINARY_NAME} .

build_and_push_docker: build_docker
	docker push ${USERNAME}/${BINARY_NAME}
