
## Golang Stuff
GOCMD=go
GORUN=$(GOCMD) run

SERVICE=sut-product-go

proto-gen:
	protoc --proto_path=proto/ --go_out=paths=source_relative,plugins=grpc:./pb proto/*/*.proto

init:
	$(GOCMD) mod init $(SERVICE)

tidy:
	$(GOCMD) mod tidy

run:
	echo "for local development, please run: make run ENV=local"
	$(GORUN) cmd/main.go