.PHONY: fmt test build package tag run push clean

BIN_NAME=node-balance-retriever
IMAGE_NAME=node-balance-retriever

default: fmt test build package tag run

fmt:
	goimports -w .
	go fmt ./...

test: fmt
	go test ./...

build: test
	GOOS=linux GOARCH=arm go build -o artifacts/${BIN_NAME} .
	chmod +x artifacts/${BIN_NAME}

package: build
	docker build --build-arg BIN_NAME=${BIN_NAME} -t $(IMAGE_NAME):local .

tag: package
	docker tag $(IMAGE_NAME):local $(IMAGE_NAME):latest

run: tag
	docker run --name ${BIN_NAME} --rm $(IMAGE_NAME):latest

push: tag
	docker push $(IMAGE_NAME):latest

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}
