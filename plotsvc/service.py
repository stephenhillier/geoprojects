import grpc
import plotsvc_pb2
import plotsvc_pb2_grpc

import base64
from concurrent import futures
from io import BytesIO
import matplotlib
import numpy as np
import matplotlib.pyplot as plt
import time

_ONE_DAY_IN_SECONDS = 60 * 60 * 24


def plot_sieve_to_png(points):
    plt.figure(figsize=(20, 12))
    plt.semilogx((*zip(*points)))
    ticks = [0.08, 0.16, 0.32, 1.0, 2.5, 5.0, 10, 20, 50, 100]
    plt.xticks(ticks, [str(t) for t in ticks])
    plt.ylim((0, 100))
    plt.grid(b=None, which='major', axis='both')
    plt.xlabel('Opening size (mm)')
    plt.ylabel('Percent passing')
    fig_bytes = BytesIO()
    plt.savefig(fig_bytes, format='png')
    fig_bytes.seek(0)
    fig_base64 = base64.b64encode(fig_bytes.read())
    return fig_base64.decode("utf-8")


class SievePlotServicer(plotsvc_pb2_grpc.SievePlotServicer):
    """ provides methods for using the SievePlot service """

    def PlotSieve(self, request, context):
        pts = [[0.08, 10], [0.16, 20], [0.32, 30], [1.25, 40], [
            2, 50], [5, 60], [10, 70], [12, 80], [16, 90], [20, 100]]
        fig = plot_sieve_to_png(pts)
        print(fig)
        return plotsvc_pb2.SievePlotResponse(
            ok=True,
            figure=fig
        )


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    plotsvc_pb2_grpc.add_SievePlotServicer_to_server(
        SievePlotServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    try:
        while True:
            time.sleep(_ONE_DAY_IN_SECONDS)
    except KeyboardInterrupt:
        server.stop(0)


if __name__ == '__main__':
    serve()
