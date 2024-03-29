# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import control_relay_service_pb2 as control__relay__service__pb2


class ControlRelayHelloServiceStub(object):
  # missing associated documentation comment in .proto file
  pass

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
  # missing associated documentation comment in .proto file
  pass

  def Hello(self, request, context):
    # missing associated documentation comment in .proto file
    pass
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
