package main

type ProtoClientResp struct {
  Speed int `json:"speed"`
  Size int `json:"size"`
  Url string `json:"url"`
}

const urlSink = "http://sink:80/load"
const urlCtrl = "http://ctrl:8080"
