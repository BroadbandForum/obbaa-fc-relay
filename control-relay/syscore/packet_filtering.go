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
* Control Relay packet filtering file
*
* Created by Filipe Claudio(Altice Labs) on 01/09/2020
*/ 

package syscore

import (
	"context"
	"control_relay/utils/log"
	pb "control_relay/pb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/peer"
)

func (c *controlRelayPacketFilterService) SetControlRelayPacketFilter(ctx context.Context, in *pb.ControlRelayPacketFilterList) (*pb.ControlRelayPacketFilterResponse, error) {
	
	if ctx == nil {
		return nil, nil
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.ControlRelayPacketFilterResponse{Response: pb.ControlRelayPacketFilterResponse_FAIL}, nil
	}
	address := p.Addr.String()

	mutex.Lock()
	if controller, ok := clientControllerList[address]; ok {
		if controller.filter != nil {
			c.ClearControlRelayPacketFilter(ctx, nil)
		}
		controller.filter = in
	}
	mutex.Unlock()
	return &pb.ControlRelayPacketFilterResponse{Response: pb.ControlRelayPacketFilterResponse_SUCCESS}, nil
}

func (c *controlRelayPacketFilterService) ClearControlRelayPacketFilter(ctx context.Context, _ *empty.Empty) (*pb.ControlRelayPacketFilterResponse, error) {

	if ctx == nil {
		return nil, nil
	}

	p, ok := peer.FromContext(ctx)
	if !ok {
		return &pb.ControlRelayPacketFilterResponse{Response: pb.ControlRelayPacketFilterResponse_FAIL}, nil
	}
	address := p.Addr.String()

	mutex.Lock()
	if controller, ok := clientControllerList[address]; ok {
		log.Warning(clientControllerList[address].ip, " cleared the filter list.")
		controller.filter = nil
	}
	mutex.Unlock()
	return &pb.ControlRelayPacketFilterResponse{Response: pb.ControlRelayPacketFilterResponse_SUCCESS}, nil
}