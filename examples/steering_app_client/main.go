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
* Control Relay SDN controller simulator file
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
	pb "sim_sdn_client/pb"
	"time"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

var address = flag.String("control_relay_address", "localhost:50055", "Control Relay Address")
var packetPath = flag.String("packet_path", "data/ont-dhcp.bin", "Path of the captured packet")
var deviceName = flag.String("olt_serial", "olt1", "Identification of the OLT")
var controlAppName = "controlAppName"
var count = 0
var packets []*ControlRelayPacketInternal
var filter = pb.ControlRelayPacketFilterList{}

type ControlRelayPacketInternal struct {
	Device_name      string
	Device_interface string
	Originating_rule string
	Packet           []byte
}

func main() {

	flag.Parse()
	buf, errReadFile := ioutil.ReadFile(*packetPath)
	if errReadFile != nil {
		log.Fatalf("Packet not found %v", errReadFile)
	}
	boot(buf)

	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Println("Failed to establish connection with the Control Relay")
		log.Println("Control Relay address: ", *address)
		log.Fatalf("Error: ", err)
		conn.Close()
	}

	addHelloService := pb.NewControlRelayHelloServiceClient(conn)
	addClientService := pb.NewCpriMessageClient(conn)
	//addFilterService := pb.NewControlRelayPacketFilterServiceClient(conn)

	response, err := addHelloService.Hello(context.Background(), &pb.HelloRequest{
		LocalEndpointHello: &pb.HelloRequest_Controller{
			// send nothing
		},
	})
	if err != nil {
		log.Println("Hello Service failed")
		log.Println("Could not create a client connection to the given target: ")
		log.Println("Target is the Control Relay: ", *address)
		log.Println("Error: ", err)
	} else {
		log.Println("Control App: OK")
		controlAppName = response.GetDevice().DeviceName
	}

	//setControlRelayPacketFilter(addFilterService)
	go waitForPacketsOnStream(addClientService)

	for {
		sendPacket(addClientService, *packets[0])
		//sendPacket(addClientService, *packets[1])
		//sendPacket(addClientService, *packets[2])
		time.Sleep(10 * time.Second)
	}
}

func sendPacket(client pb.CpriMessageClient, packet ControlRelayPacketInternal) {
	metadata := pb.CpriMetaData{
		Generic: &pb.GenericMetadata{
			DeviceName:      packet.Device_name,
			DeviceInterface: packet.Device_interface,
			OriginatingRule: packet.Originating_rule,
		},
	}

	_, err := client.CpriTx(context.Background(), &pb.CpriMsg{
		MetaData: &metadata,
		Packet:   packet.Packet,
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}

func addPacket(deviceName string, deviceInterface string, rule string, buf []byte) {
	packets = append(packets, &ControlRelayPacketInternal{
		Device_name:      deviceName,
		Device_interface: deviceInterface,
		Originating_rule: rule,
		Packet:           buf,
	})
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
			log.Println("Request number: ", count)
			/* log.Println("######### Packet Out #########")
			log.Println("DeviceName: ", packet.DeviceName)
			log.Println("DeviceInterface: ", packet.DeviceInterface)
			log.Println("Packet: ", packet.Packet)
			log.Println("############# END ############") */
		}
	}
}

func setControlRelayPacketFilter(client pb.ControlRelayPacketFilterServiceClient) {
	_, err := client.SetControlRelayPacketFilter(context.Background(), &filter)
	if err != nil {
		log.Println("Set Filters: ", err)
		return
	}
	log.Println("Set Filters: OK")
}

func addFilter(include bool, devicename string, deviceinterface string, rule string) {
	filter.Filter = append(filter.Filter, &pb.ControlRelayPacketFilterList_ControlRelayPacketFilter{
		Type: func() pb.ControlRelayPacketFilterList_ControlRelayPacketFilter_FilterType {
			if include {
				return pb.ControlRelayPacketFilterList_ControlRelayPacketFilter_INCLUDE
			}
			return pb.ControlRelayPacketFilterList_ControlRelayPacketFilter_EXCLUDE
		}(),
		DeviceName:      devicename,
		DeviceInterface: deviceinterface,
		OriginatingRule: rule,
	})
}

func clearControlRelayPacketFilter(client pb.ControlRelayPacketFilterServiceClient) {
	ok, err := client.ClearControlRelayPacketFilter(context.Background(), &empty.Empty{})
	if err != nil {
		log.Println("Clear Filters: ", err)
		return
	}
	log.Println("Clear Filters: ", ok)
	filter.Filter = nil
}

func on() {
	for {
		log.Println("Waiting for packets on stream")
		time.Sleep(10 * time.Second)
	}
}

const INCLUDE = true
const EXCLUDE = false

func boot(buf []byte) {
	addFilter(INCLUDE, "AlticeLabs_OLT2453452", "1", "")
	addFilter(INCLUDE, "AlticeLabs_OLT22", "", "")
	addFilter(INCLUDE, "AlticeLabs_OLT2223", "2", "bbb")
	addFilter(INCLUDE, "AlticeLabs_OLT221", "", "")
	addFilter(EXCLUDE, "AlticeLabs_OLT226", "2", "aaa")
	addFilter(INCLUDE, "AlticeLabs_OLT22273", "2", "bbb")
	addFilter(INCLUDE, "AlticeLabs_OLT223", "", "")
	addFilter(INCLUDE, "AlticeLabs_OLT242", "", "")
	addFilter(INCLUDE, "AlticeLabs_OLT2772", "2", "aaa")
	addFilter(INCLUDE, "AlticeLabs_OLT2225473", "2", "bbb")
	addFilter(INCLUDE, "AlticeLabs_OLT2674572", "", "")
	addFilter(INCLUDE, "AlticeLabs_OLT222", "", "")
	addFilter(INCLUDE, "AlticeLabs_OLT27467462", "2", "aaa")
	addFilter(INCLUDE, "AlticeLabs_OLT2224743", "2", "bbb")
	addFilter(INCLUDE, "AlticeLabs_OLT22467457", "", "")
	addFilter(INCLUDE, "AlticeLabs_OLT22", "1", "aaa")
	addFilter(INCLUDE, "AlticeLabs_OLT2223453453", "2", "bbb")
	addFilter(INCLUDE, "AlticeLabs_OLT235345342", "", "")
	addFilter(INCLUDE, "AlticeLabs_OLT22", "2", "aaaa")
	addFilter(INCLUDE, "AlticeLabs_OLT22253453", "2", "bbb")
	addFilter(INCLUDE, "AlticeLabs_OLT25345342", "", "")
	addFilter(EXCLUDE, "", "", "")

	addPacket("OLT1", "1", "", buf)
	addPacket("AlticeLabs_OLT222", "1", "", buf)
	addPacket("AlticeLabs_OLT223", "1", "", buf)
	addPacket("AlticeLabs_OLT2223442", "1", "", buf)
}
