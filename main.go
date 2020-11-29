package main

import (
  "crypto/tls"
  "fmt"
  "net"
  "net/http"
  "time"

  "github.com/spf13/cobra"
  "golang.org/x/net/http2"
)

func createClient() *http.Client {
  client := &http.Client{
    Timeout: 10*time.Second,
  }

  client.Transport = &http2.Transport{
    AllowHTTP: true,
    DialTLS: func(netw, addr string, cfg *tls.Config) (net.Conn, error) {
        return net.DialTimeout(netw, addr, 10*time.Second)
  }}
  return client
}



func createSink() *cobra.Command {
  var sink = &cobra.Command{
    Use:   "sink",
    Short: "run sink",
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("sink")
      sink()
    },
  }
  return sink
}


func createCli() *cobra.Command {
  var client = &cobra.Command{
    Use:   "client",
    Short: "run client",
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("client")
      client()
    },
  }
  return client
}
func createCtrl() *cobra.Command {
  var client = &cobra.Command{
    Use:   "ctrl",
    Short: "run ctrl",
    Run: func(cmd *cobra.Command, args []string) {
      fmt.Println("ctrl")
      ctrl()
    },
  }
  return client
}

func main() {

  var rootCmd = &cobra.Command{
    Use:   "cpefr",
    Short: "cperf",
  }
  rootCmd.AddCommand(createCli())
  rootCmd.AddCommand(createSink())
  rootCmd.AddCommand(createCtrl())

  if err := rootCmd.Execute(); err != nil {
    fmt.Println(err)
  }

}