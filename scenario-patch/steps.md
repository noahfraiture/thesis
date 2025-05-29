Files belong to "admin" but can be read by any user including the "web" which execute them

1. Enter in the game, add secret as player, and then add as friend
2. Decode that credentials 'admin' 'admin'
3. Find /admin and use the credentials
4. Use parameter in url to see files such at /home/web/.ssh/authorized_keys
5. Search that public key online and find the related private key in a repol
6. ssh as web in this server
7. use the cronjob to put our ssh key in authorized to enter as admin
8. as admin, use cve in systemd to elevate to root and restart server

0. At any point : analyse logs to see the problem

## Sql injection

- with test we can easily create a player in the table players by guessing it exist, it only requires 'name' which is present on the table
- [injections](https://github.com/swisskyrepo/PayloadsAllTheThings/blob/master/SQL%20Injection/SQLite%20Injection.md*)
- We can list tables with `select * from sqlite_master where type = 'table'` and get value from there
- To see values we must insert a player and add him as friend `jimmy') ('noah', 'PLACEHOLDER')` but that would require to understand that ? **TODO**
- `jimmy') ON CONFLICT(p1, p2) DO NOTHING; INSERT INTO players(name) VALUES ((SELECT pass FROM secrets LIMIT 1)); INSERT INTO friends (p1, p2) VALUES ('noah', (SELECT pass FROM secrets LIMIT 1)); --`
