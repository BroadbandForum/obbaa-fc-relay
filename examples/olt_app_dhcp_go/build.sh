# compile olt
rm  ./bin/olt.*
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/olt.amd64 main.go 
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o ./bin/olt.arm64 main.go