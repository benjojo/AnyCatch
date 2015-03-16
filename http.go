package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"net/http"
)

func StartServer() {
	m := martini.Classic()
	m.Get("/Get", LastPings)
	m.Get("/Send/:ip", SendPing)

	m.Use(func(res http.ResponseWriter, req *http.Request) {
		// if req.Header.Get("X-API-KEY") != "secret123" {
		// res.WriteHeader(http.StatusUnauthorized)
		// }
	})

	// m.Run()
	m.RunOnAddr(":2374")
}

func SendPing(rw http.ResponseWriter, req *http.Request, params martini.Params) {
	SendPingPacket(params["ip"], AnycastIP)
}

func LastPings(rw http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(lastIPs)
	if err != nil {
		http.Error(rw, "Issue in making json", http.StatusInternalServerError)
	}

	rw.Write(b)
}
