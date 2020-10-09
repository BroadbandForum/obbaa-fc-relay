# compile protobuffs
protoc --proto_path=proto --go_out=plugins=grpc:pb --go_opt=paths=source_relative ./proto/control_relay_packet_filter_service.v1.proto ./proto/control_relay_service.proto

# compile the steering app client
go build -o bin/sdn_client

# docker image
docker build .\
   -t obbaa/steering_app_client \
   -f ./sdn_client.dockerfile

# run steering app client
docker run -it \
   --name STEERING_APP_CLIENT  \
   --rm \
   --env OBBAA_ADDRESS=192.168.56.102 \
   --env OBBAA_PORT=9292 \
   obbaa/steering_app_client
