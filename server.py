from http.server import HTTPServer, BaseHTTPRequestHandler
import json
import logging
import ssl


class Server:

    def __init__(self, config, bot_handler):

        class HttpHandler(BaseHTTPRequestHandler):

            def do_POST(self):
                content_len = int(self.headers.get('Content-Length'))
                post_body = self.rfile.read(content_len)

                data = json.loads(post_body)

                bot_handler(data)

                self.send_response(200)
                self.end_headers()

        self.server = HTTPServer((config['server_host'], config['server_port']), HttpHandler)

        if config['server_use_ssl']:
            sslctx = ssl.SSLContext(ssl.PROTOCOL_TLS_SERVER)
            sslctx.check_hostname = False
            sslctx.load_cert_chain(certfile='certificate.pem', keyfile="private.pem")
            self.server.socket = sslctx.wrap_socket(self.server.socket, server_side=True)

        logging.info("Server: inited")

    def start(self):
        logging.info("Server: Starting")
        self.server.serve_forever()
