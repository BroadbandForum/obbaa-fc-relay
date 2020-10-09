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
* Control Relay file that decode and print the contents of the packet
*
* Created by Filipe Claudio(Altice Labs) on 01/09/2020
*/ 

package decode

import (
	"control_relay/utils/log"
	"strings"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

// LayerMap ...
var LayerMap = make(map[string]Layer)

// Layer ...
type Layer struct {
	callback func(layer gopacket.Layer)
}

// PrintPacketInfo ...
func PrintPacketInfo(p []byte) {

	packet := gopacket.NewPacket(p, layers.LayerTypeEthernet, gopacket.Default)

    for _, layer := range packet.Layers() {
		if layerType, ok := LayerMap[layer.LayerType().String()]; ok {
			layerType.callback(layer)
		}
	}

	applicationLayer := packet.ApplicationLayer()
    if applicationLayer != nil {
        log.Info("Application layer/Payload found.")
        log.Info("%s\n", applicationLayer.Payload())

        // Search for a string inside the payload
        if strings.Contains(string(applicationLayer.Payload()), "HTTP") {
            log.Info("HTTP found!")
        }
    }

    // Check for errors
    if err := packet.ErrorLayer(); err != nil {
        log.Info("Error decoding some part of the packet:", err)
    } 
}

func layerTypeEthernet(layer gopacket.Layer) {
	log.Warning("-------------------------- Ethernet Layer detected --------------------------")
	ethernetPacket, _ := layer.(*layers.Ethernet)
	log.Info("Source MAC: ", ethernetPacket.SrcMAC)
	log.Info("Destination MAC: ", ethernetPacket.DstMAC)
	log.Info("Ethernet type: ", ethernetPacket.EthernetType)
}

func layerTypeIPv4(layer gopacket.Layer) {
	ip, _ := layer.(*layers.IPv4)
	log.Warning("-------------------------- IPv4 Layer detected --------------------------")
	log.Info("From %s to %s\n", ip.SrcIP, ip.DstIP)
	log.Info("Protocol: ", ip.Protocol)
	log.Info("ID: ", ip.Id)
	log.Info("Version: ", ip.Version)
	log.Info("Length: ", ip.Length)
	log.Info("IHL: ", ip.IHL)
	log.Info("TOS: ", ip.TOS)
	log.Info("TTL: ", ip.TTL)
	log.Info("Chksum: ", ip.Checksum)
	log.Info("Frag: ", ip.FragOffset)
	log.Info("Flags: ", ip.Flags)
	log.Info("Options: ", ip.Options)
}

func layerTypeIPv6(layer gopacket.Layer) {
	ip, _ := layer.(*layers.IPv6)
	log.Warning("-------------------------- IPv6 Layer detected --------------------------")
	log.Info("From %s to %s\n", ip.SrcIP, ip.DstIP)
	log.Info("Protocol: ", ip.NextHeader)
	log.Info("Version: ", ip.Version)
	log.Info("Length: ", ip.Length)
}

func layerTypeTCP(layer gopacket.Layer) {
	log.Warning("-------------------------- TCP Layer detected --------------------------")
	tcp, _ := layer.(*layers.TCP)
	log.Info("From port %d to %d\n", tcp.SrcPort, tcp.DstPort)
	log.Info("Sequence number: ", tcp.Seq)
}

func layerTypeUDP(layer gopacket.Layer) {
	log.Warning("-------------------------- UDP Layer detected --------------------------")
	udp, _ := layer.(*layers.UDP)
	log.Info("From port %d to %d\n", udp.SrcPort, udp.DstPort)
	log.Info("Checksum: ", udp.Checksum)
}

func layerTypeDHCPv4(layer gopacket.Layer) {
	log.Warning("-------------------------- DHCPv4 Layer detected ------------------------")
	dhcpv4, _ := layer.(*layers.DHCPv4)
	log.Info("IP Client: ", dhcpv4.ClientIP)
	log.Info("Address Client: ", dhcpv4.ClientHWAddr)
	log.Info("Flags: ", dhcpv4.Flags)
	log.Info("Operation: ", dhcpv4.Operation)
	log.Info("DHCP Options: ", dhcpv4.Options)
}

func layerTypeDot1Q(layer gopacket.Layer) {
	log.Warning("-------------------------- Dot1Q Layer detected --------------------------")
	dot1q, _ := layer.(*layers.Dot1Q)
	log.Info("Priority: ", dot1q.Priority)
	log.Info("VLAN :", dot1q.VLANIdentifier)
	log.Info("Type :", dot1q.Type)
}

func layerTypeDNS(layer gopacket.Layer) {
	log.Warning("-------------------------- DNS Layer detected --------------------------")
	dns, _ := layer.(*layers.DNS)
	log.Info("ANCount: ", dns.ANCount)
}

func init() {
	LayerMap[layers.LayerTypeEthernet.String()] 	= Layer{callback: layerTypeEthernet}
	LayerMap[layers.LayerTypeIPv4.String()] 		= Layer{callback: layerTypeIPv4}
	LayerMap[layers.LayerTypeIPv6.String()] 		= Layer{callback: layerTypeIPv6}
	LayerMap[layers.LayerTypeDot1Q.String()] 		= Layer{callback: layerTypeDot1Q}
	LayerMap[layers.LayerTypeDHCPv4.String()] 		= Layer{callback: layerTypeDHCPv4}
	LayerMap[layers.LayerTypeUDP.String()] 			= Layer{callback: layerTypeUDP}
	LayerMap[layers.LayerTypeDNS.String()]			= Layer{callback: layerTypeDNS}
}