# Go-SSH
Simple ssh client with go

## Execute command
```
client,err:=NewSSHClient("host", "root", "passwd")

stdout, stderr, err := client.ExecuteCmd("pwd")
```
## Execute command with private key
``` 
client, err:=NewSSHClientWithKey("host","root","/Users/user/.ssh/id_rsa")
stdout, stderr, err := client.ExecuteCmd("pwd")
```
## File upload

``` 
client.UploadFile("/tmp/test.html", "/root/test/test.html")
```
**Note**: Source file contents are fully read in memory, so you should not upload very large files using this command. If you really need to upload huge file to a lot of hosts, try using bittorrent or UFTP, as they provide much higher network effeciency than SSH.

