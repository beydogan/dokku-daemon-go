# dokku-daemon-go
Dokku Daemon written with Go to interact with Dokku

# Requirements 

A server running Ubuntu 14.04 or later with Dokku installed.

# Installing 

As a user with access to `sudo`

```
curl -L https://github.com/beydogan/dokku-daemon-go/releases/download/v0.1.1/dokku-daemon-go-linux64
chmod +x dokku-daemon-go-linux64
sudo ./dokku-daemon-go-linux64 install
sudo service dokku-daemon start
```

# Usage

Daemon will create a UNIX socket at `/var/run/dokku-daemon/dokku-daemon.sock`.

To test you can use;

```
sudo socat - UNIX-CONNECT:/var/run/dokku-daemon/dokku-daemon.sock
```

and type `apps`. It should return the output in JSON format.
