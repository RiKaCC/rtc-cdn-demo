package main

import (
	"net/url"

	"github.com/pion/webrtc/v3"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/connect"}

	ws := webSocketConnection(u)
	pc := createPeerConnection(webrtc.Configuration{})
	dataChannelHandler(pc)
	createOffer(pc)

}
