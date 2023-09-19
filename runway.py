from runwaylauncher import launch_query, executable
import socketserver



class MyTCPHandler(socketserver.BaseRequestHandler):
    """
    The request handler class for our server.

    It is instantiated once per connection to the server, and must
    override the handle() method to implement communication to the
    client.
    """

    def handle(self):
        # self.request is the TCP socket connected to the client
        args = self.request.recv(4096).strip().decode().split(" ")
        print(len(args))
        print("query_args recieved".format(self.client_address[0]))
        proc = launch_query(args, executable)

        # just send back the same data, but upper-cased
        self.request.sendall(bytes(args[0], 'utf-8'))


if __name__ == "__main__":
    HOST, PORT = "localhost", 9999

    # Create the server, binding to localhost on port 9999
    with socketserver.TCPServer((HOST, PORT), MyTCPHandler) as server:
        # Activate the server; this will keep running until you
        # interrupt the program with Ctrl-C
        server.serve_forever()
