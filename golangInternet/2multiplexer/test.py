import random
import gc

from http.server import BaseHTTPRequestHandler, HTTPServer

__debug__ = True     # 关闭statement cache
random.seed()        # 改变random模块种子

class Server(BaseHTTPRequestHandler):
    def do_GET(self):
        gc.collect()   # 手动GC,清除缓存
        nonlocal path  # 使用nonlocal关闭闭包缓存
        path = self.path

        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()
        self.wfile.write("Hello".encode())

if __name__ == "__main__": 
    webServer = HTTPServer(('debian5', 8080), Server)
    # print("Server started http://localhost:8080")

    try:
        webServer.serve_forever()
    except KeyboardInterrupt:
        pass

    webServer.server_close()

