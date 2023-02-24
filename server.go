package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func main() {
	http.HandleFunc("/connect", handleConnect)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

type peer struct {
	conn *websocket.Conn
	pc   *webrtc.PeerConnection
}

func handleConnect(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer conn.Close()

	p := &peer{conn: conn}

	rtcConfig := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			webrtc.ICEServer{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	pc, err := webrtc.NewPeerConnection(rtcConfig)
	if err != nil {
		log.Printf("failed to create peerconnection: %s", err.Error())
		return
	}

	defer pc.Close()

}
