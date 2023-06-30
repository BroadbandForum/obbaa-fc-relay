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
* Control Relay gRPC Standard Plugin unit tests file
*
* Created by Filipe Claudio(Altice Labs) on 01/09/2020
 */

package main

import (
	"context"
	"control_relay/pb/tr477"
	"net"
	"reflect"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func Test_controlRelayHelloService_Hello(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *tr477.HelloCpriRequest
	}

	tests := []struct {
		name    string
		s       *controlRelayHelloService
		args    args
		want    *tr477.HelloCpriResponse
		wantErr bool
	}{
		{"test1", &controlRelayHelloService{}, args{context.TODO(), MockHelloRequest("AlticeLabs_OLT22")}, MockHelloResponse(), false},
		{"test2", &controlRelayHelloService{}, args{context.TODO(), MockHelloRequest("")}, MockHelloResponse(), true},
		{"test3", &controlRelayHelloService{}, args{nil, MockHelloRequest("AlticeLabs_OLT22")}, MockHelloResponse(), true},
		{"test4", &controlRelayHelloService{}, args{context.TODO(), nil}, MockHelloResponse(), true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _ := test.s.HelloCpri(test.args.ctx, test.args.in)
			if (!reflect.DeepEqual(got, test.want)) != test.wantErr {
				t.Errorf("controlRelayHelloService.Hello() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_controlRelayPacketService_PacketTx(t *testing.T) {
	type args struct {
		ctx context.Context
		in  *tr477.CpriMsg
	}

	tests := []struct {
		name    string
		s       *controlRelayPacketService
		args    args
		want    *empty.Empty
		wantErr bool
	}{
		{"test1", &controlRelayPacketService{}, args{context.Background(), MockControlRelayPacketRequest("AlticeLabs_OLT22")}, &empty.Empty{}, false},
		{"test2", &controlRelayPacketService{}, args{nil, MockControlRelayPacketRequest("AlticeLabs_OLT22")}, &empty.Empty{}, true},
		{"test3", &controlRelayPacketService{}, args{context.Background(), nil}, &empty.Empty{}, true},
		{"test4", &controlRelayPacketService{}, args{context.Background(), MockControlRelayPacketRequest("")}, &empty.Empty{}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, _ := test.s.CpriTx(test.args.ctx, test.args.in)
			if !reflect.DeepEqual(got, test.want) != test.wantErr {
				t.Errorf("controlRelayPacketService.PacketTx() = %v, want %v", got, test.want)
			}
		})
	}
}

func Test_controlRelayPacketService_ListenForPacketRx(t *testing.T) {
	server := controlRelayPacketService{}
	mock := &mockControlRelayPacketService_ListenForPacketRxServer{}
	server.ListenForCpriRx(&empty.Empty{}, mock)

	client := mockClientConnectionTest(t, ":12345")
	_, err := client.TransferCpri(context.Background(), &empty.Empty{})
	require.NoError(t, err)
}

func Test_plugin_PacketOutCallBack(t *testing.T) {
	type args struct {
		packet *tr477.CpriMsg
	}
	tests := []struct {
		name string
		p    plugin
		args args
	}{
		{"test1", Plugin, args{MockControlRelayPacketInternalRequest("AlticeLabs_OLT22")}},
		{"test2", Plugin, args{MockControlRelayPacketInternalRequest("AlticeLabs_OLT98346239")}},
		{"test3", Plugin, args{MockControlRelayPacketInternalRequest("AlticeLabs_StreamNull")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.PacketOutCallBack(tt.args.packet)
		})
	}
}

func Test_addDevice(t *testing.T) {
	type args struct {
		name string
		ip   string
		net  string
	}

	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{"AlticeLabs_OLT1", "127.0.0.5", "tcp"}, true},
		{"test2", args{"AlticeLabs_OLT2", "127.0.0.6", "tcp"}, true},
		{"test3", args{"AlticeLabs_OLT3", "127.0.0.7", "tcp"}, true},
		{"test4", args{"AlticeLabs_OLT3", "127.0.0.3", "tcp"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addDevice(tt.args.name, tt.args.ip, tt.args.net); got != tt.want {
				t.Errorf("addDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setDeviceStream(t *testing.T) {
	type args struct {
		ip     string
		stream *mockControlRelayPacketService_ListenForPacketRxServer
		ch     chan string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{"127.0.0.1", &mockControlRelayPacketService_ListenForPacketRxServer{}, make(chan string)}, true},
		{"test2", args{"127.0.0.2", &mockControlRelayPacketService_ListenForPacketRxServer{}, make(chan string)}, true},
		{"test3", args{"127.0.0.4", &mockControlRelayPacketService_ListenForPacketRxServer{}, make(chan string)}, false},
		{"test4", args{"127.0.0.3", &mockControlRelayPacketService_ListenForPacketRxServer{}, make(chan string)}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setDeviceStream(tt.args.ip, tt.args.stream, tt.args.ch); got != tt.want {
				t.Errorf("setDeviceStream() = %v, want %v", got, tt.want)
			}
			for d := range devices {
				if tt.args.ip == devices[d].ipDevice {
					go func(ch chan string) {
						<-ch
					}(devices[d].ch)
				}
			}
		})
	}
}

func Test_deleteDevice(t *testing.T) {
	type args struct {
		name string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test1", args{"AlticeLabs_OLT22"}, true},
		{"test2", args{"AlticeLabs_OLT222"}, true},
		{"test3", args{"AlticeLabs_OLT223"}, true},
		{"test4", args{"AlticeLabs_OLT22123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deleteDevice(tt.args.name); got != tt.want {
				t.Errorf("deleteDevice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockServerTest(t *testing.T) {
	lis, err := net.Listen("tcp", "0.0.0.0:12345")
	require.NoError(t, err)
	server := grpc.NewServer()

	addHelloServ := controlRelayHelloService{}
	addPacketServ := controlRelayPacketService{}

	// The & symbol points to the address of the stored value.
	tr477.RegisterCpriHelloServer(server, &addHelloServ)
	tr477.RegisterCpriMessageServer(server, &addPacketServ)

	go server.Serve(lis)
}

func mockClientConnectionTest(t *testing.T, address string) tr477.CpriMessageClient {
	mockServerTest(t)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	require.NoError(t, err)
	return tr477.NewCpriMessageClient(conn)
}

// Mock of grpc server stream
type mockControlRelayPacketService_ListenForPacketRxServer struct {
	grpc.ServerStream
	Packets []*tr477.CpriMsg
}

func (_m *mockControlRelayPacketService_ListenForPacketRxServer) Recv() (*tr477.CpriMsg, error) {
	return &tr477.CpriMsg{
		Header:   &tr477.CpriMsgHeader{},
		MetaData: &tr477.CpriMetaData{},
		Packet:   []byte{},
	}, nil
}

// Mock function Send of grpc server stream
func (_m *mockControlRelayPacketService_ListenForPacketRxServer) Send(packet *tr477.CpriMsg) error {
	_m.Packets = append(_m.Packets, packet)
	return nil
}

// Mock Context of grpc server stream
func (_m *mockControlRelayPacketService_ListenForPacketRxServer) Context() context.Context {
	return context.TODO()
}

// MockControlRelayPacketInternalRequest ...
func MockControlRelayPacketInternalRequest(name string) *tr477.CpriMsg {
	return &tr477.CpriMsg{
		MetaData: &tr477.CpriMetaData{
			Generic: &tr477.GenericMetadata{
				DeviceName:      name,
				DeviceInterface: "1234",
				Direction:       tr477.GenericMetadata_NNI_TO_UNI,
			},
		},
		Packet: []byte{},
	}
}

func MockHelloRequest(name string) *tr477.HelloCpriRequest {
	return &tr477.HelloCpriRequest{
		LocalEndpointHello: &tr477.Hello{},
	}
}

func MockHelloResponse() *tr477.HelloCpriResponse {
	return &tr477.HelloCpriResponse{
		RemoteEndpointHello: &tr477.Hello{},
	}
}

// MockControlRelayPacketRequestSuccess ...
func MockControlRelayPacketRequest(name string) *tr477.CpriMsg {
	metadata := tr477.CpriMetaData{
		Generic: &tr477.GenericMetadata{
			DeviceName:      name,
			DeviceInterface: "1234",
			Direction:       tr477.GenericMetadata_NNI_TO_UNI,
		},
	}
	return &tr477.CpriMsg{
		MetaData: &metadata,
		Packet:   []byte{},
	}
}

func init() {
	addDevice("AlticeLabs_OLT22", "127.0.0.1", "tcp")
	addDevice("AlticeLabs_OLT222", "127.0.0.2", "tcp")
	addDevice("AlticeLabs_OLT223", "127.0.0.3", "tcp")
	addDevice("AlticeLabs_StreamNull", "127.0.0.50", "tcp")

	setDeviceStream("127.0.0.1", &mockControlRelayPacketService_ListenForPacketRxServer{}, make(chan string))
	setDeviceStream("127.0.0.2", &mockControlRelayPacketService_ListenForPacketRxServer{}, make(chan string))
	setDeviceStream("127.0.0.3", &mockControlRelayPacketService_ListenForPacketRxServer{}, make(chan string))
	setDeviceStream("127.0.0.50", nil, make(chan string))
	for d := range devices {
		go func(ch chan string) {
			<-ch
		}(devices[d].ch)
	}
}
