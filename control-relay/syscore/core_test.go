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
* Control Relay core unit tests file
*
* Created by Filipe Claudio(Altice Labs) on 01/09/2020
 */

package syscore

import (
	"context"
	control_relay "control_relay/pb/control_relay"
	tr477 "control_relay/pb/tr477"
	"control_relay/utils/log"
	"net"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func boot(t *testing.T) {
	mockNewServer()

	if !addClientController("127.0.0.1:12345", "tcp") || // SUCCESS
		!addClientController("127.0.0.2:12345", "tcp") || // SUCCESS
		!addClientController("127.0.0.3:12345", "tcp") || // SUCCESS
		!addClientController("127.0.0.50:12345", "tcp") || // SUCCESS
		addClientController("127.0.0.1:12345", "tcp") { // FAIL
		t.Errorf("Error: creating client controller")
	}

	if !setClientControllerStream("127.0.0.1:12345", &mockControlRelayPacketService_ListenForPacketRxServer{}, make(chan string)) || // SUCCESS
		!setClientControllerStream("127.0.0.2:12345", &mockControlRelayPacketService_ListenForPacketRxServer{}, make(chan string)) || // SUCCESS
		!setClientControllerStream("127.0.0.3:12345", &mockControlRelayPacketService_ListenForPacketRxServer{}, make(chan string)) || // SUCCESS
		!setClientControllerStream("127.0.0.50:12345", nil, make(chan string)) || // SUCCESS
		setClientControllerStream("1234", nil, make(chan string)) { // FAIL
		t.Errorf("Error: setClientControllerStream")
	}

	for d := range clientControllerList {
		go func(ch chan string) {
			<-ch
		}(clientControllerList[d].ch)
	}
}

func Test_ClientServices(t *testing.T) {
	boot(t)
	c1 := mockNewControlRelayServiceClient(t, "0.0.0.0:12345")
	c2 := mockNewControlRelayServiceClient(t, "0.0.0.0:12345")
	c3 := mockNewControlRelayServiceClient(t, "0.0.0.0:12345")

	type args struct {
		s          serverController
		controller string
	}

	tests := []struct {
		name   string
		s      *serverController
		packet *tr477.CpriMsg
	}{
		{"test1", c1, MockControlRelayPacketInternalRequest("AlticeLabs_OLT22")},
		{"test2", c2, MockControlRelayPacketInternalRequest("AlticeLabs_OLT22")},
		{"test3", c1, MockControlRelayPacketInternalRequest("AlticeLabs_OLT22234234")},
		{"test4", c3, MockControlRelayPacketInternalRequest("AlticeLabs_OLT22234234")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helloService(tt.s.helloService, tt.s.sdnAddress)

			got, _ := tt.s.helloService.HelloCpri(context.Background(), MockHelloRequest("AlticeLabs_OLT22"))
			want := MockHelloResponse()
			if got.RemoteEndpointHello.EndpointName != want.RemoteEndpointHello.EndpointName ||
				got.RemoteEndpointHello.EntityName != want.RemoteEndpointHello.EntityName {
				t.Errorf("Hello() = %v, want %v", got, want)

			}

			go waitForPacketsOnStream(tt.s.clientStream, tt.s.sdnAddress)
		})
	}
	time.Sleep(1 * time.Second)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PacketInCallBack(tt.packet)
			deleteServerController(tt.s.sdnAddress)
		})
	}
}

func Test_equals(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{"test1", args{[]string{"a", "b", "c"}, []string{"a", "b", "c"}}, []string{}, false},
		{"test2", args{[]string{"a", "b", "c"}, []string{"d", "b", "c"}}, []string{"a"}, false},
		{"test3", args{[]string{}, []string{"a", "b", "c"}}, []string{}, false},
		{"test4", args{[]string{"a", "b", "c"}, []string{}}, []string{"a", "b", "c"}, false},
		{"test5", args{[]string{"a", "b", "c"}, []string{"a", "b", "c"}}, []string{"a"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := equals(tt.args.a, tt.args.b)
			for i, v := range got {
				if v != tt.want[i] {
					if tt.wantErr {
						return
					}
					t.Errorf("equals() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

/*
func TestCopyFile(t *testing.T) {
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test1", args{"aaa", "bbb"}, true},
		{"test2", args{"./bin/plugins/BBF-OLT-standard-1.0.so", "./test/test.so"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CopyFile(tt.args.src, tt.args.dst)
			if err != nil {
				if tt.wantErr {
					return
				}
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
} */

func Test_split(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		symbol  string
		want    []string
		wantErr bool
	}{
		{"test1", "a/e/i/o/u", "/", []string{"a", "e", "i", "o", "u"}, false},
		{"test2", "a/ei/o/u", "/", []string{"a", "ei", "o", "u"}, false},
		{"test3", "a", "/", []string{"b"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := split(tt.str, tt.symbol); !reflect.DeepEqual(got, tt.want) {
				if tt.wantErr {
					return
				}
				t.Errorf("split() = %v, want %v", got, tt.want)
			}
		})
	}
}

func mockNewServer() {
	lis, _ := net.Listen("tcp", "0.0.0.0:12345")
	server := grpc.NewServer()

	addHelloServ := controlRelayHelloService{}
	addPacketServ := controlRelayPacketService{}
	addFilterServ := controlRelayPacketFilterService{}

	// The & symbol points to the address of the stored value.
	tr477.RegisterCpriHelloServer(server, &addHelloServ)
	tr477.RegisterCpriMessageServer(server, &addPacketServ)
	control_relay.RegisterControlRelayPacketFilterServiceServer(server, &addFilterServ)

	go server.Serve(lis)
}

func mockNewControlRelayServiceClient(t *testing.T, address string) *serverController {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	require.NoError(t, err)
	addHelloServiceClient := tr477.NewCpriHelloClient(conn)
	addPacketServiceClient := tr477.NewCpriMessageClient(conn)
	stream, err := addPacketServiceClient.TransferCpri(context.Background())
	if err != nil {
		log.Error("")
		return nil
	}

	s := &serverController{
		dial:         conn,
		helloService: addHelloServiceClient,
		clientStream: stream,
		sdnAddress:   address,
	}
	serverControllerList[address] = s
	return s
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

// MockHelloRequest ...
func MockHelloRequest(name string) *tr477.HelloCpriRequest {
	return &tr477.HelloCpriRequest{
		LocalEndpointHello: &tr477.Hello{
			EntityName:   os.Getenv("CONTROL_RELAY_HELLO_NAME"),
			EndpointName: os.Getenv("CONTROL_RELAY_HELLO_NBI_ENDPOINT_NAME"),
		},
	}
}

// MockHelloResponse ...
func MockHelloResponse() *tr477.HelloCpriResponse {
	return &tr477.HelloCpriResponse{
		RemoteEndpointHello: &tr477.Hello{
			EntityName:   os.Getenv("CONTROL_RELAY_HELLO_NAME"),
			EndpointName: os.Getenv("CONTROL_RELAY_HELLO_NBI_ENDPOINT_NAME"),
		},
	}
}

// MockControlRelayPacketRequest ...
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
