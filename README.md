# sut-product-go
Microservice for managing user requested product

## How to run in local?

1. Create file env and name it `dev.env`. its content can be seen in code block below. 
```
PORT=:50054
NOTIF_HOST=:50055
DB_URL=
```

2. Execute command below
```
make init # initialize go.mod
make tidy # Tidy up go module
```

3. Adding go bin into path env variables
```
export PATH=$PATH:$(go env GOPATH)/bin
```

4. Adding folder with `pb` as name into ther project root directory

5. Generate protobuf by executing command below
```
make proto-gen
```

6. Run the application
```
make run
```
## How to run in local?
```
docker build --tag=sut/product-service --build-arg SERVICE=sut-product-go --build-arg PORT=50054 .
docker run -p 50054:50054 <IMAGE_ID>
```
