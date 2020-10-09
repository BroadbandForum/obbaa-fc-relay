# OB-BAA Control-Relay

This repository contains the control relay feature of the OB-BAA project [OB-BAA implementation](https://github.com/BroadbandForum/obbaa).

The control relay feature in OB-BAA provides the capability of relaying packets from an access node to an application endpoint in the SDN Management and Control functional layer.

The control relay application will run as a micro-service in a docker container, with two ports open to the "public", these ports have to be specified when the container starts up. One of this ports is used for network equipment to connect to the gRPC Standard Plugin (Southbound) and the other is used for controllers to connect as gRPC clients (Northbound).

# Requirements

- Docker
- Go 1.14.4
- Protobuf compiler


## Control-Relay components:

  - Core: Provides a function that is invoked by plugins, a function that is responsible for sending packets to network controllers;
    
  - Northbound API: It is a grpc Server/Client, it allows Control-Relay to connect to various applications, both as a server and as a client.
    
  - Plugins/Adapters: each vendor can provide its own plugin/adapter. Plugins must implement a packetOutCallback that the core will invoke whenever it wants to send a packet to a device.


The Control-Relay includes a standard plugin that communicates via gRPC **(gRPC Standard Plugin)**, which corresponds to the OB-BAA standard adapter. This plugin will communicate with equipment that comply to the control_relay_service.proto. 

The plugins are separate **“.so”** files from the Control-Relay that can be added or removed during the container runtime. The objective is for each vendor to provide its plugin so that Control-Relay supports equipment that uses different communication protocols.



# Build Intructions:

The following commands must be run in the control-relay subdirectory of the source code

### compile protos
`protoc --proto_path=proto --go_out=plugins=grpc:pb --go_opt=paths=source_relative ./proto/control_relay_packet_filter_service.v1.proto ./proto/control_relay_service.proto`.

The 'control\_relay\_packet\_filter\_service.v1.proto' is an experimental feature and allows a SDN Controller or application connecting in the Northbound to apply selective filtering. 

### compile the plugin
`go build -buildmode=plugin -o bin/plugin-standard/BBF-OLT-standard-1.0.so plugins/standard/BBF-OLT-standard-1.0.go`

### compile the control relay
`go build -o bin/control-relay`

### create a bundle with the binaries
`./create-bundle.sh`

### docker image
```
docker build . \
  -t  broadbandforum/obbaa-control-relay \
  -f ./Dockerfile
```

### run control relay
```
docker run -d \
 --name CONTROL_RELAY  \
 --rm  \
 -v $(pwd)/bin/plugins:/control_relay/plugin-repo  \
 -p50055:50055   -p50052:50052  \
 --env OBBAA_ADDRESS=$(docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' baa) \
 --env OBBAA_PORT=9292 \
  broadbandforum/obbaa-control-relay  
```
# Examples

The following examples are included
  - **olt_app_dhcp_go** - an application for simulating the device. It reads the packet from a file and sends it to the control relay service.
  - **steering_app_client** - Simulates a SDN Controller connecting as client to the Control Relay microservice
  - **steering_app_server** - Simulates a SDN Controller acting as server. Its IP address/hostname must be given in the SDN_MC_SERVER_LIST environment variable when starting the Control Relay micro-service.

