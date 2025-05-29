# Context

The network of a company is infected by a botnet but they call you simply for a dataleak on a machine. You have to find the problem.
The botnet is behaving like this :
- process running entierly in memory
- message poll from parent : other infected machine or external server. If parent down : parent = external server
- packet have following form |TYPE|HOP|PAYLOAD|. The packet are encrypted
  TYPE :    Message itself
    SLEEP :   parent -> child. Put on sleep
    SPREAD :  parent -> child. Infect a machine
    EXEC :    parent -> child. Exec bash command
    QUIT :    parent -> child. Leave the botnet
    ASK :     child -> parent. Polling, and the server answer with one of the aboves
  HOP :     Number of hop to prograpage, so put this message in our process with HOP-1
  PAYLOAD : Additional payload

- When infected, install a backdoor
- Process intercept signal and ignore them



# Steps

1. You can search binary on the disk but won't find anything
2. You can list process but zheng is disguised
3. Find the process
  - Find process with ps and top by network usage as the name is disguised
  - Can confirm by dumping memory with pmem or Volatility  
  - Can see parent in memory, either machine of the network, or external server if the user already block some machine
    Else he will need to dump packets
  - find used port with netstat

4. Dump packet to read message and reverse engineer protocol
5. Find foothold
  - can be by killing process, blocking backdoor and see what come
  - can be by analyzing logs (/var/log/auth.log)
5. Has to block it in two place :
  - initial foothold in ssh
  - backdoor ftp ?
6. Instead of killing process everywhere, send message to shutdown network
