# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import control_relay_service_pb2 as control__relay__service__pb2
from google.protobuf import empty_pb2 as google_dot_protobuf_dot_empty__pb2


class ControlRelayHelloServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.Hello = channel.unary_unary(
                '/control_relay_service.v1.ControlRelayHelloService/Hello',
                request_serializer=control__relay__service__pb2.HelloRequest.SerializeToString,
                response_deserializer=control__relay__service__pb2.HelloResponse.FromString,
                )


class ControlRelayHelloServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def Hello(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ControlRelayHelloServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'Hello': grpc.unary_unary_rpc_method_handler(
                    servicer.Hello,
                    request_deserializer=control__relay__service__pb2.HelloRequest.FromString,
                    response_serializer=control__relay__service__pb2.HelloResponse.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'control_relay_service.v1.ControlRelayHelloService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class ControlRelayHelloService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def Hello(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/control_relay_service.v1.ControlRelayHelloService/Hello',
            control__relay__service__pb2.HelloRequest.SerializeToString,
            control__relay__service__pb2.HelloResponse.FromString,
            options, channel_credentials,
            call_credentials, compression, wait_for_ready, timeout, metadata)


class ControlRelayPacketServiceStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.PacketTx = channel.unary_unary(
                '/control_relay_service.v1.ControlRelayPacketService/PacketTx',
                request_serializer=control__relay__service__pb2.ControlRelayPacket.SerializeToString,
                response_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                )
        self.ListenForPacketRx = channel.unary_stream(
                '/control_relay_service.v1.ControlRelayPacketService/ListenForPacketRx',
                request_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
                response_deserializer=control__relay__service__pb2.ControlRelayPacket.FromString,
                )


class ControlRelayPacketServiceServicer(object):
    """Missing associated documentation comment in .proto file."""

    def PacketTx(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListenForPacketRx(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_ControlRelayPacketServiceServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'PacketTx': grpc.unary_unary_rpc_method_handler(
                    servicer.PacketTx,
                    request_deserializer=control__relay__service__pb2.ControlRelayPacket.FromString,
                    response_serializer=google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            ),
            'ListenForPacketRx': grpc.unary_stream_rpc_method_handler(
                    servicer.ListenForPacketRx,
                    request_deserializer=google_dot_protobuf_dot_empty__pb2.Empty.FromString,
                    response_serializer=control__relay__service__pb2.ControlRelayPacket.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'control_relay_service.v1.ControlRelayPacketService', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class ControlRelayPacketService(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def PacketTx(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/control_relay_service.v1.ControlRelayPacketService/PacketTx',
            control__relay__service__pb2.ControlRelayPacket.SerializeToString,
            google_dot_protobuf_dot_empty__pb2.Empty.FromString,
            options, channel_credentials,
            call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListenForPacketRx(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(request, target, '/control_relay_service.v1.ControlRelayPacketService/ListenForPacketRx',
            google_dot_protobuf_dot_empty__pb2.Empty.SerializeToString,
            control__relay__service__pb2.ControlRelayPacket.FromString,
            options, channel_credentials,
            call_credentials, compression, wait_for_ready, timeout, metadata)
