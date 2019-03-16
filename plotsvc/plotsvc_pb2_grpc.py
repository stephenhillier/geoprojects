# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
import grpc

import plotsvc_pb2 as plotsvc__pb2


class SievePlotStub(object):
  # missing associated documentation comment in .proto file
  pass

  def __init__(self, channel):
    """Constructor.

    Args:
      channel: A grpc.Channel.
    """
    self.PlotSieve = channel.unary_unary(
        '/plotsvc.SievePlot/PlotSieve',
        request_serializer=plotsvc__pb2.SievePlotRequest.SerializeToString,
        response_deserializer=plotsvc__pb2.SievePlotResponse.FromString,
        )


class SievePlotServicer(object):
  # missing associated documentation comment in .proto file
  pass

  def PlotSieve(self, request, context):
    # missing associated documentation comment in .proto file
    pass
    context.set_code(grpc.StatusCode.UNIMPLEMENTED)
    context.set_details('Method not implemented!')
    raise NotImplementedError('Method not implemented!')


def add_SievePlotServicer_to_server(servicer, server):
  rpc_method_handlers = {
      'PlotSieve': grpc.unary_unary_rpc_method_handler(
          servicer.PlotSieve,
          request_deserializer=plotsvc__pb2.SievePlotRequest.FromString,
          response_serializer=plotsvc__pb2.SievePlotResponse.SerializeToString,
      ),
  }
  generic_handler = grpc.method_handlers_generic_handler(
      'plotsvc.SievePlot', rpc_method_handlers)
  server.add_generic_rpc_handlers((generic_handler,))
