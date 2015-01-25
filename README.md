# goIRC

##IRC server written in Go


[IRC Spec](https://tools.ietf.org/html/rfc1459)

###Connection Steps
 1. telnet 192.168.128.213 3030
 1. JOIN #<Channel Name>:<User Name>
 1. MSG #<Channel Name>:<Message> 

### Local Steps
 1. go run bus.go connection.go help.go
 2. telnet localhost 3030
 3. PASS <your nick>
 4. JOIN #gophers:
 5. MSG #gophers:hello!