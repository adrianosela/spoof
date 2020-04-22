# smurf

carry out a smurf attack

![](.media/smurf_attack_diagram.png)

### Usage:

* modify `main.go` to specify the desired victim IP address, as well as the IP corresponding to the network's broadcast address.

* build this directory with `go build`

```
go build .
```

* run the program with privileged access:

```
sudo ./smurf
```
