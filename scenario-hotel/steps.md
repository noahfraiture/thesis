1. Find camera ip and port to find hikvision running
2. find version of hikvision on '/' page
3. Exploit to get black screenshop, and config with credentials of operator
4. Use these credentials of operator to enter in dashboard server with ssh
5. On the server, find token in ENV and route /next in .bash_history
6. Make mosquitto crash to turn on the light
7. exploit c file to control the robot and make it go further

light
can use `mosquitto_pub -h localhost -t light/control -m "off"` to sned message
