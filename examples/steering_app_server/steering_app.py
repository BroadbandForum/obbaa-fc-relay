"""
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
"""

"""
 * Control Relay SDN controller simulator file
 *
 * Created by Joao Silva (Altice Labs) on 08/11/2019
""" 

from concurrent import futures
from dataclasses import dataclass
import logging
import sys
import time
import grpc
import os
from scapy import *
from ncclient import manager
from lxml import etree
from scapy.all import *
from google.protobuf import struct_pb2
import control_relay_service_pb2 as service
import control_relay_service_pb2_grpc as rpc



filter_olt = '<network-manager xmlns="urn:bbf:yang:obbaa:network-manager"><managed-devices><device><name>AlticeLabs_OLT</name><root><if:interfaces xmlns:if="urn:ietf:params:xml:ns:yang:ietf-interfaces" /></root></device></managed-devices></network-manager>'
devices = []
packets = []
count = 0

@dataclass
class ConnectedDevices:
    deviceName: str
        

with open("control-relay-v2/3-onu1-activation-to-bng1.xml") as f:
    activate_onu1 = etree.fromstring(f.read())

with open("control-relay-v2/3-onu2-activation-to-bng2.xml") as f:
    activate_onu2 = etree.fromstring(f.read())

    
#slot=2 port=1 onuId=1
#slot=2 port=2 onuId=2 

cmds = {
    "olt1":{ # Serial
        2: { # Slot
            1:{ # port
                1: activate_onu1 #onu
            },
            2:{
                2: activate_onu2 #onu
            }
        }
    }
}

class ControllerHelloService(rpc.ControlRelayHelloServiceServicer):
    
    obbaaAddress = "host"
    obbaaPort = "155"

    def __init__(self, obbaaAddress, obbaaPort):
        self.obbaaAddress = obbaaAddress
        self.obbaaPort = obbaaPort
        super(rpc.ControlRelayHelloServiceServicer, self).__init__()
        
    def Hello(self, request, context):
        devices.append(
            ConnectedDevices(request.device.device_name)
        )
        print("Message: OK")
    
        return service.HelloResponse()
           
class ControllerPacketService(rpc.ControlRelayPacketServiceServicer):

    obbaaAddress = "host"
    obbaaPort = "155"

    def __init__(self, obbaaAddress, obbaaPort):
        self.obbaaAddress = obbaaAddress
        self.obbaaPort = obbaaPort
        super(rpc.ControlRelayPacketServiceServicer, self).__init__()

    def PacketTx(self, request, context):
        """
        packet = Ether(request.packet)
        s = (
            '*** PACKET IN ***\n'
            'DeviceName: {deviceName}\n'
            'DeviceInterface: {deviceInterface}\n'
            'OriginatingRule: {originatingRule}\n'
            'Packet: \n{packet}\n'
        ).format(
            deviceName=request.device_name,
            deviceInterface=request.device_interface,
            originatingRule=request.originating_rule,
            packet=packet.show(dump=True)
        )
        
        logging.info(s)
        """
        global count
        count += 1
        print("Request number: ", count)
        """
        try:
            [request.oltPort][request.onuId]
            cmd = cmds[request.device_name][request.device_interface] 
        except KeyError as _:
            logging.error("Onu not found!")
            
        with manager.connect_ssh(
                self.obbaaAddress,
                port=self.obbaaPort,
                username="admin",
                password="password",
                hostkey_verify=False,
                look_for_keys=False,
                device_params={'name': 'default'}) as m:
            
            try:
                m.dispatch(cmd, filter=('subtree', filter_olt))
            except Exception as ex:
                logging.error(
                    "Failed to dispatch new configuration on OBBAA: {}".format(str(ex)))
        """
        
        return struct_pb2.Value()
        
    def ListenForPacketRx(self, request, context):
        a = service.ControlRelayPacket(
            device_name="AlticeLabs_OLT22",
            device_interface="1",
            originating_rule="",
            packet=None  
        )
        b = service.ControlRelayPacket(
            device_name="AlticeLabs_OLT222",
            device_interface="1",
            originating_rule="",
            packet=None  
        )
        c = service.ControlRelayPacket(
            device_name="AlticeLabs_OLT223",
            device_interface="2",
            originating_rule="",
            packet=None
        )
        d = service.ControlRelayPacket(
            device_name="AlticeLabs_OLT1994",
            device_interface="2",
            originating_rule="",
            packet=None
        )
        
        packets.append(a)
        packets.append(b)
        packets.append(c)
        packets.append(d)
        

        time.sleep(20)

        while True:
            for packet in packets:
                yield packet
                
           
        
def serve(port, obbaaAddress, obbaaPort):

    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    
    rpc.add_ControlRelayHelloServiceServicer_to_server(
        ControllerHelloService(obbaaAddress, obbaaPort), server)
    
    rpc.add_ControlRelayPacketServiceServicer_to_server(
        ControllerPacketService(obbaaAddress, obbaaPort), server)
    
    server.add_insecure_port('[::]:' + port)
    server.start()
 
    while True:
        try:
            time.sleep(1)
        except KeyboardInterrupt:
            print("Bye")
            sys.exit()

if __name__ == '__main__':
    logging.basicConfig(level=logging.INFO)
    
    obbaaAddress = os.environ['OBBAA_ADDRESS'] if 'OBBAA_ADDRESS' in os.environ else "10.112.106.236"
    obbaaPort = os.environ['OBBAA_PORT'] if 'OBBAA_PORT' in os.environ else "32000"
    port = os.environ['STEERING_APP_PORT'] if 'STEERING_APP_PORT' in os.environ else "50051"
    serve(port, obbaaAddress, obbaaPort)
