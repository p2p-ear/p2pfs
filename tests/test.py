import threading
import socket
import subprocess
from time import sleep, time

HOST = '127.0.0.1'

# because fuck nc utility

class P2PClient(threading.Thread):

    def __init__(self, port):
        self.stdout = None
        self.stderr = None
        self.port = port
        threading.Thread.__init__(self)

    def run(self):
        p = subprocess.Popen('./main {}'.format(self.port).split(),
                             shell=False,
                             stdout=subprocess.PIPE,
                             stderr=subprocess.PIPE)

        self.stdout, self.stderr = p.communicate()

    def flush(self):
        self.stdout = None
        self.stderr = None


def test_recieve():

    PORT = 9000
    iter = 10

    peer = P2PClient(PORT)
    peer.start()

    with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
        s.connect((HOST, PORT))

        for _ in range(10):
            payload = 'fuck'
            s.sendall(payload)
            print('Sent {}'.format(payload))
            print('Peer output {}'.format(peer.stdout))
            peer.flush()

        print("Closing connection")
        s.sendall("STOP")
        print('Peer output {}'.format(peer.stdout))
        peer.flush()

    peer.join()

    assert 0
