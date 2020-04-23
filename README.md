# smurf

[![Go Report Card](https://goreportcard.com/badge/github.com/adrianosela/smurf)](https://goreportcard.com/report/github.com/adrianosela/smurf)
[![Documentation](https://godoc.org/github.com/adrianosela/smurf?status.svg)](https://godoc.org/github.com/adrianosela/smurf)
[![license](https://img.shields.io/github/license/adrianosela/smurf.svg)](https://github.com/adrianosela/smurf/blob/master/LICENSE)

carry out a smurf attack

![](.media/smurf_attack_diagram.png)

### Usage:

* modify `main.go` to specify the desired victim IP address, as well as the IP corresponding to the network's broadcast address.

* build this directory with `go build`

```
$ go build
```

* run the program with privileged access:

```
$ sudo ./smurf
```

### Wireshark:

> example using my own host as the victim, with my LAN's broadcast as the network address:

![](.media/run.gif)
