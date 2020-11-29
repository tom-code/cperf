package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var ctrlSpeed = 1
var ctrlSize = 1000

func ctrlHandler(w http.ResponseWriter, r *http.Request) {
  log.Println(r.URL.Path)
  if r.URL.Path == "/client" {
    cfg := ProtoClientResp {
      Speed: ctrlSpeed,
      Url: urlSink,
      Size: ctrlSize,
    }
    data, err := json.Marshal(&cfg)
    if err != nil {
      log.Println(err)
    }
    w.Write(data)
    w.WriteHeader(200)
  }

  if r.URL.Path == "/set" {
    sp := r.URL.Query().Get("speed")
    if len(sp) > 0 {
      var err error
      ctrlSpeed, err = strconv.Atoi(sp)
      if err == nil {
        log.Printf("speed set to %d\n", ctrlSpeed)
      }
    }
    size := r.URL.Query().Get("size")
    if len(size) > 0 {
      var err error
      ctrlSize, err = strconv.Atoi(size)
      if err == nil {
        log.Printf("size set to %d\n", ctrlSize)
      }
    }
  }
}

func ctrl() {
  h2s := &http2.Server {
  }
  handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    ctrlHandler(w, r)
  })
  h1s := &http.Server {
    Addr:    ":8080",
    Handler: h2c.NewHandler(handler, h2s),
  }
  hs := &http.Server {
    Addr: ":8081",
    Handler: handler,
  }
  go func() {
    hs.ListenAndServe()
  }()
  h1s.ListenAndServe()
}