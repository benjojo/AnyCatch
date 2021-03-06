package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"net/http"
)

func StartServer(password string) {
	m := martini.Classic()
	m.Get("/Get", LastPings)
	m.Get("/Send/:ip/:token", SendPing)

	m.Use(func(res http.ResponseWriter, req *http.Request) {
		if req.Header.Get("X-API-KEY") != password {
			res.WriteHeader(http.StatusUnauthorized)
		}
	})

	// m.Run()
	m.RunOnAddr(":2374")
}

func SendPing(rw http.ResponseWriter, req *http.Request, params martini.Params) {
	for i := 0; i < 3; i++ {
		SendPingPacket(params["ip"], AnycastIP, params["token"])
	}
}

func LastPings(rw http.ResponseWriter, req *http.Request) {
	b, err := json.Marshal(lastIPs)
	if err != nil {
		http.Error(rw, "Issue in making json", http.StatusInternalServerError)
	}

	rw.Write(b)
}
