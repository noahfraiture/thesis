import http.server
import socketserver
import sys


# Function to print the POST payload
def print_payload(payload):
    print("Received POST payload:")
    print(payload.decode("utf-8"))


# Create a custom request handler
class SimpleHTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        # Handle GET requests
        self.send_response(200)
        self.send_header("Content-type", "text/html")
        self.end_headers()
        self.wfile.write(b"Hello, this is a simple HTTP server!")

    def do_POST(self):
        # Handle POST requests
        content_length = int(self.headers["Content-Length"])
        post_data = self.rfile.read(content_length)

        # Print the POST payload
        print_payload(post_data)

        # Send a simple response
        self.send_response(200)
        self.send_header("Content-type", "text/plain")
        self.end_headers()
        self.wfile.write(b"POST request received and processed")


# Check if the correct number of command-line arguments is provided
if len(sys.argv) != 2:
    print("Usage: python script_name.py <interface>")
    sys.exit(1)

# Get the interface from command-line argument
interface = sys.argv[1]
port = 80

# Create the server
with socketserver.TCPServer((interface, port), SimpleHTTPRequestHandler) as httpd:
    print(f"Serving on {interface}:{port}")
    # Start the server
    httpd.serve_forever()
