sockbench
=========

Sockbench is a simple and compact socket testing utility. To measure
performance. You can experiment with Go language and analyze your server.

The source code of the utility can be used for self-training.

Only the standard Go library is used. Demonstrates all the fantastic power of
Goroutine. It shows how to start many servers and clients in one process, how to
intercept system signals SIGTERM (), how to bypass the list of servers and
correctly complete the work.

You can safely use the keyboard shortcut Ctrl + C in the Linux console. Except
crazy mode.

It shows how to safely use variables for a list of running servers (sync.Mutex).
How to wait for goroutine (threads) to execute using sync.WaitGroup.

If you do not understand deeply what is happening inside the operating system,
it is recommended to read man 2 socket.

The Benchmark suite demonstrates what a socket buffer is and how the size of the
buffer affects performance. Also shows the difference between a file socket and
a TCP socket. This code will allow you to better understand why sometimes the
speed is up to 725.46 Kbit / sec, and sometimes the speed can reach 20.20 Gbit /
sec. This is a dramatic difference of 30,000 times.

It is not recommended to use this code for Production systems.

Macintosh compatibility has not been tested. If you slightly correct the code
regarding socket paths, then most likely it will work.

On Windows Unix, the socket will most likely not work. But don’t worry. You can
experiment with TCP / IP.




## Compatibility

Utility tested with version:

* runtime.Version(): **go1.12.6**
* linux uname -a
```
Linux s **5.1.16**-1-default #1 SMP Wed Jul 3 12:37:47 UTC 2019 (2af8a22) x86_64 x86_64 x86_64 GNU/Linux.
```
* lsb_release -a
```
    LSB Version:    n/a
    Distributor ID: openSUSE
    Description:    **openSUSE Tumbleweed**
    Release:        **20190721**
    Codename:       n/a
```
* lscpu
```
    Architecture:        x86_64
    CPU op-mode(s):      32-bit, 64-bit
    Byte Order:          Little Endian
    CPU(s):              6
    On-line CPU(s) list: 0-5
    Thread(s) per core:  1
    Core(s) per socket:  6
    Socket(s):           1
    NUMA node(s):        1
    Vendor ID:           GenuineIntel
    Model name:          **Intel(R) Core(TM) i7-8850H CPU @ 2.60GHz**
    CPU max MHz:         4300.0000
    L3 cache:            9216K
    NUMA node0 CPU(s):   0-5
```



## Quick Start

Get source code:

```Bash
$ git clone git@github.com:Clark-Devis/Sockbench.git
```

For testing, you need to specify your address, example:

```Go
  var ifaceIpAddr string = "192.168.10.40"
```

For testing, you need to specify a function that interests you, just comment out, example:
```Go
func main() {
  // EXAMPLE1 - Multiple connections on single server:
  Example1()
}
```

```Bash
$ time go run sockbench.go
```


## Example 1

Multiple connections on single server.

Socket type: "unix".

Runtime: 1 proccess (main go) -> multithreading (go routines). Each connection in its own goroutine.

You can use the linux command "nc" to connect multiple times.
Shutdown processing of multiple connection is performed correctly.
You can safely use the keyboard shortcut Ctrl + C in the Linux console.

NOTE from manual Nc command: "After the file has been transferred, the connection will close automatically."

```Go
func main() {
  // EXAMPLE1 - Multiple connections on single server:
  Example1()
}
```

Use the linux command "nc" to connect:

```Bash
$ pid=`pgrep sockbench`; path="/run/user/1000/sockbench/1/socket_${pid}_1"; nc -U $path
```

Output:
```
[2019-08-20T11:01:50.7260601440][1][INFO][SERVER] Uid: 1000; Dir: /run/user/1000/sockbench/1; Socket: /run/user/1000/sockbench/1/socket_7264_1; Addr: /run/user/1000/sockbench/1/socket_7264_1
[2019-08-20T11:01:50.7261272610][1][INFO][SERVER] Dir was created: /run/user/1000/sockbench/1
[2019-08-20T11:01:50.7261488830][1][INFO][SERVER] Visited files or directory: /run/user/1000/sockbench/1
[2019-08-20T11:01:50.7261867550][1][INFO][SERVER] Listener was created: &{0xc0000da000 /run/user/1000/sockbench/1/socket_7264_1 true {{0 0} 0}}
[2019-08-20T11:01:50.7261907730][1][INFO][SERVER] SetUnlinkOnClose(true) executed.
[2019-08-20T11:01:50.7261943870][1][INFO][SERVER] The action AcceptUnix() will be executed and go into standby mode.
[2019-08-20T11:01:50.7261973440][1][INFO][SERVER][STANDBY]
[2019-08-20T11:01:54.1714497870][1][INFO][SERVER][CONN=&{{0xc0000da200}}] The action AcceptUnix() for unix socket was executed.
[2019-08-20T11:01:54.1714615470][1][INFO][SERVER][CONN=&{{0xc0000da200}}] SetReadBuffer(1) executed
[2019-08-20T11:01:54.1714690080][1][INFO][SERVER] The action AcceptUnix() will be executed and go into standby mode.
[2019-08-20T11:01:54.1714714540][1][INFO][SERVER][STANDBY]
[2019-08-20T11:02:03.4108496960][1][INFO][SERVER][CONN=&{{0xc0000da200}}] Read of 1 bytes: "1"
[2019-08-20T11:02:03.4108990450][1][INFO][SERVER][CONN=&{{0xc0000da200}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:03.8267567920][1][INFO][SERVER][CONN=&{{0xc0000da200}}] Read of 1 bytes: "2"
[2019-08-20T11:02:03.8268016800][1][INFO][SERVER][CONN=&{{0xc0000da200}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:04.1468912950][1][INFO][SERVER][CONN=&{{0xc0000da200}}] Read of 1 bytes: "3"
[2019-08-20T11:02:04.1469350090][1][INFO][SERVER][CONN=&{{0xc0000da200}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:08.9461621400][1][INFO][SERVER][CONN=&{{0xc000104080}}] The action AcceptUnix() for unix socket was executed.
[2019-08-20T11:02:08.9465091900][1][INFO][SERVER][CONN=&{{0xc000104080}}] SetReadBuffer(1) executed
[2019-08-20T11:02:08.9466232600][1][INFO][SERVER] The action AcceptUnix() will be executed and go into standby mode.
[2019-08-20T11:02:08.9466581400][1][INFO][SERVER][STANDBY]
[2019-08-20T11:02:09.4671027680][1][INFO][SERVER][CONN=&{{0xc000104080}}] Read of 1 bytes: "a"
[2019-08-20T11:02:09.4671448830][1][INFO][SERVER][CONN=&{{0xc000104080}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:09.7709520150][1][INFO][SERVER][CONN=&{{0xc000104080}}] Read of 1 bytes: "b"
[2019-08-20T11:02:09.7709919010][1][INFO][SERVER][CONN=&{{0xc000104080}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:10.9883056500][1][INFO][SERVER][CONN=&{{0xc000104080}}] Read of 1 bytes: "c"
[2019-08-20T11:02:10.9887348100][1][INFO][SERVER][CONN=&{{0xc000104080}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:37.1263499400][1][INFO][SERVER][CONN=&{{0xc0000da200}}][DISCONNECT] EOF
[2019-08-20T11:02:40.1711930480][1][INFO][SERVER][CONN=&{{0xc000104080}}][DISCONNECT] EOF
[2019-08-20T11:02:44.9655460070][1][INFO][SERVER][CONN=&{{0xc0000da400}}] The action AcceptUnix() for unix socket was executed.
[2019-08-20T11:02:44.9655578720][1][INFO][SERVER][CONN=&{{0xc0000da400}}] SetReadBuffer(1) executed
[2019-08-20T11:02:44.9655638830][1][INFO][SERVER] The action AcceptUnix() will be executed and go into standby mode.
[2019-08-20T11:02:44.9655664300][1][INFO][SERVER][STANDBY]
[2019-08-20T11:02:47.6454187520][1][INFO][SERVER][CONN=&{{0xc000104200}}] The action AcceptUnix() for unix socket was executed.
[2019-08-20T11:02:47.6454305330][1][INFO][SERVER][CONN=&{{0xc000104200}}] SetReadBuffer(1) executed
[2019-08-20T11:02:47.6454648710][1][INFO][SERVER] The action AcceptUnix() will be executed and go into standby mode.
[2019-08-20T11:02:47.6454679680][1][INFO][SERVER][STANDBY]
[2019-08-20T11:02:49.6035253860][1][INFO][SERVER][CONN=&{{0xc0000da400}}] Read of 1 bytes: "4"
[2019-08-20T11:02:49.6035660070][1][INFO][SERVER][CONN=&{{0xc0000da400}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:50.1953435600][1][INFO][SERVER][CONN=&{{0xc0000da400}}] Read of 1 bytes: "5"
[2019-08-20T11:02:50.1957628500][1][INFO][SERVER][CONN=&{{0xc0000da400}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:50.8285889400][1][INFO][SERVER][CONN=&{{0xc0000da400}}] Read of 1 bytes: "6"
[2019-08-20T11:02:50.8286296560][1][INFO][SERVER][CONN=&{{0xc0000da400}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:54.7712522410][1][INFO][SERVER][CONN=&{{0xc000104200}}] Read of 1 bytes: "d"
[2019-08-20T11:02:54.7712948660][1][INFO][SERVER][CONN=&{{0xc000104200}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:55.2356403460][1][INFO][SERVER][CONN=&{{0xc000104200}}] Read of 1 bytes: "e"
[2019-08-20T11:02:55.2356817300][1][INFO][SERVER][CONN=&{{0xc000104200}}] Read of 1 bytes: "\n"
[2019-08-20T11:02:55.6432210500][1][INFO][SERVER][CONN=&{{0xc000104200}}] Read of 1 bytes: "f"
[2019-08-20T11:02:55.6432619670][1][INFO][SERVER][CONN=&{{0xc000104200}}] Read of 1 bytes: "\n"
^C[2019-08-20T11:02:59.4756524540][0][PRIM][MAIN][INTERCEPTOR] Signal is intercepted by channel: interrupt
[2019-08-20T11:02:59.4756894810][0][PRIM][MAIN][INTERCEPTOR] len(sl.servers): 1
[2019-08-20T11:02:59.4757043040][0][PRIM][MAIN][INTERCEPTOR]: Clean up server id = 1 with type = unix
[2019-08-20T11:02:59.4757872940][1][INFO][SERVER] Socket file was removed: /run/user/1000/sockbench/1/socket_7264_1
[2019-08-20T11:02:59.4758305600][1][INFO][SERVER] Socket dir was removed: /run/user/1000/sockbench/1
exit status 1
```
*A similar output to the terminal for the other examples. But a detailed output is disabled for benchmark.*


## Example 2

*A similar output to the terminal for the Example1.*

Local unix socket file: **func Example2Unix()**

For loopback interface with Port = 6080: **func Example2TCP("127.0.0.1")**

Local interface with IP-address with Port = 6080: **func Example2TCP(ifaceIpAddr)**


```Go
func main() {
  Example2Unix()
}
```

```Go
func main() {
  Example2TCP("127.0.0.1")
}
```

For testing, you need to specify your address, example:

```Go
var ifaceIpAddr string = "192.168.10.40"
```

```Go
func main() {
  Example2TCP(ifaceIpAddr)
}
```

Runtime: 1 proccess (main go) -> multithreading (go routines). Single server in
its own goroutine (1). Single client in its own goroutine (1).

The client does not close the connection. The server closes the connection after
receiving data and receiving a signal from the client.



## Example 3

*A similar output to the terminal for the Example1.*

Local unix socket file. Two client and two server in parallel mode with parallel start.

Runtime: 1 proccess (main go) -> multithreading (go routines).

Goroutines:

1. First single server in its own goroutine.
2. First single client in its own goroutine.
3. Secondary single server in its own goroutine.
4. Secondary single client in its own goroutine.

The client does not close the connection. The server closes the connection after
receiving data and receiving a signal from the client.

```Go
func main() {
  Example3UnixParallel()
}
```



## Example Crazy 2

**Next using crazy mode. Be careful.**

Cycle launch of a huge variety client/server in parallel by batching.

Local unix socket file. Lots of client and lots of server in parallel mode.

Runtime: 1 proccess (main go) -> multithreading (go routines).

Goroutines:

1. First single server in its own goroutine.
2. First single client in its own goroutine.
3. Secondary single server in its own goroutine.
4. Secondary single client in its own goroutine.
5. ...
6. ...
7. N single server in its own goroutine.
8. N single client in its own goroutine.

The client does not close the connection. The server closes the connection after
receiving data and receiving a signal from the client.

Batching go routines by 1 pair:

```Go
func main() {
  ExampleCrazy2(1)
}
```

Output:

```
[2019-08-20T15:12:19.2214989490][0][PRIM][MAIN] x=0, y=196
[2019-08-20T15:12:19.2215039600][0][PRIM][MAIN] x=0, y=197
[2019-08-20T15:12:19.2215087000][0][PRIM][MAIN] x=0, y=198
[2019-08-20T15:12:19.2215155050][0][PRIM][MAIN] x=0, y=199
[2019-08-20T15:12:19.2215209430][0][PRIM][MAIN] x=0, y=200
[2019-08-20T15:12:19.2215281930][0][PRIM][MAIN] runtime.NumGoroutine(): 403
```

Batching go routines by 1000 pair, 1 pair it is 2 goroutine:

```Go
func main() {
  ExampleCrazy2(1000)
}
```

Output:
```
[2019-08-20T15:23:12.3191320340][0][PRIM][MAIN] x=1000, y=196
[2019-08-20T15:23:12.3191353130][0][PRIM][MAIN] x=1000, y=197
[2019-08-20T15:23:12.3191387840][0][PRIM][MAIN] x=1000, y=198
[2019-08-20T15:23:12.3191420820][0][PRIM][MAIN] x=1000, y=199
[2019-08-20T15:23:12.3191464190][0][PRIM][MAIN] x=1000, y=200
[2019-08-20T15:23:12.3191498120][0][PRIM][MAIN] runtime.NumGoroutine(): 403
```



## Example Bench Unix 1

Benchmark. Unix socket. Single client and single server in parallel.

Quantitites of iterations: 100

Buffer size in bytes: 1, 4096, 16384, 4194304, 6291456, 62914560

Runtime: 1 proccess (main go) -> multithreading (go routines).

Goroutines:

1. First single server in its own goroutine.
2. First single client in its own goroutine.
3. Result.
4. Secondary single server in its own goroutine.
5. Secondary single client in its own goroutine.
6. Result.
6. ...
7. N single server in its own goroutine.
8. N single client in its own goroutine.
9. Result for N pair.

The client does not close the connection. The server closes the connection after
receiving data and receiving a signal from the client.

Batching go routines by 1 pair.

Output:

```
[2019-08-20T15:31:34.6028334200][0][PRIM][MAIN][BENCH] Comment: FS=tmpfs(RAM); Type: Unix socket; Speed=970.65 KByte/s; Bytes=200; Iterations=100; Buffer size=1; Time: 0.20121899999999998ms 201.219microsec 201219ns
[2019-08-20T15:31:34.6059535750][0][PRIM][MAIN][BENCH] Comment: FS=tmpfs(RAM); Type: Unix socket; Speed=560.31 MByte/s; Bytes=409700; Iterations=100; Buffer size=4096; Time: 0.697331ms 697.331microsec 697331ns
[2019-08-20T15:31:34.6135672010][0][PRIM][MAIN][BENCH] Comment: FS=tmpfs(RAM); Type: Unix socket; Speed=313.33 MByte/s; Bytes=1638500; Iterations=100; Buffer size=16384; Time: 4.987067ms 4987.067microsec 4987067ns
[2019-08-20T15:31:35.1572411570][0][PRIM][MAIN][BENCH] Comment: FS=tmpfs(RAM); Type: Unix socket; Speed=771.09 MByte/s; Bytes=419430500; Iterations=100; Buffer size=4194304; Time: 518.743239ms 518743.239microsec 518743239ns
[2019-08-20T15:31:35.6801575780][0][PRIM][MAIN][BENCH] Comment: FS=tmpfs(RAM); Type: Unix socket; Speed=1.19 GByte/s; Bytes=629145700; Iterations=100; Buffer size=6291456; Time: 493.07163199999997ms 493071.632microsec 493071632ns
[2019-08-20T15:31:38.1115598330][0][PRIM][MAIN][BENCH] Comment: FS=tmpfs(RAM); Type: Unix socket; Speed=2.51 GByte/s; Bytes=6291456100; Iterations=100; Buffer size=62914560; Time: 2330.48072ms 2.33048072e+06microsec 2330480720ns
```



## Example Bench Suite 1

**Bench suite may take a long time.**

Benchmark. Unix/TCPloopback/TCPlocal socket. Single client and single server in parallel.

Output:

```
time go run sockbench.go
..............................
---------------------------------------------------------------------------------------------------------------------------------------------
|   # | Comment+Type:                                      |    Speed (bits): |   Speed (bytes): |  It: |         Buffer size: |      Time: |
---------------------------------------------------------------------------------------------------------------------------------------------
|   1 | FS=tmpfs(RAM), Unix socket                         |   7951.25 Kbit/s |   993.91 KByte/s |  100 |                    1 |   0.00 sec |
|   2 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |   2570.92 Kbit/s |   321.37 KByte/s |  100 |                    1 |   0.00 sec |
|   3 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |   2566.63 Kbit/s |   320.83 KByte/s |  100 |                    1 |   0.00 sec |
|   4 | FS=tmpfs(RAM), Unix socket                         |   3356.28 Kbit/s |   419.53 KByte/s |  100 |                    2 |   0.00 sec |
|   5 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |   1536.67 Kbit/s |   192.08 KByte/s |  100 |                    2 |   0.00 sec |
|   6 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |   1004.78 Kbit/s |   125.60 KByte/s |  100 |                    2 |   0.00 sec |
|   7 | FS=tmpfs(RAM), Unix socket                         |   5169.42 Kbit/s |   646.18 KByte/s |  100 |                    4 |   0.00 sec |
|   8 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |   1636.60 Kbit/s |   204.57 KByte/s |  100 |                    4 |   0.00 sec |
|   9 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |   1646.44 Kbit/s |   205.80 KByte/s |  100 |                    4 |   0.00 sec |
|  10 | FS=tmpfs(RAM), Unix socket                         |     10.71 Mbit/s |     1.34 MByte/s |  100 |                   16 |   0.00 sec |
|  11 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |   4496.56 Kbit/s |   562.07 KByte/s |  100 |                   16 |   0.00 sec |
|  12 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |   4528.92 Kbit/s |   566.11 KByte/s |  100 |                   16 |   0.00 sec |
|  13 | FS=tmpfs(RAM), Unix socket                         |    148.98 Mbit/s |    18.62 MByte/s |  100 |                  256 |   0.00 sec |
|  14 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |    115.93 Mbit/s |    14.49 MByte/s |  100 |                  256 |   0.00 sec |
|  15 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |     65.28 Mbit/s |     8.16 MByte/s |  100 |                  256 |   0.00 sec |
|  16 | FS=tmpfs(RAM), Unix socket                         |   1163.18 Mbit/s |   145.40 MByte/s |  100 |                 4096 |   0.00 sec |
|  17 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |    184.54 Kbit/s |    23.07 KByte/s |  100 |                 4096 |  17.34 sec |
|  18 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |    185.13 Kbit/s |    23.14 KByte/s |  100 |                 4096 |  17.29 sec |
|  19 | FS=tmpfs(RAM), Unix socket                         |   3271.75 Mbit/s |   408.97 MByte/s |  100 |                16384 |   0.00 sec |
|  20 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |    647.21 Kbit/s |    80.90 KByte/s |  100 |                16384 |  19.78 sec |
|  21 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |    653.86 Kbit/s |    81.73 KByte/s |  100 |                16384 |  19.58 sec |
|  22 | FS=tmpfs(RAM), Unix socket                         |     16.79 Gbit/s |     2.10 GByte/s |  100 |              4194304 |   0.19 sec |
|  23 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |     11.83 Gbit/s |     1.48 GByte/s |  100 |              4194304 |   0.26 sec |
|  24 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |     12.01 Gbit/s |     1.50 GByte/s |  100 |              4194304 |   0.26 sec |
|  25 | FS=tmpfs(RAM), Unix socket                         |     13.29 Gbit/s |     1.66 GByte/s |  100 |              6291456 |   0.35 sec |
|  26 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |     13.10 Gbit/s |     1.64 GByte/s |  100 |              6291456 |   0.36 sec |
|  27 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |     14.87 Gbit/s |     1.86 GByte/s |  100 |              6291456 |   0.32 sec |
|  28 | FS=tmpfs(RAM), Unix socket                         |     20.20 Gbit/s |     2.52 GByte/s |  100 |             62914560 |   2.32 sec |
|  29 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |     18.00 Gbit/s |     2.25 GByte/s |  100 |             62914560 |   2.60 sec |
|  30 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |     18.22 Gbit/s |     2.28 GByte/s |  100 |             62914560 |   2.57 sec |
---------------------------------------------------------------------------------------------------------------------------------------------

real    1m23.778s
user    0m9.029s
sys     0m6.769s
```



## Example Bench Suite 2

**Bench suite may take a long time.**

Benchmark. Unix socket. Single client and single server in parallel.

Output:

```
time go run sockbench.go
..........
---------------------------------------------------------------------------------------------------------------------------------------------
|   # | Comment+Type:                                      |    Speed (bits): |   Speed (bytes): |  It: |         Buffer size: |      Time: |
---------------------------------------------------------------------------------------------------------------------------------------------
|   1 | FS=tmpfs(RAM), Unix socket                         |     13.90 Mbit/s |     1.74 MByte/s |  100 |                    1 |   0.00 sec |
|   2 | FS=tmpfs(RAM), Unix socket                         |     11.22 Mbit/s |     1.40 MByte/s |  100 |                    2 |   0.00 sec |
|   3 | FS=tmpfs(RAM), Unix socket                         |     18.56 Mbit/s |     2.32 MByte/s |  100 |                    4 |   0.00 sec |
|   4 | FS=tmpfs(RAM), Unix socket                         |     39.28 Mbit/s |     4.91 MByte/s |  100 |                   16 |   0.00 sec |
|   5 | FS=tmpfs(RAM), Unix socket                         |    542.28 Mbit/s |    67.78 MByte/s |  100 |                  256 |   0.00 sec |
|   6 | FS=tmpfs(RAM), Unix socket                         |    992.07 Mbit/s |   124.01 MByte/s |  100 |                 4096 |   0.00 sec |
|   7 | FS=tmpfs(RAM), Unix socket                         |   2319.80 Mbit/s |   289.98 MByte/s |  100 |                16384 |   0.01 sec |
|   8 | FS=tmpfs(RAM), Unix socket                         |   7593.57 Mbit/s |   949.20 MByte/s |  100 |              4194304 |   0.42 sec |
|   9 | FS=tmpfs(RAM), Unix socket                         |     14.01 Gbit/s |     1.75 GByte/s |  100 |              6291456 |   0.33 sec |
|  10 | FS=tmpfs(RAM), Unix socket                         |     19.65 Gbit/s |     2.46 GByte/s |  100 |             62914560 |   2.39 sec |
---------------------------------------------------------------------------------------------------------------------------------------------

real    0m3.482s
user    0m3.407s
sys     0m1.634s
```

This Bench suite demonstrates fluctuations in the network subsystem. It is also
important to understand that the results suite duration of 0.00 sec should be
treated with caution.



## Example Bench Suite 3

**Bench suite may take a long time.**

Benchmark. TCP loopback socket. Single client and single server in parallel.

Output:

```
time go run sockbench.go
..........
---------------------------------------------------------------------------------------------------------------------------------------------
|   # | Comment+Type:                                      |    Speed (bits): |   Speed (bytes): |  It: |         Buffer size: |      Time: |
---------------------------------------------------------------------------------------------------------------------------------------------
|   1 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |   2668.72 Kbit/s |   333.59 KByte/s |  100 |                    1 |   0.00 sec |
|   2 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |   3581.53 Kbit/s |   447.69 KByte/s |  100 |                    2 |   0.00 sec |
|   3 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |   6570.00 Kbit/s |   821.25 KByte/s |  100 |                    4 |   0.00 sec |
|   4 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |     16.08 Mbit/s |     2.01 MByte/s |  100 |                   16 |   0.00 sec |
|   5 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |     14.34 Mbit/s |     1.79 MByte/s |  100 |                  256 |   0.01 sec |
|   6 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |    186.47 Kbit/s |    23.31 KByte/s |  100 |                 4096 |  17.16 sec |
|   7 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |    639.87 Kbit/s |    79.98 KByte/s |  100 |                16384 |  20.01 sec |
|   8 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |     13.37 Gbit/s |     1.67 GByte/s |  100 |              4194304 |   0.23 sec |
|   9 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |     14.06 Gbit/s |     1.76 GByte/s |  100 |              6291456 |   0.33 sec |
|  10 | FS=tmpfs(RAM), TCP4 socket 127.0.0.1               |     16.50 Gbit/s |     2.06 GByte/s |  100 |             62914560 |   2.84 sec |
---------------------------------------------------------------------------------------------------------------------------------------------

real    0m40.911s
user    0m3.627s
sys     0m2.580s
```


## Example Bench Suite 4

**Bench suite may take a long time.**

Benchmark. TCP local IP address. Single client and single server in parallel.

Output:

```
time go run sockbench.go
..........
---------------------------------------------------------------------------------------------------------------------------------------------
|   # | Comment+Type:                                      |    Speed (bits): |   Speed (bytes): |  It: |         Buffer size: |      Time: |
---------------------------------------------------------------------------------------------------------------------------------------------
|   1 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |    725.46 Kbit/s |    90.68 KByte/s |  100 |                    1 |   0.00 sec |
|   2 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |   1099.28 Kbit/s |   137.41 KByte/s |  100 |                    2 |   0.00 sec |
|   3 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |   1816.77 Kbit/s |   227.10 KByte/s |  100 |                    4 |   0.00 sec |
|   4 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |   4782.85 Kbit/s |   597.86 KByte/s |  100 |                   16 |   0.00 sec |
|   5 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |     69.86 Mbit/s |     8.73 MByte/s |  100 |                  256 |   0.00 sec |
|   6 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |    182.96 Kbit/s |    22.87 KByte/s |  100 |                 4096 |  17.49 sec |
|   7 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |    667.89 Kbit/s |    83.49 KByte/s |  100 |                16384 |  19.17 sec |
|   8 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |     15.31 Gbit/s |     1.91 GByte/s |  100 |              4194304 |   0.20 sec |
|   9 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |     14.35 Gbit/s |     1.79 GByte/s |  100 |              6291456 |   0.33 sec |
|  10 | FS=tmpfs(RAM), TCP4 socket 192.168.10.40           |     17.34 Gbit/s |     2.17 GByte/s |  100 |             62914560 |   2.70 sec |
---------------------------------------------------------------------------------------------------------------------------------------------

real    0m40.216s
user    0m3.427s
sys     0m2.563s
```