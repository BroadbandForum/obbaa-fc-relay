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
* Control Relay OLT simulator file
*
* Created by Filipe Claudio(Altice Labs) on 01/09/2020
 */

package main

import (
	"context"
	"flag"
	"io"
	"io/ioutil"
	"log"
	pb "olt_app_dhcp/pb"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var address = flag.String("address", "localhost:50052", "Address for standar grpc plugin")
var packetPath = flag.String("packet_path", "./ont-dhcp.bin", "Path of the captured packet")
var deviceName = flag.String("olt_serial", "AlticeLabs_OLT22", "Identification of the OLT")
var slot = flag.String("olt_port", "", "Olt port number")
var pon = flag.String("pon", "", "Pon")
var onuID = flag.String("onu_id", "", "Onu ID")
var deviceInterface = ""
var count = 0

func main() {
	flag.Parse()
	deviceInterface = *slot + "-" + *pon + "-" + *onuID
	log.Println(*packetPath)
	buf, errReadFile := ioutil.ReadFile(*packetPath)
	if errReadFile != nil {
		log.Fatalf("Packet not found %v", errReadFile)
	}

	// Set up a connection to the server.
	conn, errConnServer := grpc.Dial(*address, grpc.WithInsecure())
	if errConnServer != nil {
		log.Fatalf("Did not connect: %v", errConnServer)
	}
	defer conn.Close()

	addHelloService := pb.NewControlRelayHelloServiceClient(conn)
	addClientService := pb.NewCpriMessageClient(conn)

	_, err := addHelloService.Hello(context.Background(), &pb.HelloRequest{
		LocalEndpointHello: &pb.HelloRequest_Device{
			Device: &pb.DeviceHello{
				DeviceName: *deviceName,
			},
		},
	})
	if err != nil {
		log.Println("Hello Service failed")
		log.Println("Could not create a client connection to the given target: ")
		log.Println("Target is the grpc standard plugin: ", *address)
		log.Println("Error: ", err)
	} else {
		log.Println("Control App: OK")
	}

	go waitForPacketsOnStream(addClientService)
	for {
		sendPacket(addClientService, buf)
		time.Sleep(5 * time.Second)
	}
}

func sendPacket(client pb.CpriMessageClient, buf []byte) {
	metadata := pb.CpriMetaData{
		Generic: &pb.GenericMetadata{
			DeviceName:      *deviceName,
			DeviceInterface: deviceInterface,
		},
	}

	_, err := client.CpriTx(context.Background(), &pb.CpriMsg{
		MetaData: &metadata,
		Packet:   buf})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func waitForPacketsOnStream(client pb.CpriMessageClient) {
	stream, err := client.ListenForCpriRx(context.Background(), &empty.Empty{})

	if err != nil {
		log.Println("ListenForPacketRx Service failed")
		log.Println("Could not create a client connection to the given target: ")
		log.Println("Target is the Control Relay: ", *address)
		log.Println("Error: ", err)
	}

	go on()
	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if packet != nil {
			count++
			log.Println("Packet Received. Request number: ", count)
			/* log.Println("Successfully received package: ", packet)
			log.Println("######### Packet Out #########")
			log.Println("DeviceName: ", packet.DeviceName)
			log.Println("DeviceInterface: ", packet.DeviceInterface)
			log.Println("Originating rule: ", packet.OriginatingRule)
			log.Println("Packet: ", packet.Packet)
			log.Println("############# END ############") */
		}
	}
}

func on() {
	for {
		log.Println("Waiting for packets on stream")
		time.Sleep(10 * time.Second)
	}
}
