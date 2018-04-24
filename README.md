# exec-docker
[![Build Status](https://travis-ci.org/wenvlem/exec-docker.svg?branch=master)](https://travis-ci.org/wenvlem/exec-docker)
[![GoDoc](https://godoc.org/github.com/nanopack/portal?status.svg)](https://godoc.org/github.com/wenvlem/exec-docker)

exec-docker is a simple binary to be used with telegraf's `exec` plugin. It outputs in influx line protocol. It collects and reports several image, volume, and container metrics.

Metrics gathered include:
 - Containers:
   - Total created
   - Total in each "state"
   - Total disk consumed
 - Images:
   - Total created
   - Total unused
   - Total dangling
   - Total disk consumed
 - Volumes:
   - Total created
   - Total unused
   - Total disk consumed

(Metrics gathered may not be applicable to all users)

Errors are printed to stderr.

#### Example Chronograf Dashboard:
![chronograf](assets/chron.png?raw=true "chronograf")

#### Example Output:
```
containers total=5,size_rw=7661352,running=5
images total=4,dangling=0,unused=0,size=744879901
volumes,volume=064986933f size=32840
volumes,volume=07041bc47a size=65536
volumes,volume=90c22d20b6 size=0
volumes,volume=total total=3,unused=0,size=98376
```

#### Example telegraf config:
```
[[inputs.exec]]
  commands = ["exec-docker"]
  timeout = "5s"
  name_suffix = "_docker"
  data_format = "influx"
```
