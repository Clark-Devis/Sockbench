package main

import (
  "bufio"
  "fmt"
  "net"
  "os"
  "os/signal"
  "os/user"
  "path"
  "path/filepath"
  "runtime"
  "strconv"
  "strings"
  "sync"
  "syscall"
  "time"
  // "io"
  // "github.com/davecgh/go-spew/spew"
)

func init() {
  if runtime.GOOS != "linux" {
    panic("Delete the line \"runtime.GOOS != \"linux\"\" if you want to run in a different environment. Be careful.")
  }

  // if runtime.Version() != "go1.12.6" {
  //   panic("Delete the line \"runtime.Version() != \"go1.12.6\"\" if you want to run in a different environment. Be careful.")
  // }
}

var printPrim   bool   = true
var printInfo   bool   = true
var ifaceIpAddr string = "192.168.10.40"


func main() {
  // prim(0, fmt.Sprintf("[MAIN] os.Getpid():            %30v", os.Getpid()))
  // prim(0, fmt.Sprintf("[MAIN] os.Getppid():           %30v", os.Getppid()))
  // prim(0, fmt.Sprintf("[MAIN] runtime.GOOS:           %30v", runtime.GOOS))
  // prim(0, fmt.Sprintf("[MAIN] runtime.Version:        %30v", runtime.Version()))
  // prim(0, fmt.Sprintf("[MAIN] runtime.NumCPU():       %30v", runtime.NumCPU()))
  // prim(0, fmt.Sprintf("[MAIN] runtime.NumGoroutine(): %30v", runtime.NumGoroutine()))
  ngBefore := runtime.NumGoroutine()

  // EXAMPLE1 - Multiple connections on single server:
  // Example1()

  // EXAMPLE2 - Single client and single server in parallel goroutine mode:
  // fmt.Println("Example2Unix()")
  // ngBefore++
  // Example2Unix()
  // fmt.Println("Example2TCP(\"127.0.0.1\")")
  // ngBefore++
  // Example2TCP("127.0.0.1") // is Loopback; PORT=6080
  // fmt.Println("Example2TCP(ifaceIpAddr)")
  // ngBefore++
  // Example2TCP(ifaceIpAddr) // is not Loopback; Local inet address; PORT=6080

  // EXAMPLE3 - Two client and two server in parallel mode with parallel start:
  // fmtPrintln("Example3UnixParallel()")
  // ngBefore++
  // Example3UnixParallel()

  // EXAMPLECRAZY1 - Cycle launch of a single client/server in parallel:
  // fmt.Println("ExampleCrazy1()")
  // ngBefore++
  // ExampleCrazy1()

  // EXAMPLECRAZY2 - Cycle launch of a huge variety client/server in parallel:
  // Next using crazy mode. Be careful.
  //
  // fmt.Println("ExampleCrazy2(1)")
  // ngBefore++
  // ExampleCrazy2(1)
  // fmt.Println("ExampleCrazy2(1000)")
  // ngBefore++
  // ExampleCrazy2(1000)

  // EXAMPLEBENCHUNIX1 - Benchmark. Unix socket. Single client and single server in parallel.
  // ngBefore++
  // ExampleBenchUnix1()

  // EXAMPLEBENCHSUITE1 - Benchmark. Unix/TCPloopback/TCPlocal socket. Single client and single server in parallel.
  // ngBefore++
  // ExampleBenchSuite1()

  // EXAMPLEBENCHSUITE2 - Benchmark. Unix socket. Single client and single server in parallel.
  // ngBefore++
  // ExampleBenchSuite2()

  // EXAMPLEBENCHSUITE3 - Benchmark. TCP loopback socket. Single client and single server in parallel.
  // ngBefore++
  // ExampleBenchSuite3()

  // EXAMPLEBENCHSUITE4 - Benchmark. TCP local IP address. Single client and single server in parallel.
  // ngBefore++
  // ExampleBenchSuite4()


  ngAfter := runtime.NumGoroutine()
  // prim(0, fmt.Sprintf("[MAIN] runtime.NumGoroutine(): %v", ngAfter))
  if ngAfter - ngBefore != 0 {
    panic("The number of running go rouitine is leaking.")
  } else {
    // prim(0, "[MAIN] With the number of running go rouitine, everything is fine.")
  }
}

func Example1() {
  // Use "EXAMPLE1" and you can use the linux command "nc" to connect multiple times.
  // Shutdown processing of multiple connection is performed correctly.
  // BASH CMD: nc -U /run/user/1000/sockbench/socket
  // BASH CMD: pid=`pgrep sockbench`; path="/run/user/1000/sockbench/1/socket_${pid}_1"; nc -U $path

  // HINT from man nc: "After the file has been transferred, the connection will close automatically."

  // LINUX KERNEL: Minimum, Default, Maximum memory size values in bytes.
  // > cat /proc/sys/net/ipv4/tcp_rmem
  // 4096    131072  6291456
  // cat /proc/sys/net/ipv4/tcp_wmem
  // 4096    16384   4194304

  // EXAMPLE1 - Multiple connections on single server:
  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  printInfo = true
  var id int = 1
  ex1UnixServerMultiCon(id, &sl)
}

func Example2Unix() {
  // EXAMPLE2 - Single client and single server in parallel goroutine mode:
  // Single client and single server in parallel goroutine mode. The client does
  // not close the connection. The server closes the connection after receiving
  // data and receiving a signal from the client.

  // Try looping the program. It is understood that there will be no leakage of the
  // processor, memory, socket files, etc. To improve performance, it is recommended
  // to compile the program.
  // MAN LSOF: lsof - list open files
  // BASH CMD: while [ : ]; do go run sockbench.go; if [ $? -ne 0 ]; then break; fi done
  // BASH CMD: while [ : ]; do ./sockbench; if [ $? -ne 0 ]; then break; fi done
  // BASH CMD: while [ : ]; do find /run/user/1000/sockbench/; sleep 0.01; done
  // BASH CMD: while [ : ]; do echo 'D'; ss -x -a | grep sockbench; sleep 0.01; done
  // BASH CMD: while [ : ]; do lsof -p `pgrep sockbench` | wc -l; sleep 0.01; done
  // BASH CMD: while [ : ]; do lsof -p `pgrep sockbench` | wc -l >> sockbench_lsof_test.txt; sleep 0.01; done
  //           cat sockbench_lsof_test.txt | sort -g
  //           100 goroutine ~~~~ 300 file discriptors

  // UNIX:
  printInfo = true
  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  var wg sync.WaitGroup
  var bs benchState
  var it int = 1
  var buf int = 1
  var id int = 1
  var b byte = 255
  var s []byte
  var r benchResult
  for i := 0; i < buf; i++ {
    s = append(s, b)
  }
  s = append(s, '|')

  wg.Add(1)
  go ex2UnixServer(&bs, &wg, it, buf, id, &sl)
  wg.Add(1)
  go ex2UnixClient(&bs, &wg, it, buf, id, s)
  wg.Wait()
  r.calcResult(&bs, it, buf, "FS=tmpfs(RAM)", "Unix socket")
  r.printResult()
}

func Example2TCP(ip string) {
  // BASH CMD: nc 127.0.0.1 6080
  // TCP:
  printInfo = true
  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  var wg sync.WaitGroup
  var bs benchState
  var it int = 1
  var buf int = 1
  var id int = 1
  var b byte = 255
  var s []byte;
  var r benchResult
  for i := 0; i < buf; i++ {
    s = append(s, b)
  }
  s = append(s, '|')

  wg.Add(1)
  go ex2TCPServer(&bs, &wg, it, buf, id, &sl, ip)
  wg.Add(1)
  go ex2TCPClient(&bs, &wg, it, buf, id, s, ip)
  wg.Wait()
  r.calcResult(&bs, it, buf, "FS=tmpfs(RAM)", "Unix socket")
  r.printResult()
}

func Example3UnixParallel() {
  // EXAMPLE3 - Two client and two server in parallel mode with parallel start:
  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  var wg sync.WaitGroup
  var b byte = 255
  var s []byte;
  for i := 0; i < 1; i++ {
    s = append(s, b)
  }
  s = append(s, '|')

  prim(0, "[MAIN] BENCH 1 1 1")
  var bs  benchState
  var it  int = 1
  var buf int = 1
  var id  int = 1
  var r   benchResult

  wg.Add(1)
  go ex2UnixServer(&bs, &wg, it, buf, id, &sl)
  wg.Add(1)
  go ex2UnixClient(&bs, &wg, it, buf, id, s)


  prim(0, "[MAIN] BENCH 1 1 2")
  var bs2  benchState
  var it2  int = 1
  var buf2 int = 1
  var id2  int = 2
  var r2   benchResult

  wg.Add(1)
  go ex2UnixServer(&bs2, &wg, it2, buf2, id2, &sl)
  wg.Add(1)
  go ex2UnixClient(&bs2, &wg, it2, buf2, id2, s)

  wg.Wait()
  r.calcResult(&bs, it, buf, "FS=tmpfs(RAM)", "Unix socket")
  r2.calcResult(&bs2, it2, buf2, "FS=tmpfs(RAM)", "Unix socket")
  r.printResult()
  r2.printResult()
}

func ExampleCrazy1() {
  // EXAMPLECRAZY1 - Cycle launch of a single client/server in parallel:

  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  var wg sync.WaitGroup
  var b byte = 255
  var s []byte
  for i := 0; i < 1; i++ {
    s = append(s, b)
  }
  s = append(s, '|')

  prim(0, "[MAIN] BENCH 1 1 1")
  var bs  benchState
  var it  int = 1
  var buf int = 1
  var id  int = 1
  var r   benchResult

  wg.Add(1)
  go ex2UnixServer(&bs, &wg, it, buf, id, &sl)

  wg.Add(1)
  go ex2UnixClient(&bs, &wg, it, buf, id, s)

  prim(0, "[MAIN] BENCH 1 1 2")
  var bs2  benchState
  var it2  int = 1
  var buf2 int = 1
  var id2  int = 2
  var r2   benchResult

  wg.Add(1)
  go ex2UnixServer(&bs2, &wg, it2, buf2, id2, &sl)

  wg.Add(1)
  go ex2UnixClient(&bs2, &wg, it2, buf2, id2, s)

  prim(0, fmt.Sprintf("[MAIN] runtime.NumGoroutine(): %v", runtime.NumGoroutine()))
  wg.Wait()

  r.calcResult(&bs, it, buf, "FS=tmpfs(RAM)", "Unix socket")
  r.printResult()
  r2.calcResult(&bs2, it2, buf2, "FS=tmpfs(RAM)", "Unix socket")
  r2.printResult()
}

func ExampleCrazy2(batching int) {
  // EXAMPLECRAZY2 - Cycle launch of a huge variety client/server in parallel:
  // Next using crazy mode. Be careful.

  // You can safely use the keyboard shortcut Ctrl + C in the Linux console. The
  // server is supposed to clean up after itself.
  // The number of goroutine can increase infinitely if the server does not have
  // time to process the data.
  // Be careful: panic: accept unix /run/user/1000/sockbench/1685/socket_12810_1685: accept4: too many open files

  // Single client and single server in parallel. Parallel start. May be infinity loop.
  // Batch start and stop.
  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  var b byte = 255
  var s []byte
  for i := 0; i < 1; i++ {
    s = append(s, b)
  }
  s = append(s, '|')

  printInfo = false
  for x := 1; x <= batching; x++ {
    // for x := 0; x < 1; x++ {
    prim(0, fmt.Sprintf("[MAIN] x=%v", x))
    prim(0, fmt.Sprintf("[MAIN] runtime.NumGoroutine(): %v", runtime.NumGoroutine()))

    var wg sync.WaitGroup
    // 200 goroutines run in parallel.
    // ulimit -n 1024   goroutines=200   Y=100 -> OK
    // ulimit -n 2048   goroutines=400   Y=200 -> OK
    // Need change the ulimit to avoid the error "too many open files".

    for y := 1; y <= 200; y++ {
      // 1 iteration -> 2 goroutines
      prim(0, fmt.Sprintf("[MAIN] x=%v, y=%v", x, y))

      var bs benchState
      var it int = 1
      var buf int = 1

      wg.Add(1)
      go ex2UnixServer(&bs, &wg, it, buf, y, &sl)
      wg.Add(1)
      go ex2UnixClient(&bs, &wg, it, buf, y, s)
    }

    // Waiting for the completion of the package goroutines.
    prim(0, fmt.Sprintf("[MAIN] runtime.NumGoroutine(): %v", runtime.NumGoroutine()))
    wg.Wait()
  }
}

func ExampleBenchUnix1() {
  // BENCHUNIX1 - Benchmark. Unix socket. Single client and single server in parallel.
  // Iterations = 100   Buffer size = 1 byte
  // 1 (1 b)   4096 (4 kB)    16384 (16 kB)   4194304 (4 mB)   6291456 (6 mB)
  // 1 (1 b)        * 100 = bytes
  // 4096 (4 kB)    * 100 =  ~400 kB
  // 16384 (16 kB)  * 100 = ~1600 kB
  // 4194304 (4 mB) * 100 =  ~400 mB
  // 6291456 (6 mB) * 100 =  ~600 mB
  printInfo = false
  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  var it int
  var buf int
  var r benchResult

  it  = 100
  buf = 1
  r   = sockBenchUnix(it, buf, &sl)
  r.printResult()

  it  = 100
  buf = 4096
  r   = sockBenchUnix(it, buf, &sl)
  r.printResult()

  it  = 100
  buf = 16384
  r   = sockBenchUnix(it, buf, &sl)
  r.printResult()

  it  = 100
  buf = 4194304
  r   = sockBenchUnix(it, buf, &sl)
  r.printResult()

  it  = 100
  buf = 6291456
  r   = sockBenchUnix(it, buf, &sl)
  r.printResult()

  it  = 100
  buf = 62914560
  r   = sockBenchUnix(it, buf, &sl)
  r.printResult()
}

func ExampleBenchSuite1() {
  // BENCHUNIX1 - Benchmark. Unix socket. Single client and single server in parallel.
  // Iterations = 100   Buffer size = 1 byte
  // 1 (1 b)   4096 (4 kB)    16384 (16 kB)   4194304 (4 mB)   6291456 (6 mB)
  // 1 (1 b)        * 100 = bytes
  // 4096 (4 kB)    * 100 =  ~400 kB
  // 16384 (16 kB)  * 100 = ~1600 kB
  // 4194304 (4 mB) * 100 =  ~400 mB
  // 6291456 (6 mB) * 100 =  ~600 mB
  printInfo = false

  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  var it  int
  var r   benchResult
  var rs  []benchResult

  it  = 100
  buf := []int{1, 2, 4, 16, 256, 4096, 16384, 4194304, 6291456, 62914560}

  for _, b := range buf {
    print(".")
    r   = sockBenchUnix(it, b, &sl)
    rs  = append(rs, r)

    print(".")
    r   = sockBenchTCP(it, b, &sl, "127.0.0.1") // is Loopback; PORT=6080
    rs  = append(rs, r)

    print(".")
    r   = sockBenchTCP(it, b, &sl, ifaceIpAddr) // is not Loopback; Local inet address; PORT=6080
    rs  = append(rs, r)
  }
  print("\n")

  benchPrintTable(rs)
}

func ExampleBenchSuite2() {
  printInfo = false

  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  var it  int
  var r   benchResult
  var rs  []benchResult

  it  = 100
  buf := []int{1, 2, 4, 16, 256, 4096, 16384, 4194304, 6291456, 62914560}

  for _, b := range buf {
    print(".")
    r   = sockBenchUnix(it, b, &sl)
    rs  = append(rs, r)
  }
  print("\n")

  benchPrintTable(rs)
}

func ExampleBenchSuite3() {
  printInfo = false

  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  var it  int
  var r   benchResult
  var rs  []benchResult

  it  = 100
  buf := []int{1, 2, 4, 16, 256, 4096, 16384, 4194304, 6291456, 62914560}

  for _, b := range buf {
    print(".")
    r   = sockBenchTCP(it, b, &sl, "127.0.0.1") // is Loopback; PORT=6080
    rs  = append(rs, r)
  }
  print("\n")

  benchPrintTable(rs)
}

func ExampleBenchSuite4() {
  printInfo = false

  var sl serversList
  sl.build()
  sl.regSyscallSIGTERM()

  var it  int
  var r   benchResult
  var rs  []benchResult

  it  = 100
  buf := []int{1, 2, 4, 16, 256, 4096, 16384, 4194304, 6291456, 62914560}

  for _, b := range buf {
    print(".")
    r   = sockBenchTCP(it, b, &sl, ifaceIpAddr) // is not Loopback; Local inet address; PORT=6080
    rs  = append(rs, r)
  }
  print("\n")

  benchPrintTable(rs)
}

func level(id int, m string, l string) {
  t := time.Now()
  ns := strconv.Itoa(t.Nanosecond())
  ns = ns + strings.Repeat("0", 10-len(ns))
  tf := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d.%v", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second(), ns)
  fmt.Printf("[%v][%v][%v]%v\n", tf, id, l, m)
}

func info(id int, m string) {
  // You must turn off for more accurate performance measurement. Just comment out.
  if printInfo {
    level(id, m, "INFO")
  }
}

func prim(id int, m string) {
  if printPrim {
    level(id, m, "PRIM")
  }
}

type benchResult struct {
  begin               int64
  end                 int64
  nanosec             int64
  microsec            float64
  millisec            float64
  sec                 float64
  it                  int
  bufSize             int
  bytes               int
  kbytes              float64
  mbytes              float64
  gbytes              float64
  comment             string
  t                   string
  speedB              float64
  speedKB             float64
  speedMB             float64
  speedGB             float64
  speedKb             float64
  speedMb             float64
  speedGb             float64
  speedPrettyBytes    string
  speedPrettyBits     string
  timePretty          string
}

func (br *benchResult) calcResult(bs *benchState, it int, bufSize int, com string, t string) {
  // One iteration contains bufSize (slice of bytes) plus 1 rune delimeter (1 byte)
  br.begin    = bs.begin
  br.end      = bs.end
  br.nanosec  = (br.end - br.begin)
  br.microsec = float64(br.nanosec) / 1000
  br.millisec = br.microsec / 1000
  br.sec      = br.millisec / 1000
  br.it       = it
  br.bufSize  = bufSize
  br.bytes    = it * (bufSize + 1)
  br.kbytes   = float64(br.bytes) / 1024
  br.mbytes   = br.kbytes / 1024
  br.gbytes   = br.mbytes / 1024
  br.comment  = com
  br.t        = t

  // bytes ---- N sec
  // speed ---- 1 sec
  br.speedB   = float64(br.bytes)  / br.sec
  br.speedKB  = br.kbytes / br.sec
  br.speedMB  = br.mbytes / br.sec
  br.speedGB  = br.gbytes / br.sec

  if br.speedGB > 1 {
    br.speedPrettyBytes = fmt.Sprintf("%.2f GByte/s", br.speedGB)
    br.speedPrettyBits  = fmt.Sprintf("%.2f Gbit/s", br.speedGB * 8)
  } else if br.speedMB > 1 {
    br.speedPrettyBytes = fmt.Sprintf("%.2f MByte/s", br.speedMB)
    br.speedPrettyBits  = fmt.Sprintf("%.2f Mbit/s", br.speedMB * 8)
  } else if br.speedKB > 1 {
    br.speedPrettyBytes = fmt.Sprintf("%.2f KByte/s", br.speedKB)
    br.speedPrettyBits  = fmt.Sprintf("%.2f Kbit/s", br.speedKB * 8)
  } else {
    br.speedPrettyBytes = fmt.Sprintf("%.2f Byte/s", br.speedB)
    br.speedPrettyBytes = fmt.Sprintf("%.2f bit/s", br.speedB * 8)
  }

  br.timePretty = fmt.Sprintf("%.2f sec", br.sec)
}

func (br *benchResult) printResult() {
  prim(0, fmt.Sprintf(
    "[MAIN][BENCH] Comment: %v; Type: %v; Speed=%v; Bytes=%v; Iterations=%v; Buffer size=%v; Time: %vms %vmicrosec %vns",
    br.comment,
    br.t,
    br.speedPrettyBytes,
    br.bytes,
    br.it,
    br.bufSize,
    br.millisec,
    br.microsec,
    br.nanosec))
}

func benchPrintTable(s []benchResult) {
  l := strings.Repeat("-", 141)

  fmt.Println(l)
  fmt.Printf(
    "| %3v | %-50v | %16v | %16v | %4v | %20v | %10v |\n",
    "#", "Comment+Type:", "Speed (bits):", "Speed (bytes):", "It:", "Buffer size:", "Time:")
  fmt.Println(l)

  for i, v := range s {
    fmt.Printf(
      "| %3v | %-50v | %16v | %16v | %4v | %20v | %10v |\n",
      i+1, v.comment+", "+v.t, v.speedPrettyBits, v.speedPrettyBytes, v.it, v.bufSize, v.timePretty)
  }
  fmt.Println(l)
}

func sockBenchUnix(it int, bufSize int, sl *serversList) benchResult {
  var wg sync.WaitGroup
  var bs benchState
  var id int = 1
  var s []byte
  var b byte = 255
  var r benchResult

  // prim(id, fmt.Sprintf("[MAIN][BENCH] Start with it=%v; bufSize=%v", it, bufSize))
  // prim(id, fmt.Sprintf("[MAIN] base2=%b, base10=%d\n", b, b))
  for i := 0; i < bufSize; i++ {
    s = append(s, b)
  }
  s = append(s, '|')
  // prim(id, fmt.Sprintf("[MAIN] len(s)=%v; cap(s)=%v", len(s), cap(s)))

  wg.Add(1)
  go ex2UnixServer(&bs, &wg, it, bufSize, id, sl)
  wg.Add(1)
  go ex2UnixClient(&bs, &wg, it, bufSize, id, s)
  wg.Wait()

  r.calcResult(&bs, it, bufSize, "FS=tmpfs(RAM)", "Unix socket")
  return r
}

func sockBenchTCP(it int, bufSize int, sl *serversList,ip string) benchResult {
  var wg sync.WaitGroup
  var bs benchState
  var id int = 1
  var s []byte
  var b byte = 255
  var r benchResult

  // prim(id, fmt.Sprintf("[MAIN][BENCH] Start with it=%v; bufSize=%v", it, bufSize))
  // prim(id, fmt.Sprintf("[MAIN] base2=%b, base10=%d\n", b, b))
  for i := 0; i < bufSize; i++ {
    s = append(s, b)
  }
  s = append(s, '|')
  // prim(id, fmt.Sprintf("[MAIN] len(s)=%v; cap(s)=%v", len(s), cap(s)))

  wg.Add(1)
  go ex2TCPServer(&bs, &wg, it, bufSize, id, sl, ip)
  wg.Add(1)
  go ex2TCPClient(&bs, &wg, it, bufSize, id, s, ip)
  wg.Wait()

  r.calcResult(&bs, it, bufSize, "FS=tmpfs(RAM)", "TCP4 socket "+ip)
  return r
}

type serversList struct {
  servers map[string]serverConf
  mutex   sync.Mutex
}

func (sl *serversList) build() {
  sl.servers = make(map[string]serverConf)
}

func (sl *serversList) addServer(sc serverConf) {
  sl.mutex.Lock()
  sl.servers[sc.key] = sc
  sl.mutex.Unlock()
}

func (sl *serversList) removeServer(sc serverConf) {
  sl.mutex.Lock()
  delete(sl.servers, sc.key)
  sl.mutex.Unlock()
}

func (sl *serversList) regSyscallSIGTERM() {
  signalChannel := make(chan os.Signal, 1)
  signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
  go func(c chan os.Signal, sl *serversList) {
    signal := <-c
    prim(0, fmt.Sprintf("[MAIN][INTERCEPTOR] Signal is intercepted by channel: %v", signal))
    prim(0, fmt.Sprintf("[MAIN][INTERCEPTOR] len(sl.servers): %v", len(sl.servers)))
    // GODOC: Close stops listening on the Unix address. Already accepted connections are not closed.
    for k, s := range sl.servers {
      prim(0, fmt.Sprintf("[MAIN][INTERCEPTOR]: Clean up server id = %v with type = %v", k, s.addrType))
      s.listenerClose()
      if s.addrType == "unix" {
        s.removeSockFile()
        s.removeSockDir()
      }
    }

    os.Exit(1)
  }(signalChannel, sl)
}

type benchState struct {
  state string
  mutex sync.Mutex
  begin int64
  end   int64
}

func (bs *benchState) setRunningServer() {
  bs.mutex.Lock()
  bs.state = "running_server"
  bs.mutex.Unlock()
}

func (bs *benchState) getRunningServer() bool {
  bs.mutex.Lock()
  defer bs.mutex.Unlock()
  return bs.state == "running_server"
}

func (bs *benchState) setClientStopped() {
  bs.mutex.Lock()
  bs.state = "client_stopped"
  bs.mutex.Unlock()
}

func (bs *benchState) getClientStopped() bool {
  bs.mutex.Lock()
  defer bs.mutex.Unlock()
  return bs.state == "client_stopped"
}

func (bs *benchState) setServerStopped() {
  bs.mutex.Lock()
  bs.state = "server_stopped"
  bs.mutex.Unlock()
}

func (bs *benchState) getServerStopped() bool {
  bs.mutex.Lock()
  defer bs.mutex.Unlock()
  return bs.state == "server_stopped"
}

type serverConf struct {
  id           int
  key          string
  userCur      *user.User
  addrType     string

  sockDir      string
  sockFile     string
  postfix      string
  unixAddr     *net.UnixAddr
  unixListener *net.UnixListener

  ipAddr       net.IP
  ipPort       int
  tcpAddr      *net.TCPAddr
  tcpListener  *net.TCPListener
}

func (sc *serverConf) build(id int) {
  sc.id = id
  sc.key = strconv.Itoa(id)
  sc.userCur, _ = user.Current()
}

func (sc *serverConf) buildUnix(id int) {
  sc.build(id)
  sc.addrType = "unix"
  sc.sockDir  = path.Join("/", "run", "user", sc.userCur.Uid, "sockbench", sc.key)
  sc.postfix  = "socket_" + strconv.Itoa(os.Getpid()) + "_" + sc.key
  sc.sockFile = path.Join(sc.sockDir, sc.postfix)
  sc.unixAddr = &net.UnixAddr{sc.sockFile, sc.addrType}
  info(id, "[SERVER] Uid: "+sc.userCur.Uid+"; Dir: "+sc.sockDir+"; Socket: "+sc.sockFile+"; Addr: "+sc.unixAddr.String())
}

func (sc *serverConf) buildTcp(id int, ip string) {
  sc.build(id)
  sc.addrType = "tcp4"
  sc.ipAddr   = net.ParseIP(ip)
  sc.ipPort   = 6080
  sc.tcpAddr  = &net.TCPAddr{sc.ipAddr, sc.ipPort, ""}
  info(id, "[SERVER] IpAddr: "+sc.ipAddr.String()+"; IpPort: "+strconv.Itoa(sc.ipPort)+"; Addr: "+sc.tcpAddr.String())
}

func (sc *serverConf) createSockDir() {
  if _, err1 := os.Stat(sc.sockDir); os.IsNotExist(err1) {
    err2 := os.MkdirAll(sc.sockDir, 0700)
    if err2 != nil {
      panic(err2)
    } else {
      info(sc.id, fmt.Sprintf("[SERVER] Dir was created: %v", sc.sockDir))
    }
  } else {
    info(sc.id, fmt.Sprintf("[SERVER] Directory already exists: %v", sc.sockDir))
  }
}

func (sc *serverConf) checkSockDir() {
  var oldSockFiles []string
  err := filepath.Walk(sc.sockDir, func(path string, fileinfo os.FileInfo, err error) error {
    if err != nil {
      panic(fmt.Sprintf("[SERVER] Prevent panic by handling failure accessing a path: %v; err: %v\n", path, err))
    }
    if path != sc.sockDir {
      oldSockFiles = append(oldSockFiles, path)
    }
    info(sc.id, fmt.Sprintf("[SERVER] Visited files or directory: %v\n", path))
    return nil
  })
  if err != nil {
    panic(err)
  }
  if len(oldSockFiles) != 0 {
    panic("[SERVER] A unix file socket leak occurs. Be careful. To continue, you need to delete the extra socket files manually. Example: rm -rfv /run/user/1000/sockbench")
  }
}

func (sc *serverConf) removeSockFile() {
  os.Remove(sc.sockFile)
  if _, err := os.Stat(sc.sockFile); os.IsExist(err) {
    panic(err)
  }
  info(sc.id, fmt.Sprintf("[SERVER] Socket file was removed: %v", sc.sockFile))
}

func (sc *serverConf) removeSockDir() {
  os.Remove(sc.sockDir)
  if _, err := os.Stat(sc.sockDir); os.IsExist(err) {
    panic(err)
  }
  info(sc.id, fmt.Sprintf("[SERVER] Socket dir was removed: %v", sc.sockDir))
}

func (sc *serverConf) addUnixListener(l *net.UnixListener) {
  sc.unixListener = l
}

func (sc *serverConf) addTCPListener(l *net.TCPListener) {
  sc.tcpListener = l
}

func (sc *serverConf) listenerClose() {
  switch sc.addrType {
  case "unix":
    sc.unixListener.Close()
  case "tcp4":
    sc.tcpListener.Close()
  }
}

type clientConf struct {
  id           int
  key          string
  userCur      *user.User
  addrType     string

  sockDir      string
  sockFile     string
  postfix      string
  unixAddr     *net.UnixAddr
  unixListener *net.UnixListener

  ipAddr       net.IP
  ipPort       int
  tcpAddr      *net.TCPAddr
  tcpListener  *net.TCPListener
}

func (cc *clientConf) build(id int) {
  cc.id = id
  cc.key = strconv.Itoa(id)
  cc.userCur, _ = user.Current()
}

func (cc *clientConf) buildUnix(id int) {
  cc.build(id)
  cc.addrType = "unix"
  cc.sockDir  = path.Join("/", "run", "user", cc.userCur.Uid, "sockbench", cc.key)
  cc.postfix  = "socket_" + strconv.Itoa(os.Getpid()) + "_" + cc.key
  cc.sockFile = path.Join(cc.sockDir, cc.postfix)
  cc.unixAddr = &net.UnixAddr{cc.sockFile, cc.addrType}
  info(id, "[CLIENT] Uid: "+cc.userCur.Uid+"; Dir: "+cc.sockDir+"; Socket: "+cc.sockFile+"; Addr: "+cc.unixAddr.String())
}

func (cc *clientConf) buildTcp(id int, ip string) {
  cc.build(id)
  cc.addrType = "tcp4"
  cc.ipAddr   = net.ParseIP(ip)
  cc.ipPort   = 6080
  cc.tcpAddr  = &net.TCPAddr{cc.ipAddr, cc.ipPort, ""}
  info(id, "[CLIENT] IpAddr: "+cc.ipAddr.String()+"; IpPort: "+strconv.Itoa(cc.ipPort)+"; Addr: "+cc.tcpAddr.String())
}

func ex2UnixServer(bs *benchState, wg *sync.WaitGroup, maxIt int, bufSize int, id int, sl *serversList) {
  var sc serverConf
  sc.buildUnix(id)
  sl.addServer(sc)

  defer func(sc serverConf, sl *serversList) {
    sc.removeSockFile()
    sc.removeSockDir()
    bs.setServerStopped()
    wg.Done()
    sl.removeServer(sc)
    info(id, "[SERVER] Call wg.Done()")
  }(sc, sl)
  // Mount: tmpfs on /run/user/1000 type tmpfs (rw,nosuid,nodev,relatime,size=1610456k,mode=700,uid=1000,gid=100)

  // tmpfs - a virtual memory filesystem The tmpfs facility allows the
  // creation of filesystems whose contents reside in virtual memory. Since
  // the files on such filesystems typically re- side in RAM, file access is
  // extremely fast.

  // net.ListenUnix(addrType): There are three types of unix domain socket.
  // “unix”       corresponds to SOCK_STREAM
  // “unixdomain” corresponds to SOCK_DGRAM
  // “unixpacket” corresponds to SOCK_SEQPACKET
  // man 2 socket

  // > ps -T -o pid,ppid,%cpu,cputime,cputimes,%mem,rss,sz,vsz --pid 6819
  // > ps -L -o pid,ppid,%cpu,cputime,cputimes,%mem,rss,sz,vsz --pid 6819

  // var err error = nil

  sc.createSockDir()
  defer func() {
    sc.removeSockDir()
    info(id, "[SERVER] Call sc.removeSockDir()")
  }()
  sc.checkSockDir()

  listener, err := net.ListenUnix(sc.addrType, sc.unixAddr)
  if err != nil {
    panic(err)
  }
  defer func() {
    listener.Close()
    info(id, fmt.Sprintf("[SERVER] Listener was closed: %v", listener))
  }()
  info(id, fmt.Sprintf("[SERVER] Listener was created: %v", listener))
  sc.addUnixListener(listener)

  listener.SetUnlinkOnClose(true)
  info(id, "[SERVER] SetUnlinkOnClose(true) executed.")
  info(id, "[SERVER][STANDBY]")

  bs.setRunningServer()
  info(id, "[SERVER] bs.setRunningServer() executed.")

  info(id, "[SERVER] The action AcceptUnix() will be executed and go into standby mode.")
  conn, err := listener.AcceptUnix()
  if err != nil {
    panic(err)
  }
  info(id, fmt.Sprintf("[SERVER] The action AcceptUnix() for unix socket was executed: %v", conn))
  defer func() {
    conn.Close()
    info(id, fmt.Sprintf("[SERVER] Connection was closed: %v", conn))
  }()

  // DOC: SetReadBuffer sets the size of the operating system's receive buffer associated with the connection.
  // func (*UnixConn) SetReadBuffer
  // SRC: src/net/sockopt_posix.go
  // func setReadBuffer(fd *netFD, bytes int) error {
  //   err := fd.pfd.SetsockoptInt(syscall.SOL_SOCKET, syscall.SO_RCVBUF, bytes)
  //   runtime.KeepAlive(fd)
  //   return wrapSyscallError("setsockopt", err)
  // }
  err = conn.SetReadBuffer(bufSize)
  if err != nil {
    panic(err)
  }
  info(id, fmt.Sprintf("[SERVER] SetReadBuffer(%v) executed.", bufSize))

  // CloseWrite shuts down the writing side of the Unix domain connection. Most callers should just use Close.
  err = conn.CloseWrite()
  if err != nil {
    panic(err)
  }
  info(id, "[SERVER] CloseWrite() executed.")

  reader := bufio.NewReaderSize(conn, bufSize)
  var i int

  start := time.Now().UnixNano()
  for i = 0; i < maxIt; i++ {
    // info(id, fmt.Sprintf("[SERVER] Iteration=%v", i))
    // info(id, fmt.Sprintf("[SERVER] Iteration=%v; Before Read()...", i))
    // var s []byte
    // s, err := reader.ReadBytes('|')
    _, err := reader.ReadBytes('|')
    // info(id, fmt.Sprintf("[SERVER] Iteration=%v; After Read()...", i))
    if err != nil {
      panic(err)
      // info(id, fmt.Sprintf("[SERVER][DISCONNECT] %v", err))
      // break
    }

    // info(id, fmt.Sprintf("[SERVER] Iteration=%v len(s)=%v; cap(s)=%v; s: %v", i, len(s), cap(s), s))
    // for _, v := range s {
    //   info(id, fmt.Sprintf("[SERVER] Iteration=%v; Read byte2=%b byte10=%d, symbol=%c, symbol_quote=%q", i, v, v, v, v));
    // }
  }
  bs.begin = start
  bs.end = time.Now().UnixNano()

  for {
    if bs.getClientStopped() {
      info(id, "[SERVER] The server received a signal from the client and shut down")

      err = conn.Close()
      info(id, fmt.Sprintf("[SERVER] Connection was closed by signal from client: %v", conn))
      if err != nil {
        panic(err)
      }

      err = listener.Close()
      info(id, fmt.Sprintf("[SERVER] Listener was closed by signal from client: %v", listener))
      if err != nil {
        panic(err)
      }

      sc.removeSockFile()
      sc.removeSockDir()
      bs.setServerStopped()
      sl.removeServer(sc)
      info(id, "[SERVER][FINISH]")

      break
    }
    time.Sleep(1 * time.Millisecond)
  }
}

func ex2UnixClient(bs *benchState, wg *sync.WaitGroup, maxIt int, bufSize int, id int, s []byte) {
  defer func() {
    for {
      time.Sleep(1 * time.Millisecond)
      if bs.getServerStopped() {
        wg.Done()
        info(id, "[CLIENT] Call wg.Done().")
        break
      }
    }
  }()

  var err error = nil
  var conn *net.UnixConn
  var connErr error = nil
  var runningServer bool = false
  var cc clientConf
  cc.buildUnix(id)

  for {
    info(id, "[CLIENT] Try again later...")
    time.Sleep(1 * time.Millisecond)

    if !runningServer && bs.getRunningServer() {
      info(id, "[CLIENT] The client received a signal from the server: state=running_server")
      runningServer = true
    } else if runningServer {
      info(id, "[CLIENT] The client received a signal from the server: state=running_server")
    } else {
      info(id, "[CLIENT] The client not received a signal from the server: state=running_server")
      continue
    }

    conn, connErr = net.DialUnix(cc.addrType, nil, cc.unixAddr)
    if connErr != nil {
      info(id, fmt.Sprintf("[CLIENT] Error connection: %v", connErr))
      connErr = nil
      continue
    }

    if runningServer && connErr == nil {
      info(id, "[CLIENT] Get server state=running_server and connection was created.")
      break
    }
  }

  // << 10 // 1024        kb
  // << 20 // 1048576     kb=1mb
  // << 25 // 33554432    kb
  // << 30 // 1073741824  kb
  // GODOC: SetWriteBuffer sets the size of the operating system's transmit buffer associated with the connection.
  // func (*UnixConn) SetWriteBuffer
  // SRC: src/net/sockopt_posix.go
  // func setWriteBuffer(fd *netFD, bytes int) error {
  //   err := fd.pfd.SetsockoptInt(syscall.SOL_SOCKET, syscall.SO_SNDBUF, bytes)
  //   runtime.KeepAlive(fd)
  //   return wrapSyscallError("setsockopt", err)
  // }
  err = conn.SetWriteBuffer(bufSize)
  if err != nil {
    panic(err)
  }
  info(id, fmt.Sprintf("[CLIENT] SetWriteBuffer(%v) executed", bufSize))

  // CloseRead shuts down the reading side of the Unix domain connection. Most callers should just use Close.
  err = conn.CloseRead()
  if err != nil {
    panic(err)
  }
  info(id, "[CLIENT] CloseRead() executed")

  for i := 0; i < maxIt; i++ {
    // info(id, fmt.Sprintf("[CLIENT] Iteration=%v", i))
    // info(id, "[CLIENT] The action Write() will be executed...")

    // var n int
    // info(id, "[CLIENT] Start writing...")

    // n, err = conn.Write(s)
    _, err = conn.Write(s)
    if err != nil {
      panic(err)
    }
    // info(id, fmt.Sprintf("[CLIENT] Write %v bytes. Err=%v", n, err))
  }

  // It is expected that the connection will be closed by the server:
  bs.setClientStopped()
  info(id, "[CLIENT] The client has finished and sent a signal to the server")
  info(id, "[CLIENT][FINISH]")
}

func ex2TCPServer(bs *benchState, wg *sync.WaitGroup, maxIt int, bufSize int, id int, sl *serversList, ip string) {
  // GODOC: func (*TCPConn) SetLinger
  //   SetLinger sets the behavior of Close on a connection which still has data
  //   waiting to be sent or to be acknowledged.
  // GODOC: func (*TCPConn) SetNoDelay
  //   SetNoDelay controls whether the operating system should delay packet
  //   transmission in hopes of sending fewer packets (Nagle's algorithm). The default
  //   is true (no delay), meaning that data is sent as soon as possible after a
  //   Write.
  // GODOC: func (*TCPListener) File
  //   File returns a copy of the underlying os.File. It is the caller's responsibility
  //   to close f when finished. Closing l does not affect f, and closing f does not
  //   affect l.The returned os.File's file descriptor is different from the
  //   connection's. Attempting to change properties of the original using this
  //   duplicate may or may not have the desired effect.

  var sc serverConf
  sc.buildTcp(id, ip)
  sl.addServer(sc)
  info(id, fmt.Sprintf("[SERVER] sc.ipAddr.IsLoopback()=%v", sc.ipAddr.IsLoopback()))

  defer func(sc serverConf, sl *serversList) {
    bs.setServerStopped()
    wg.Done()
    sl.removeServer(sc)
    info(id, "[SERVER] Call wg.Done()")
  }(sc, sl)

  listener, err := net.ListenTCP(sc.addrType, sc.tcpAddr)
  if err != nil {
    panic(err)
  }
  defer func() {
    listener.Close()
    info(id, fmt.Sprintf("[SERVER] Listener was closed: %v", listener))
  }()
  info(id, fmt.Sprintf("[SERVER] Listener was created: %v", listener))
  sc.addTCPListener(listener)

  info(id, "[SERVER][STANDBY]")
  bs.setRunningServer()
  info(id, "[SERVER] bs.setRunningServer() executed.")
  info(id, "[SERVER] The action AcceptTCP() will be executed and go into standby mode.")
  conn, err := listener.AcceptTCP()
  if err != nil {
    panic(err)
  }
  info(id, fmt.Sprintf("[SERVER] The action AcceptTCP() for tcp socket was executed: %v", conn))
  defer func() {
    conn.Close()
    info(id, fmt.Sprintf("[SERVER] Connection was closed: %v", conn))
  }()

  // DOC: SetReadBuffer sets the size of the operating system's receive buffer associated with the connection.
  // func (*UnixConn) SetReadBuffer
  // SRC: src/net/sockopt_posix.go
  // func setReadBuffer(fd *netFD, bytes int) error {
  //   err := fd.pfd.SetsockoptInt(syscall.SOL_SOCKET, syscall.SO_RCVBUF, bytes)
  //   runtime.KeepAlive(fd)
  //   return wrapSyscallError("setsockopt", err)
  // }
  err = conn.SetReadBuffer(bufSize)
  if err != nil {
    panic(err)
  }
  info(id, fmt.Sprintf("[SERVER] SetReadBuffer(%v) executed.", bufSize))

  // CloseWrite shuts down the writing side of the Unix domain connection. Most callers should just use Close.
  err = conn.CloseWrite()
  if err != nil {
    panic(err)
  }
  info(id, "[SERVER] CloseWrite() executed.")

  reader := bufio.NewReaderSize(conn, bufSize)
  var i int

  start := time.Now().UnixNano()
  for i = 0; i < maxIt; i++ {
    // info(id, fmt.Sprintf("[SERVER] Iteration=%v", i))
    // info(id, fmt.Sprintf("[SERVER] Iteration=%v; Before Read()...", i))
    // var s []byte
    // s, err := reader.ReadBytes('|')
    _, err := reader.ReadBytes('|')
    // info(id, fmt.Sprintf("[SERVER] Iteration=%v; After Read()...", i))
    if err != nil {
      panic(err)
      // info(id, fmt.Sprintf("[SERVER][DISCONNECT] %v", err))
      // break
    }

    // info(id, fmt.Sprintf("[SERVER] Iteration=%v len(s)=%v; cap(s)=%v; s: %v", i, len(s), cap(s), s))
    // for _, v := range s {
    //   info(id, fmt.Sprintf("[SERVER] Iteration=%v; Read byte2=%b byte10=%d, symbol=%c, symbol_quote=%q", i, v, v, v, v));
    // }
  }
  bs.begin = start
  bs.end = time.Now().UnixNano()

  for {
    if bs.getClientStopped() {
      info(id, "[SERVER] The server received a signal from the client and shut down")

      err = conn.Close()
      info(id, fmt.Sprintf("[SERVER] Connection was closed by signal from client: %v", conn))
      if err != nil {
        panic(err)
      }

      err = listener.Close()
      info(id, fmt.Sprintf("[SERVER] Listener was closed by signal from client: %v", listener))
      if err != nil {
        panic(err)
      }

      bs.setServerStopped()
      sl.removeServer(sc)
      info(id, "[SERVER][FINISH]")

      break
    }
    time.Sleep(1 * time.Millisecond)
  }
}

func ex2TCPClient(bs *benchState, wg *sync.WaitGroup, maxIt int, bufSize int, id int, s []byte, ip string) {
  defer func() {
    for {
      time.Sleep(1 * time.Millisecond)
      if bs.getServerStopped() {
        wg.Done()
        info(id, "[CLIENT] Call wg.Done().")
        break
      }
    }
  }()

  var err error = nil
  var conn *net.TCPConn
  var connErr error = nil
  var runningServer bool = false
  var cc clientConf
  cc.buildTcp(id, ip)
  info(id, fmt.Sprintf("[CLIENT] cc.ipAddr.IsLoopback()=%v", cc.ipAddr.IsLoopback()))

  for {
    info(id, "[CLIENT] Try again later...")
    time.Sleep(1 * time.Millisecond)

    if !runningServer && bs.getRunningServer() {
      info(id, "[CLIENT] The client received a signal from the server: state=running_server")
      runningServer = true
    } else if runningServer {
      info(id, "[CLIENT] The client received a signal from the server: state=running_server")
    } else {
      info(id, "[CLIENT] The client not received a signal from the server: state=running_server")
      continue
    }

    conn, connErr = net.DialTCP(cc.addrType, nil, cc.tcpAddr)
    if connErr != nil {
      info(id, fmt.Sprintf("[CLIENT] Error connection: %v", connErr))
      connErr = nil
      continue
    }

    if runningServer && connErr == nil {
      info(id, "[CLIENT] Get server state=running_server and connection was created.")
      break
    }
  }

  // << 10 // 1024        kb
  // << 20 // 1048576     kb=1mb
  // << 25 // 33554432    kb
  // << 30 // 1073741824  kb
  // GODOC: SetWriteBuffer sets the size of the operating system's transmit buffer associated with the connection.
  // func (*UnixConn) SetWriteBuffer
  // SRC: src/net/sockopt_posix.go
  // func setWriteBuffer(fd *netFD, bytes int) error {
  //   err := fd.pfd.SetsockoptInt(syscall.SOL_SOCKET, syscall.SO_SNDBUF, bytes)
  //   runtime.KeepAlive(fd)
  //   return wrapSyscallError("setsockopt", err)
  // }
  err = conn.SetWriteBuffer(bufSize)
  if err != nil {
    panic(err)
  }
  info(id, fmt.Sprintf("[CLIENT] SetWriteBuffer(%v) executed", bufSize))

  // CloseRead shuts down the reading side of the Unix domain connection. Most callers should just use Close.
  err = conn.CloseRead()
  if err != nil {
    panic(err)
  }
  info(id, "[CLIENT] CloseRead() executed")

  for i := 0; i < maxIt; i++ {
    // info(id, fmt.Sprintf("[CLIENT] Iteration=%v", i))
    // info(id, "[CLIENT] The action Write() will be executed...")

    // var n int
    // info(id, "[CLIENT] Start writing...")

    // n, err = conn.Write(s)
    _, err = conn.Write(s)
    if err != nil {
      panic(err)
    }
    // info(id, fmt.Sprintf("[CLIENT] Write %v bytes. Err=%v", n, err))
  }

  // It is expected that the connection will be closed by the server:
  bs.setClientStopped()
  info(id, "[CLIENT] The client has finished and sent a signal to the server")
  info(id, "[CLIENT][FINISH]")
}

func ex1UnixServerMultiCon(id int, sl *serversList) {
  var sc serverConf
  sc.buildUnix(id)
  sl.addServer(sc)

  defer func(sc serverConf, sl *serversList) {
    sc.removeSockFile()
    sc.removeSockDir()
    sl.removeServer(sc)
  }(sc, sl)

  // Use "EXAMPLE1" and you can use the linux command "nc" to connect multiple times.
  // Shutdown processing of multiple connection is performed correctly.
  // BASH CMD: nc -U nc -U /run/user/1000/sockbench/1/socket_<PID>_1
  // HINT from man nc: "After the file has been transferred, the connection will close automatically."

  // Mount: tmpfs on /run/user/1000 type tmpfs (rw,nosuid,nodev,relatime,size=1610456k,mode=700,uid=1000,gid=100)

  // tmpfs - a virtual memory filesystem The tmpfs facility allows the
  // creation of filesystems whose contents reside in virtual memory. Since
  // the files on such filesystems typically re- side in RAM, file access is
  // extremely fast.

  // net.ListenUnix(addrType): There are three types of unix domain socket.
  // “unix”       corresponds to SOCK_STREAM
  // “unixdomain” corresponds to SOCK_DGRAM
  // “unixpacket” corresponds to SOCK_SEQPACKET
  // man 2 socket

  // > BASH CMD ps -AH -o pid,ppid,spid,c,pcpu,%cpu,cputime,cputimes,%mem,rss,sz,vsz,cmd | grep 7197

  sc.createSockDir()
  defer func() {
    sc.removeSockDir()
    info(id, "[SERVER] Call sc.removeSockDir()")
  }()
  sc.checkSockDir()

  listener, err := net.ListenUnix(sc.addrType, sc.unixAddr)
  if err != nil {
    panic(err)
  }
  defer func() {
    listener.Close()
    info(id, fmt.Sprintf("[SERVER] Listener was closed: %v", listener))
  }()
  info(id, fmt.Sprintf("[SERVER] Listener was created: %v", listener))
  sc.addUnixListener(listener)

  listener.SetUnlinkOnClose(true)
  info(id, "[SERVER] SetUnlinkOnClose(true) executed.")

  // Maximum number of connections. You can choose an infinite loop.
  for {
    info(id, "[SERVER] The action AcceptUnix() will be executed and go into standby mode.")
    info(id, "[SERVER][STANDBY]")
    conn, err := listener.AcceptUnix()
    if err != nil {
      panic(err)
    }
    info(id, fmt.Sprintf("[SERVER][CONN=%v] The action AcceptUnix() for unix socket was executed.", conn))
    defer func() {
      conn.Close()
      info(id, fmt.Sprintf("[SERVER] Connection was closed: %v", conn))
    }()

    // GODOC: SetReadBuffer sets the size of the operating system's receive buffer associated with the connection.
    // func (*UnixConn) SetReadBuffer
    // SRC: src/net/sockopt_posix.go
    // func setReadBuffer(fd *netFD, bytes int) error {
    //   err := fd.pfd.SetsockoptInt(syscall.SOL_SOCKET, syscall.SO_RCVBUF, bytes)
    //   runtime.KeepAlive(fd)
    //   return wrapSyscallError("setsockopt", err)
    // }
    conn.SetReadBuffer(1)
    if err != nil {
      panic(err)
    }
    info(id, fmt.Sprintf("[SERVER][CONN=%v] SetReadBuffer(1) executed", conn))

    go func(c *net.UnixConn) {
      for {
        var buf [1]byte
        n, err := c.Read(buf[:])
        if err != nil {
          // panic(err)
          info(id, fmt.Sprintf("[SERVER][CONN=%v][DISCONNECT] %v", conn, err))
          break
        }
        info(id, fmt.Sprintf("[SERVER][CONN=%v] Read of %v bytes: %q", conn, n, string(buf[:n])))
      }
    }(conn)
  }
}
