# Leaning Hydrometer Manager for Raspberry Pi

[![go-build-and-test](https://github.com/fetherolfjd/leaning-hydro-pi/actions/workflows/go-build-and-test.yml/badge.svg)](https://github.com/fetherolfjd/leaning-hydro-pi/actions/workflows/go-build-and-test.yml)

## Build

### Compiling From Non Raspberry Pi

For Raspberry Pi 3 B:

`env GOOS=linux GOARCH=arm GOARM=7 go build`

For Raspberry Pi Zero W:

`env GOOS=linux GOARCH=arm GOARM=6 go build`

## OS Configuration

It's probable that a standard user will not have access to the bluetooth device(s)
without `root` or `sudo`. To overcome this, we'll add some do `setcap` on the binary:

`sudo setcap 'cap_net_raw,cap_net_admin+eip' /absolute/path/to/leaning-hydro-pi`

### Dependencies

This requires bluetooth packages:

`sudo apt-get install bluetooth`

## Tested Devices

So far this has been developed and tested with the following Raspberry Pi models:

 - Raspberry Pi 3 Model B
 - Raspberry Pi Zero W

## Resources

### Tilt Hydrometer iBeacon Data Format

Example data message:

```
04 3E 27 02 01 00 00 5A 09 9B 16 A3 04 1B 1A FF 4C 00 02 15 A4 95 BB 10 C5 B1 4B 44 B5 12 13 70 F0 2D 74 DE 00 44 03 F8 C5 C7
```

Message breakdown:

```
04: HCI Packet Type HCI Event
3E: LE Meta event
27: Parameter total length (39 octets)
02: LE Advertising report sub-event
01: Number of reports (1)
00: Event type connectable and scannable undirected advertising
00: Public address type
5A: address
09: address
9B: address
16: address
A3: address
04: address
1B: length of data field (27 octets)
1A: length of first advertising data (AD) structure (26)
FF: type of first AD structure - manufacturer specific data
4C: manufacturer ID - Apple iBeacon
00: manufacturer ID - Apple iBeacon
02: type (constant, defined by iBeacon spec)
15: length (constant, defined by iBeacon spec)
A4: device UUID
95: device UUID
BB: device UUID
10: device UUID
C5: device UUID
B1: device UUID
4B: device UUID
44: device UUID
B5: device UUID
12: device UUID
13: device UUID
70: device UUID
F0: device UUID
2D: device UUID
74: device UUID
DE: device UUID
00: major - temperature (in degrees fahrenheit)
44: major - temperature (in degress fahrenheit)
03: minor - specific gravity (x1000)
F8: minor - specific gravity (x1000)
C5: TX power in dBm
C7: RSSI in dBm
```

| Color | Address |
| ----- | ------- |
| Red |    A495BB10C5B14B44B5121370F02D74DE |
| Green |  A495BB20C5B14B44B5121370F02D74DE |
| Black |  A495BB30C5B14B44B5121370F02D74DE |
| Purple | A495BB40C5B14B44B5121370F02D74DE |
| Orange | A495BB50C5B14B44B5121370F02D74DE |
| Blue |   A495BB60C5B14B44B5121370F02D74DE |
| Yellow | A495BB70C5B14B44B5121370F02D74DE |
| Pink |   A495BB80C5B14B44B5121370F02D74DE |


### Links

 - [Tilt Hydrometer iBeacon Data Format](https://kvurd.com/blog/tilt-hydrometer-ibeacon-data-format/)
