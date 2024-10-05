import socketserver
import threading

# Multi-Thread가 되는 TCP Server
class ThreadedTCPServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
    pass

# class MyTCPSocketHandler(socketserver.BaseRequestHandler):
class ThreadedTCPRequestHandler(socketserver.BaseRequestHandler):

    def handle(self):
        print('> client connected by IP address {0} with Port number {1}'.format(self.client_address[0], self.client_address[1]))
        while True:
            RecvData = self.request.recv(1024)
            # Asynchronous Socket Server
            cur_thread = threading.current_thread()
            print('> echoed:', RecvData.decode('utf-8'), 'by', cur_thread.name)
            self.request.sendall(RecvData)

            if RecvData.decode('utf-8') == 'quit':
                break


if __name__ == "__main__":
    HOST, PORT = '127.0.0.1', 65456
    print('> echo-server is activated')

    # with socketserver.TCPServer((HOST, PORT), MyTCPSocketHandler) as server:
    #     server.serve_forever()

    server = ThreadedTCPServer((HOST, PORT), ThreadedTCPRequestHandler)

    with server:
        ip, port = server.server_address

        server_thread = threading.Thread(target=server.serve_forever())

        server_thread.daemon = True
        server_thread.start()
        print('> server loop running in thread (main thread):', server_thread.name)

        baseThreadNumber = threading.active_count()
        while True:
            msg = input('> ')
            if msg == 'quit':
                if baseThreadNumber == threading.active_count():
                    print('> stop procedure started')
                    break
                else:
                    print('> active threads are remained :', threading.active_count() - baseThreadNumber, "threads")

        print('> echo-server is de-activated')
        server.shutdown()