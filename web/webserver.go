package web

import (
	"fmt"
	"github.com/eyethereal/go-config"
	"html"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/tomseago/go-eltee"
)

//////////////////////////////////////////////////////////////////////
// Package level logging
// This only needs to be in one file per-package

var log = config.Logger("eltee")

///////////////////

type WebServer struct {
	s *eltee.Server

	sockets []*WebServerSocketClient
}

func PingHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ping at %v", html.EscapeString(r.URL.Path))
}

func NewWebServer(cfg *config.AclNode, s *eltee.Server) *WebServer {
	log.Error("************** Starting web server maybe????")
	ws := &WebServer{
		s: s,
	}

	// The router
	rtr := mux.NewRouter()
	rtr.Path("/ping").HandlerFunc(PingHandler)
	rtr.Path("/websocket").Handler(NewWebSocketHandler(ws))

	rtr.PathPrefix("/").Handler(http.FileServer(http.Dir("../www/root")))

	//http.Handle

	srv := &http.Server{
		Handler: rtr,
		Addr:    cfg.DefChildAsString(":8000", "address"),
	}

	log.Infof("Starting web server on %v", srv.Addr)
	go (func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Errorf("Unable to start web server: %v", err)
			return
		}
	})()

	return ws
}

// Update any control points with values from this input adapter
func (*WebServer) UpdateControlPoints() {

}

// Update this input adapter with any values from the control points
func (*WebServer) ObserveControlPoints() {

}
