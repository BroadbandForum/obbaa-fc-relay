/*
 * Copyright 2020 Broadband Forum
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
* Control Relay gRPC Standard Plugin file
*
* Created by Filipe Claudio(Altice Labs) on 01/09/2020
*/

package main

import (
	"context"
	"net"
	"os"
	"sync"

	"control_relay/pb"
	core "control_relay/syscore"

	"control_relay/utils/log"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	tcpListener, tcpListenerErr = net.Listen("tcp", "0.0.0.0:"+PLUGIN_PORT)
	grpcServer                  = grpc.NewServer()
	PLUGIN_PORT                 = os.Getenv("PLUGIN_PORT")
)

var devices = make(map[string]*connectedDevice)
var mutex = sync.RWMutex{}

type plugin string

type connectedDevice struct {
	deviceName string
	stream     pb.ControlRelayPacketService_ListenForPacketRxServer
	ipDevice   string
	network    string
	ch         chan string
}

type controlRelayHelloService struct {
	pb.UnimplementedControlRelayHelloServiceServer
}

type controlRelayPacketService struct {
	pb.UnimplementedControlRelayPacketServiceServer
}

// Hello is one of the services provided by the proto, and is used mainly to
// establish and keep the connection between the Control Relay and devices.
// Is invoked by the device, and if the connection is successfully established,
// the control relay stores the reference of that device for later sending packets
func (s *controlRelayHelloService) Hello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {

	if in == nil || ctx == nil {
		log.Warning("Standard gRPC Plugin: Device without information.")
		return nil, nil
	}

	if in.GetDevice().DeviceName == "" {
		log.Warning("Standard gRPC Plugin: Device without name.")
		return nil, nil
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		log.Warning("Standard gRPC Plugin: Without information about the device")
	} else {
		if addDevice(in.GetDevice().DeviceName, p.Addr.String(), p.Addr.Network()) {
			log.Info("Standard gRPC Plugin: A new device is connected and is waiting to receive packets: ")
			log.Info("Standard gRPC Plugin: -----------> Device Name: ", in.GetDevice().DeviceName)
			log.Info("Standard gRPC Plugin: -----------> IP Device: ", p.Addr.String())
		}
	}

	return &pb.HelloResponse{
		RemoteEndpointHello: &pb.HelloResponse_Controller{
			Controller: &pb.ControllerHello{
				// nothing
			},
		},
	}, nil
}

// PacketTx is one of the services provided by the proto, and is used to send
// packages to the SDN. PacketTx is invoked by the device.
// When invoked, it receives a new packet, and will forward the packet to the
// the core through the PacketInCallBack function.
func (s *controlRelayPacketService) PacketTx(ctx context.Context, in *pb.ControlRelayPacket) (*empty.Empty, error) {

	if in == nil || ctx == nil {
		log.Warning("Standard gRPC Plugin: Device without information.")
		return nil, nil
	}

	if in.DeviceName == "" {
		log.Warning("Standard gRPC Plugin: Packet without name.")
		return nil, nil
	}

	/* log.Info(
		"Standard gRPC Plugin: ******** PACKET IN ********",
		"\nStandard gRPC Plugin: Device Name: ", in.DeviceName,
		"\nStandard gRPC Plugin: Device interface: ", in.DeviceInterface,
		"\nStandard gRPC Plugin: Packet: ",
		"\n*********** END ***********",
	) */

	go func() {
		core.PacketInCallBack(
			core.ControlRelayPacketInternal{
				Device_name:      in.DeviceName,
				Device_interface: in.DeviceInterface,
				Originating_rule: in.OriginatingRule,
				Packet:           in.Packet,
			},
		)
	}() 

	log.Debug("Standard gRPC Plugin: PacketIn sent and received successfully by the Control Relay")

	return &empty.Empty{}, nil
}

func (s *controlRelayPacketService) ListenForPacketRx(e *empty.Empty, stream pb.ControlRelayPacketService_ListenForPacketRxServer) error {

	if stream == nil {
		return nil
	}
	
	ch := make(chan string)
	p, ok := peer.FromContext(stream.Context())

	if !ok {
		log.Warning("Standard gRPC Plugin: Without information about the device")
		return nil
	}
	if setDeviceStream(p.Addr.String(), stream, ch) {
		log.Info("Standard gRPC Plugin: Stream open and ready to send packets")
	} else {
		deleteDevice(p.Addr.String())
		return nil
	}

	for {
		<-ch
		break
	}

	return nil
}

func (p plugin) PacketOutCallBack(packet core.ControlRelayPacketInternal) {

	if _, ok := devices[packet.Device_name]; ok {
		if devices[packet.Device_name].stream == nil {
			return
		}
		if err := devices[packet.Device_name].stream.Send(&pb.ControlRelayPacket{
			DeviceName:      packet.Device_name,
			DeviceInterface: packet.Device_interface,
			OriginatingRule: packet.Originating_rule,
			Packet:          packet.Packet,
		}); err != nil {
			log.Warning("Standard gRPC Plugin: Failed to send the package to the device")
			log.Warning("Standard gRPC Plugin: Device Name: ", packet.Device_name)
			log.Warning("Standard gRPC Plugin: Error: ", err)
			if deleteDevice(packet.Device_name) {
				log.Info("Standard gRPC Plugin: Device cleared from internal cache")
			}
		}
	} else {
		log.Warning("Standard gRPC Plugin: Device not found in internal cache")
		log.Warning("Standard gRPC Plugin: DeviceName: ", packet.Device_name)
	}

}

func addDevice(name string, ip string, net string) bool {
	mutex.Lock()
	if devices[name] != nil {
		if devices[name].stream.Context().Err() != nil {
			goto SkipToEnd
		} 
		log.Warning("Standard gRPC Plugin: Device ", name, " already connected")
		mutex.Unlock()
		return false
		SkipToEnd:	
	}

	devices[name] = &connectedDevice{
		deviceName: name,
		stream:     nil,
		ipDevice:   ip,
		network:    net,
		ch:         nil,
	}
	mutex.Unlock()
	return true
}

func setDeviceStream(ip string, stream pb.ControlRelayPacketService_ListenForPacketRxServer, ch chan string) bool {
	mutex.Lock()
	for d := range devices {
		if ip == devices[d].ipDevice {
			devices[d].stream = stream
			devices[d].ch = ch
			mutex.Unlock()
			return true
		}
	}
	mutex.Unlock()
	return false
}

func deleteDevice(name string) bool {
	mutex.Lock()
	if _, ok := devices[name]; ok {
		devices[name].ch <- "close"
		delete(devices, name)
		mutex.Unlock()
		return true
	}
	mutex.Unlock()
	return false
}

func (p plugin) Start() {
	log.Info("Standard gRPC Plugin: Initializing Southbound gRPC server")
	if tcpListenerErr != nil {
		log.Error("Standard gRPC Plugin: Could not initialize Southbound gRPC server")
		log.Error("Standard gRPC Plugin: Error: ", tcpListenerErr)
	}

	addHelloServ := controlRelayHelloService{}
	addPacketServ := controlRelayPacketService{}

	// The & symbol points to the address of the stored value.
	pb.RegisterControlRelayHelloServiceServer(grpcServer, &addHelloServ)
	pb.RegisterControlRelayPacketServiceServer(grpcServer, &addPacketServ)

	if err := grpcServer.Serve(tcpListener); err != nil {
		log.Fatal("Standard gRPC Plugin: Error: ", err)
	}
}

func (p plugin) Stop() {
	log.Info("Standard gRPC Plugin is closing...")
	grpcServer.Stop()
	tcpListener.Close()
	log.Warning("gRPC plugin is closed.")
}

// Plugin is a symbol that is being exported
var Plugin plugin

// go build -buildmode=plugin -o bin/plugins/BBF-OLT-standard-1.0.so plugins/standard/BBF-OLT-standard-1.0.go
