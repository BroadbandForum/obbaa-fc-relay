FROM golang:1.14.4 AS builder

RUN apt update && apt-get install protobuf-compiler unzip -y
RUN go get -u github.com/golang/protobuf/protoc-gen-go
WORKDIR /opt/control-relay
COPY . .

RUN protoc --proto_path=proto --go_out=plugins=grpc:pb --go_opt=paths=source_relative \
  --go_opt=Mtr477_cpri_message.proto=control_relay/pb   \
  --go_opt=Mtr477_cpri_service.proto=control_relay/pb   \
  ./proto/control_relay_packet_filter_service.v1.proto  \
  ./proto/control_relay_service.proto                   \
  ./proto/tr477_cpri_message.proto    \
  ./proto/tr477_cpri_service.proto

### compile the plugin
RUN go build -buildmode=plugin -o bin/plugin-standard/BBF-OLT-standard-1.0.so plugins/standard/BBF-OLT-standard-1.0.go

### compile the control relay
RUN go build -o bin/control-relay

### create a bundle with the binaries
RUN chmod +x create-bundle.sh && ./create-bundle.sh


#####################################
# Build final docker image
#####################################
FROM ubuntu:18.04

RUN apt-get update && \
    apt-get upgrade -y

WORKDIR /control_relay

COPY --from=builder /opt/control-relay/dist/control-relay.tgz /control_relay
RUN tar xvfz control-relay.tgz && \
    rm -f control-relay.tgz

CMD ["./control-relay"]    
