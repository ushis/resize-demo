package main

import (
  "flag"
  "fmt"
  "os"
  "os/signal"
  "resizer"
  "syscall"
)

var addr string

func init() {
  flag.StringVar(&addr, "listen", ":8080", "addr or socket to listen on")
  flag.Usage = usage
}

func main() {
  flag.Parse()

  if flag.NArg() < 1 {
    usage()
    os.Exit(1)
  }
  server, err := resizer.New(flag.Arg(0), addr)

  if err != nil {
    die(err)
  }
  defer server.Stop()

  go log(server.Log)

  go server.Run()

  sig := make(chan os.Signal)
  signal.Notify(sig, syscall.SIGINT)
  <-sig
}

func usage() {
  fmt.Fprintf(os.Stderr, "Usage: %s [options] <root>\n\nOptions:\n", os.Args[0])
  flag.PrintDefaults()
}

func die(err error) {
  fmt.Fprintln(os.Stderr, err)
  os.Exit(1)
}

func log(log <-chan string) {
  for msg := range log {
    fmt.Fprintln(os.Stderr, msg)
  }
}
