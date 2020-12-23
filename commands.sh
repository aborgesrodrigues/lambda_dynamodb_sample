# compile
GOOS=linux go build -o main

# update code to sam
sam build

# call lambda function
sam local invoke "GoLambda" -e input.json

# start api for the local lambda function
sam local start-api

#
GOARCH=amd64 GOOS=linux go build -o ./dlv/dlv github.com/go-delve/delve/cmd/dlv

# compile to debug
GOARCH=amd64 GOOS=linux go build -gcflags='-N -l' -o main .

# start api in debug mode
sam local invoke -d 2345 --debugger-path ./dlv --debug-args "-delveAPI=2" --debug