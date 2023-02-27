package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func webSocketConnection(u url.URL) *websocket.Conn {
	log.Printf("connecting to %s", u.String())
	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dail failed", err)
	}
	defer ws.Close()

	return ws
}

func createPeerConnection(config webrtc.Configuration) *webrtc.PeerConnection {
	log.Println("create new peer connection")
	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal("new peer connection failed: ", err)
	}
	defer func() {
		if err = pc.Close(); err != nil {
			fmt.Printf("cannot close peerConnection: %v\n", err)
		}
	}()

	return pc
}

func createOffer(pc *webrtc.PeerConnection) {
	log.Println("create offer")
	offer, err := pc.CreateOffer(nil)
	if err != nil {
		log.Printf("failed to create offer: %s", err)
		return
	}
	if err = pc.SetLocalDescription(offer); err != nil {
		log.Printf("failed to set local description: %s", err)
		return
	}
}

func createAnswer(pc *webrtc.PeerConnection) {
	log.Println("create answer")
	answer, err := pc.CreateAnswer(nil)
	if err != nil {
		log.Printf("failed to create answer: %s", err)
		return
	}
	if err = pc.SetLocalDescription(answer); err != nil {
		log.Printf("failed to set local description: %s", err)
		return
	}
}

func dataChannelHandler(pc *webrtc.PeerConnection) {
	log.Println("create new data channel")
	dataCh, err := pc.CreateDataChannel("dataChannel", nil)
	if err != nil {
		log.Fatal("create data channel failed: ", err)
	}

	dataCh.OnOpen(func() {
		log.Println("data channel opened")
	})

	dataCh.OnClose(func() {
		log.Println("data channel closed")
	})

	dataCh.OnMessage(func(msg webrtc.DataChannelMessage) {
		chunkMessage(msg)
	})
}

func chunkMessage(msg webrtc.DataChannelMessage) {
	fmt.Printf("Message from DataChannel : message data '%s'\n", string(msg.Data))
}

// 从信令服务器获取远端的session描述SDP
func getSessionDescriptionFromSignalServer() webrtc.SessionDescription {
	// TODO: 从信令服务器获取远端的session描述SDP
	return webrtc.SessionDescription{}
}

// 将answer SDP发送给信令服务器
func sendAnswerSDPToSignalServer(answerSDP webrtc.SessionDescription) {
	// TODO: 将answer SDP发送给信令服务器
}
