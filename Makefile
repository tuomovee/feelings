.PHONY: deps clean build

deps:
	dep ensure

clean: 
	rm ./cmd/get/get
	rm ./cmd/put/put
	
build:
	GOOS=linux GOARCH=amd64 go build -o ./cmd/get/get ./cmd/get
	GOOS=linux GOARCH=amd64 go build -o ./cmd/put/put ./cmd/put

package ./cmd/get/get ./cmd/put/put:
	aws cloudformation package --template-file template.yaml --s3-bucket feelings-lambdas --output-template-file packaged.yaml

deploy:
	aws cloudformation deploy --template-file packaged.yaml --stack-name feelings --capabilities CAPABILITY_IAM
