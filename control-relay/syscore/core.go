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
* Control Relay core file
*
* Created by Filipe Claudio(Altice Labs) on 01/09/2020
 */

package syscore

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"plugin"
	"reflect"
	"strings"
	"sync"
	"time"

	control_relay "control_relay/pb/control_relay"
	tr477 "control_relay/pb/tr477"
	"control_relay/utils/log"

	"github.com/Juniper/go-netconf/netconf"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
)

var (
	OBBAA_ADDRESS      = os.Getenv("OBBAA_ADDRESS")
	OBBAA_PORT         = os.Getenv("OBBAA_PORT")
	CONTROL_RELAY_PORT = os.Getenv("CONTROL_RELAY_PORT")
	SDN_MC_SERVER_PORT = os.Getenv("SDN_MC_SERVER_PORT")
	SHARED_FOLDER      = os.Getenv("SHARED_FOLDER")
	PRIVATE_FOLDER     = os.Getenv("PRIVATE_FOLDER")
	STANDARD_FOLDER    = "./plugin-standard"
	GET_DEVICE_LIST    = "./netconf/get-device-list.xml"
)

var (
	grpcServer           = grpc.NewServer()
	sshConn, errSSH      = netconf.DialSSH("", &ssh.ClientConfig{})
	runningPlugins       = make(map[string]*RunningPluginsStruct)
	clientControllerList = make(map[string]*clientController)
	serverControllerList = make(map[string]*serverController)
	obbaaDeviceList      = make(map[string]string)
	recordPluginsNames   = []string{}
	mutex                = sync.RWMutex{}
)

// RunningPluginsStruct ...
type RunningPluginsStruct struct {
	plugin_name      string
	plugin_interface Plugin
	plugin_path      string
	plugin_vendor    string
	plugin_type      string
	plugin_model     string
	plugin_version   string
}

// Plugin ...
type Plugin interface {
	Start()
	PacketOutCallBack(packet *tr477.CpriMsg)
	Stop()
}

// Control App is the client and SDN M&C as gRPC Server
// CONTROL APP CLIENT  ----->  SERVER SDN CONTROLLER
// ######################################################
// helloService is responsable for sending a hello to the SDN server
// and also sending his own name, in this case the control relay
type serverController struct {
	dial         *grpc.ClientConn
	helloService tr477.CpriHelloClient
	clientStream tr477.CpriMessage_TransferCpriClient
	sdnAddress   string
}

func helloService(addHelloService tr477.CpriHelloClient, controller string) {
	_, err := addHelloService.HelloCpri(context.Background(), &tr477.HelloCpriRequest{
		LocalEndpointHello: &tr477.Hello{
			EntityName:   os.Getenv("CONTROL_RELAY_HELLO_NAME"),
			EndpointName: os.Getenv("CONTROL_RELAY_HELLO_NBI_ENDPOINT_NAME"),
		},
	},
	)

	if err != nil {
		log.Warning("Core: Hello Service failed")
		log.Warning("Core: Could not create a client connection to the given target: ")
		log.Warning("Core: Target is a server SDN controller: ", controller)
		log.Error("Core: Error: ", err)
	} else {
		log.Info("Core: Connection with Server SDN Controller: ok -->", controller)
	}
}

// waitForPacketsOnStream is responsable for receiving packages through the
// SDN server stream and fowarding to a plugin, invoking PacketOutCallBack
func waitForPacketsOnStream(stream tr477.CpriMessage_TransferCpriClient, controller string) {
	go on()
	for {
		packet, err := stream.Recv()
		if err == io.EOF {
			log.Error("Core: Error: ", err)
			break
		}
		if packet != nil {
			log.Debug("Core: Successfully received package")

			var plugin Plugin
			log.Info("Received packet from vnf: ", packet)
			if plugin = getPlugin(packet.MetaData.Generic.DeviceName); plugin == nil {
				go retryPacketOutCallBack(packet.MetaData.Generic.DeviceName, packet)
			} else {
				plugin.PacketOutCallBack(packet)
			}

		}
	}
}

func on() {
	for {
		log.Debug("Core: Waiting for packets on stream")
		time.Sleep(10 * time.Second)
	}
}

// PacketInCallBack is a default function to send packets for the SDN's
// This function takes only one argument, a packet from the device
func PacketInCallBack(packet *tr477.CpriMsg) {

	for _, controller := range clientControllerList {
		log.Debug(fmt.Sprintf("client controller: %+v", controller))
		if controller.stream == nil {
			log.Info("Null stream for controller " + controller.ip)
			continue
		}

		if !packetAllowedInFilter(*controller, packet) {
			log.Debug("Packet not allowed in filter for controller" + controller.ip)
			continue
		}

		log.Info("Received packet from olt: ", packet)
		if err := controller.stream.Send(packet); err != nil {
			log.Warning("Core: Failed to send the package")
			log.Warning("Core: DeviceName: ", packet.MetaData.Generic.DeviceName)
			log.Warning("Core: Error: ", err)
			if deleteClientController(controller.ip) {
				log.Info("Core: ClientController cleared from internal cache")
			}
		}
	}

	for _, controller := range serverControllerList {
		log.Debug(fmt.Sprintf("Server controller: %+v", controller))
		err := controller.clientStream.Send(packet)

		if err != nil {
			log.Warning("Core: PacketTx Service failed")
			log.Warning("Core: Could not create a client connection to the given target: ")
			log.Warning("Core: Target is the Server SDN Controller: ", controller.sdnAddress)
			log.Warning("Core: Error: ", err)
			if deleteServerController(controller.sdnAddress) {
				log.Warning("Core: ServerController cleared from list")
			}
		} else {
			log.Debug("Core: Packet sent successfully")
		}
	}
}

// packetAllowedInFilter is a function that is actually responsible for the
// checking if the packet matches the filter
func packetAllowedInFilter(controller clientController, packet *tr477.CpriMsg) bool {
	if controller.filter == nil {
		return true
	}
	var genMetadata *tr477.GenericMetadata
	if metadata := packet.GetMetaData(); metadata != nil {
		genMetadata = metadata.GetGeneric()
		if genMetadata == nil {
			return false
		}
	} else {
		return false
	}

	for _, filter := range controller.filter.Filter {
		if filter.Type == control_relay.ControlRelayPacketFilterList_ControlRelayPacketFilter_EXCLUDE {
			continue
		}
		if !doFiltering(filter.DeviceName, genMetadata.DeviceName) {
			continue
		}

		if !doFiltering(filter.DeviceInterface, genMetadata.DeviceInterface) {
			continue
		}

		if !doFiltering(filter.OriginatingRule, genMetadata.Direction.String()) {
			continue
		}
		return true
	}
	return false
}

/* --------------- Example: ---------------
   Packet:   0,   a,   b,   c
   ----------------------------------------
   Filter:   0,   a,   nil, c   -> true
   Filter:   1,   a,   nil, c   -> false
   Filter:   0,   a,   b,   c	-> true
   Filter:   0,   a,   b,   c	-> true
   Filter:   0,   a,   b,   d	-> false
   Filter:   0,   a,   b,   nil	-> true
*/
func doFiltering(filter string, prop string) bool {
	if filter == "" || filter == prop {
		return true
	}
	return false
}

// Control App is the server and SDN M&C as gRPC Client
// CONTROL APP SERVER  <-----  CLIENT SDN CONTROLLER
type clientController struct {
	stream  tr477.CpriMessage_TransferCpriServer
	filter  *control_relay.ControlRelayPacketFilterList
	ip      string
	network string
	ch      chan string
}

type controlRelayHelloService struct {
	tr477.UnimplementedCpriHelloServer
}

type controlRelayPacketService struct {
	tr477.UnimplementedCpriMessageServer
}

type controlRelayPacketFilterService struct {
	control_relay.UnimplementedControlRelayPacketFilterServiceServer
}

// Hello is one of the services provided by the proto, and is used mainly to
// establish and keep the connection between the Control Relay and SDN client.
// Is invoked by the SDN client, and if the connection is successfully established,
// the control relay stores the reference of that SDN for later sending packets
func (c *controlRelayHelloService) HelloCpri(ctx context.Context, in *tr477.HelloCpriRequest) (*tr477.HelloCpriResponse, error) {

	if ctx == nil {
		return nil, nil
	}

	log.Info("Core: A new Client SDN Controller is trying to connect... ")
	p, ok := peer.FromContext(ctx)
	if !ok {
		log.Warning("Core: Without information about the SDN")
	} else {
		if addClientController(p.Addr.String(), p.Addr.Network()) {
			log.Info("Core: A new Client SDN Controller are connected and he is waiting to receive packages: ")
			log.Info("Core: -----------> IP SDN: ", p.Addr.String())
			log.Info("Core: -----------> SDN Name: ", p.Addr.Network())
		}
	}

	return &tr477.HelloCpriResponse{
		RemoteEndpointHello: &tr477.Hello{
			EntityName:   os.Getenv("CONTROL_RELAY_HELLO_NAME"),
			EndpointName: os.Getenv("CONTROL_RELAY_HELLO_NBI_ENDPOINT_NAME"),
		},
	}, nil
}

func (*controlRelayPacketService) TransferCpri(stream tr477.CpriMessage_TransferCpriServer) error {
	ch := make(chan string)
	p, ok := peer.FromContext(stream.Context())
	if !ok {
		log.Warning("Error getting ip information from SDN")
		return errors.New("Error getting ip information from SDN")
	}
	setClientControllerStream(p.Addr.String(), stream, ch)
	log.Info("Transfer cpri core")
	for {
		in, err := stream.Recv()
		if err != nil {
			log.Error("Error receiving packet:", err)
			break
		}

		var plugin Plugin
		if plugin = getPlugin(in.MetaData.Generic.DeviceName); plugin == nil {
			go retryPacketOutCallBack(in.MetaData.Generic.DeviceName, in)
			continue
		}

		plugin.PacketOutCallBack(in)
	}
	return nil
}

var channelToAccessOBBAA = make(chan int, 1)

func retryPacketOutCallBack(deviceName string, packet *tr477.CpriMsg) {
	// ch := make(chan int, 1)  create a channel that supports only one goroutine at a time
	// ch <- 1     				will block if there is MAX ints in channel
	// <- ch					removes an int from channel, allowing another to proceed

	var plugin Plugin

	channelToAccessOBBAA <- 1 // Lock
	if plugin = getPlugin(deviceName); plugin != nil {
		<-channelToAccessOBBAA // Unlock
		plugin.PacketOutCallBack(packet)
		return
	}

	log.Warning("Updating the device list...")
	if startSSHConnection() {
		getListOfDevicesFromOBBAA()
	}

	if plugin = getPlugin(deviceName); plugin == nil {
		log.Warning("Core: There is no plugin for the ", deviceName, " device")
		<-channelToAccessOBBAA // Unlock
		return
	}
	<-channelToAccessOBBAA // Unlock
	plugin.PacketOutCallBack(packet)
}

func getPlugin(name string) Plugin {
	plugin := func() Plugin {
		if val, ok := obbaaDeviceList[name]; ok {
			if plugin, ok := runningPlugins[val]; ok {
				return plugin.plugin_interface
			}
		}
		return nil
	}()
	return plugin
}

func retryConnectionServerController(ip string) {
	for {
		time.Sleep(30 * time.Second)
		if startNorthboundClient(ip) {
			break
		}
	}
}

// addClientController is used to add new SDN's to the clientControllerList data structure
func addClientController(ip string, net string) bool {
	mutex.Lock()
	if clientControllerList[ip] != nil {
		log.Warning("Core: Client SDN controller already connected")
		mutex.Unlock()
		return false
	}

	// this key must be changed
	clientControllerList[ip] = &clientController{
		stream:  nil,
		filter:  nil,
		ip:      ip,
		network: net,
		ch:      nil,
	}
	mutex.Unlock()
	return true
}

func setClientControllerStream(ip string, stream tr477.CpriMessage_TransferCpriServer, ch chan string) bool {
	mutex.Lock()
	if controller, ok := clientControllerList[ip]; ok {
		controller.stream = stream
		controller.ch = ch
		mutex.Unlock()
		return true
	}
	mutex.Unlock()
	return false
}

func deleteServerController(ip string) bool {
	mutex.Lock()
	if controller, ok := serverControllerList[ip]; ok {
		controller.dial.Close()
		delete(serverControllerList, ip)
		go retryConnectionServerController(ip)
		mutex.Unlock()
		return true
	}
	mutex.Unlock()
	return false
}

func deleteClientController(ip string) bool {
	mutex.Lock()
	if controller, ok := clientControllerList[ip]; ok {
		controller.ch <- "close"
		delete(clientControllerList, ip)
		mutex.Unlock()
		return true
	}
	mutex.Unlock()
	return false
}

// startNorthboundServer starts the grpc server and register the services
// that are necessary for the normal operation of the application, in this
// case the services that are declared in the proto file
func startNorthboundServer() {

	lis, errLis := net.Listen("tcp", "0.0.0.0:"+CONTROL_RELAY_PORT)
	log.Info("Core: Initializing Northbound gRPC server on:", lis.Addr().String())
	if errLis != nil {
		log.Fatal("Core: Could not initialize Northbound gRPC server")
		log.Fatal("Core: Error:", errLis)
	}
	defer lis.Close()

	addHelloServiceServer := controlRelayHelloService{}
	addPacketServiceServer := controlRelayPacketService{}
	addFilterServiceServer := controlRelayPacketFilterService{}

	tr477.RegisterCpriHelloServer(grpcServer, &addHelloServiceServer)
	tr477.RegisterCpriMessageServer(grpcServer, &addPacketServiceServer)
	control_relay.RegisterControlRelayPacketFilterServiceServer(grpcServer, &addFilterServiceServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

// startNorthboundClient initiates a connection with the SDN server, the
// necessary services will be started for the normal functioning of the application
// and some of them will be immediately executed to establish a connection
func startNorthboundClient(sdnAddress string) bool {

	log.Info("Core: Initializing connection with the Northbound gRPC client")
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithBlock(),
	}

	connSteeringApp, errConnSteeringApp := grpc.Dial(sdnAddress, opts...)

	if errConnSteeringApp != nil {
		log.Fatal("Core: Failed to establish connection with the Steering App")
		log.Fatal("Core: Steering App address: ", sdnAddress)
		log.Fatal("Core: Error: ", errConnSteeringApp)
		return false
	}

	// Send hello for steering app
	addHelloServiceClient := tr477.NewCpriHelloClient(connSteeringApp)
	helloService(addHelloServiceClient, sdnAddress)

	addPacketServiceClient := tr477.NewCpriMessageClient(connSteeringApp)
	clientStream, err := addPacketServiceClient.TransferCpri(context.Background())
	if err != nil {
		log.Error("Error getting stream from client:", err)
	}
	go waitForPacketsOnStream(clientStream, sdnAddress)

	serverControllerList[sdnAddress] = &serverController{
		dial:         connSteeringApp,
		helloService: addHelloServiceClient,
		clientStream: clientStream,
		sdnAddress:   sdnAddress,
	}
	return true
}

func startSSHConnection() bool {
	log.Info("Core: Initializing ssh connection with OB-BAA")
	sshConfig := &ssh.ClientConfig{
		User:            os.Getenv("SSH_USER"),
		Auth:            []ssh.AuthMethod{ssh.Password(os.Getenv("SSH_PASSWORD"))},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	sshConn, errSSH = netconf.DialSSH(OBBAA_ADDRESS+":"+OBBAA_PORT, sshConfig)
	if errSSH != nil {
		log.Warning("Core: Failed to establish connection with OB-BAA")
		log.Warning("Core: OB-BAA address: ", OBBAA_ADDRESS, ":", OBBAA_PORT)
		log.Warning("Core: Error: ", errSSH)
		return false
	}

	log.Info("Core: SSH connection: ok")
	return true
}

// Devices ...
type Devices struct {
	Devices []Device `xml:"network-manager>managed-devices>device"`
}

// Device ...
type Device struct {
	DeviceName string `xml:"name"`
	Type       string `xml:"device-management>type"`
	Vendor     string `xml:"device-management>vendor"`
	Model      string `xml:"device-management>model"`
	Version    string `xml:"device-management>interface-version"`
}

func getListOfDevicesFromOBBAA() {
	reply, err := sshConn.Exec(netconf.RawMethod(readXMLfile()))
	if err != nil {
		log.Warning("Core: Could not execute a method")
		log.Warning("Core: Make sure the XML is correct")
		log.Warning("Core: Error: ", err)
		return
	}
	var devices Devices

	log.Debug("Reply from OB-BAA", reply)

	err = xml.Unmarshal([]byte(reply.Data), &devices)
	if err != nil {
		log.Warning("Core: Erro: ", err)
	}

	// Convert array to map
	log.Info("Core: List of devices: ")
	newListOfDevicesFromOBBAA := make(map[string]string)
	for _, d := range devices.Devices {
		newListOfDevicesFromOBBAA[d.DeviceName] = d.Vendor + "-" + d.Type + "-" + d.Model + "-" + d.Version
	}
	log.Info(newListOfDevicesFromOBBAA)

	// if the internal list is different from yhe new one, it will update
	if !reflect.DeepEqual(obbaaDeviceList, newListOfDevicesFromOBBAA) {
		mutex.Lock()
		obbaaDeviceList = newListOfDevicesFromOBBAA
		mutex.Unlock()
	}
}

func readXMLfile() string {
	if filepath.Ext(GET_DEVICE_LIST) != ".xml" {
		log.Warning("Core: Unknown file was found ", GET_DEVICE_LIST)
		log.Warning("Core: Only files with the extension \".xml\" are allowed in this folder")
		return ""
	}

	xmlFile, err := os.Open(GET_DEVICE_LIST)
	if err != nil {
		log.Warning("Core: Error opening XML file: ", err)
		return ""
	}
	defer xmlFile.Close()

	bytes, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Warning("Core: Error reading the xml file: ", err)
		return ""
	}

	return string(bytes)
}

func startPlugin(path string) error {

	log.Info("Core: starting plugin", path)
	// load module
	// 1. open the so file to load symbols, if exist's
	plug, err := plugin.Open(path)
	if err != nil {
		log.Error("Core: Error loading plugin:", path)
		log.Error("Core: Error: ", err)
		return err
	}

	// 2. look up a symbol (an exported function or variable)
	// in this case, variable Start and Stop
	symPlug, err := plug.Lookup("Plugin")
	if err != nil {
		log.Error("Core: Symbol not found: ", plug)
		log.Error("Core: Error: ", err)
		return err
	}

	// 3. Assert that loaded symbol is of a desired type
	// in this case interface type Plugin (defined above)
	var newplugin Plugin
	newplugin, ok := symPlug.(Plugin)
	if !ok {
		log.Error("Core: Unexpected type from module symbol")
		log.Error("Core: Error: ", err)
		return err
	}

	addPlugin(newplugin, path)
	return nil
}

func split(str string, symbol string) []string {
	return strings.Split(str, symbol)
}

func addPlugin(plug Plugin, path string) {
	go plug.Start()

	// split the path by the slash
	s := split(path, "/")
	// split the file extension ".so"
	name := strings.TrimSuffix(s[len(s)-1], filepath.Ext(s[len(s)-1]))
	// split the plugin name by the dashes
	props := []string(split(name, "-"))

	recordPluginsNames = append(recordPluginsNames, name)

	mutex.Lock()
	runningPlugins[name] = &RunningPluginsStruct{
		plugin_name:      name,
		plugin_interface: plug,
		plugin_path:      path,
		plugin_vendor:    props[0],
		plugin_type:      props[1],
		plugin_model:     props[2],
		plugin_version:   props[3],
	}
	mutex.Unlock()
}

func stopPlugin(p string) {
	log.Info("Stopping plugin:", p)
	name := strings.TrimSuffix(p, filepath.Ext(p))

	plugin := func() Plugin {
		for value, a := range runningPlugins {
			if value == name {
				p = a.plugin_path
				return a.plugin_interface
			}
		}
		return nil
	}()
	if plugin == nil {
		log.Warning("Plugin was not loaded")
		return
	}
	plugin.Stop() // Turn off the listener and stop grpcServer
	os.Remove(p)  // Remove the plugin from the folder
}

func getPluginsFromTheFolder(folder string) []string {
	var arrayOfLinkedPlugins []string

	filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if path == folder {
			return nil
		}
		if info.IsDir() {
			log.Warning("Core: A folder was found with the name ", path)
			log.Warning("Core: Only files with the extension \".so\" are allowed in this folder")
			return nil
		}
		if filepath.Ext(path) != ".so" {
			log.Warning("Core: Unknown file was found ", path)
			log.Warning("Core: Only files with the extension \".so\" are allowed in this folder")
			return nil
		}
		arrayOfLinkedPlugins = append(arrayOfLinkedPlugins, path)
		return nil
	})

	return arrayOfLinkedPlugins
}

// CopyFile copies a file from src to dst. If src and dst files exist, and are
// the same, then return success. Otherise, attempt to create a hard link
// between the two files. If that fail, copy the file contents from src to dst.
func CopyFile(src, dst string) (err error) {
	sfi, err := os.Stat(src)
	if err != nil {
		return
	}
	if !sfi.Mode().IsRegular() {
		// cannot copy non-regular files (e.g., directories,
		// symlinks, devices, etc.)
		return fmt.Errorf("Core: CopyFile: non-regular source file %s (%q)", sfi.Name(), sfi.Mode().String())
	}
	dfi, err := os.Stat(dst)
	if err != nil {
		if !os.IsNotExist(err) {
			return
		}
	} else {
		if !(dfi.Mode().IsRegular()) {
			return fmt.Errorf("Core: CopyFile: non-regular destination file %s (%q)", dfi.Name(), dfi.Mode().String())
		}
		if os.SameFile(sfi, dfi) {
			return
		}
	}
	if err = os.Link(src, dst); err == nil {
		return
	}
	err = copyFileContents(src, dst)
	return
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}

func equals(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		s := strings.Split(x, "/")
		x = s[len(s)-1]
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		s := strings.Split(x, "/")
		x = s[len(s)-1]
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

func checkPluginName(p string) bool {
	name := strings.TrimSuffix(p, filepath.Ext(p))

	for _, record := range recordPluginsNames {
		if record == name {
			log.Error("The name of the new plugin has already have been used, please choose another name!")
			os.Remove(SHARED_FOLDER + "/" + p)
			return false
		}
	}
	return true
}

// StartControlRelay ...
func StartControlRelay() {

	log.Info("############################## Starting Control Relay v1.0.0 ##############################")

	go startNorthboundServer()

	go func() {
		if startSSHConnection() {
			getListOfDevicesFromOBBAA()
		}
	}()

	log.Info("Starting client Northbound Connections")

	servers := split(os.Getenv("SDN_MC_SERVER_LIST"), ";")
	for _, address := range servers {
		sdnAddress := address + ":" + SDN_MC_SERVER_PORT
		go startNorthboundClient(sdnAddress)
	}

	log.Info("Checking plugins")

	for _, scr := range getPluginsFromTheFolder(SHARED_FOLDER) {

		s := strings.Split(scr, "/")
		x := s[len(s)-1]
		dst := PRIVATE_FOLDER + "/" + x

		err := CopyFile(scr, dst)
		if err != nil {
			log.Error("Core: Copy file failed: ", err)
		} else {
			log.Warning("Core: CopyFile succeeded")
		}
	}

	arrayOfStandardPlugins := getPluginsFromTheFolder(STANDARD_FOLDER)
	for _, path := range arrayOfStandardPlugins {
		startPlugin(path)
	}

	arrayOfLinkedPlugins := getPluginsFromTheFolder(PRIVATE_FOLDER)
	for _, path := range arrayOfLinkedPlugins {
		err := startPlugin(path)
		if err != nil {
			os.Remove(path)
		}
	}

	log.Info("Control Relay app running")
	for {
		time.Sleep(5 * time.Second)
		arrayOfNewPlugins := getPluginsFromTheFolder(SHARED_FOLDER)

		// doFiltering for plugins to stop
		if diff := equals(arrayOfLinkedPlugins, arrayOfNewPlugins); len(diff) != 0 {
			log.Warning("Core: Plugins to stop: ", diff)
			for _, p := range diff {
				go stopPlugin(p)
			}
		}

		// doFiltering for new plugins to upload
		if diff := equals(arrayOfNewPlugins, arrayOfLinkedPlugins); len(diff) != 0 {
			log.Warning("Core: Found new plugin: ", diff)
			for _, p := range diff {
				if checkPluginName(p) {
					CopyFile(SHARED_FOLDER+"/"+p, PRIVATE_FOLDER+"/"+p)
					go startPlugin(PRIVATE_FOLDER + "/" + p)
				}
			}
		}

		arrayOfLinkedPlugins = arrayOfNewPlugins
	}
}

func init() {
	if os.Getenv("CONTROL_RELAY_HELLO_NAME") == "" {
		os.Setenv("CONTROL_RELAY_HELLO_NAME", "control_relay_service")
	}

	if os.Getenv("CONTROL_RELAY_HELLO_NBI_ENDPOINT_NAME") == "" {
		os.Setenv("CONTROL_RELAY_HELLO_NBI_ENDPOINT_NAME", "control_relay_service_nbi")

	}

	if os.Getenv("PLUGIN_PORT") == "" {
		log.Info("PLUGIN_PORT environment variable was not specified, the value default is :50052")
		os.Setenv("PLUGIN_PORT", "50052")
	}

	if os.Getenv("SDN_MC_SERVER_PORT") == "" {
		SDN_MC_SERVER_PORT = "50053"
		log.Info("SDN_MC_SERVER_PORT environment variable was not specified, the value default is :" + SDN_MC_SERVER_PORT)
		//make sure the local variables and the env variables are in sync
		os.Setenv("SDN_MC_SERVER_PORT", SDN_MC_SERVER_PORT)
	}

	if os.Getenv("CONTROL_RELAY_PORT") == "" {
		CONTROL_RELAY_PORT = "50055"
		log.Info("CONTROL_RELAY_PORT environment variable was not specified, the value default is :" + CONTROL_RELAY_PORT)
		os.Setenv("CONTROL_RELAY_PORT", CONTROL_RELAY_PORT)
	}

	if os.Getenv("OBBAA_ADDRESS") == "" {
		OBBAA_ADDRESS = "192.168.56.102"
		log.Info("OBBAA_ADDRESS environment variable was not specidied, the value default is " + OBBAA_ADDRESS)
		os.Setenv("OBBAA_ADDRESS", OBBAA_ADDRESS)
	}

	if os.Getenv("OBBAA_PORT") == "" {
		OBBAA_PORT = "9292"
		log.Info("OBBAA_PORT environment variable was not specidied, the value default is :" + OBBAA_PORT)
		os.Setenv("OBBAA_PORT", "9292")
	}
	log.Info("OBBAA_PORT environment variable:", os.Getenv("OBBAA_PORT"))

	if os.Getenv("SSH_USER") == "" || os.Getenv("SSH_PASSWORD") == "" {
		log.Info("The SSH_USER or SSH_PASSWORD environment variables are not specified, default values have been set.")
		os.Setenv("SSH_USER", "admin")
		os.Setenv("SSH_PASSWORD", "password")
	}

	if os.Getenv("SHARED_FOLDER") == "" {
		SHARED_FOLDER = "./plugin-repo"
		os.Setenv("SHARED_FOLDER", SHARED_FOLDER)
	}
	log.Info("SHARED_FOLDER environment variable:", os.Getenv("SHARED_FOLDER"))

	if os.Getenv("PRIVATE_FOLDER") == "" {
		PRIVATE_FOLDER = "./plugin-enabled"
		os.Setenv("PRIVATE_FOLDER", PRIVATE_FOLDER)
	}
	log.Info("PRIVATE_FOLDER environment variable:", os.Getenv("PRIVATE_FOLDER"))
}
