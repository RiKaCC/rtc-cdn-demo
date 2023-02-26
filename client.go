package main

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func main() {
	u := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	ws, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dail failed", err)
	}

	defer ws.Close()

	config := webrtc.Configuration{}
	pc, err := webrtc.NewPeerConnection(config)
	if err != nil {
		log.Fatal("new peer connection failed: ", err)
	}

	defer pc.Close()

	// create data channel
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

	// Create an offer and set it as the local description
	offer, err := pc.CreateOffer(nil)
	if err != nil {
		log.Printf("failed to create offer: %s", err)
		return
	}
	if err = pc.SetLocalDescription(offer); err != nil {
		log.Printf("failed to set local description: %s", err)
		return
	}

	// end offer
	if err = ws.WriteJSON(offer); err != nil {
		log.Fatal("client send offer failed ", err)
	}

	_, msg, err := ws.ReadMessage()
	if err != nil {
		log.Fatal("read message failed ", err)
	}

	// handle the answer
	var answer webrtc.SessionDescription
	sdp, err := answer.Unmarshal()
	if err != nil {
		log.Fatal("unmarshal answer failed ", err)
	}

}

func chunkMessage(msg webrtc.DataChannelMessage) {
	fmt.Printf("Message from DataChannel : message data '%s'\n", string(msg.Data))
}

func main1() {
	// 创建一个新的PeerConnection对象
	config := webrtc.Configuration{}
	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}

	// 注册onTrack函数，当收到远程音视频轨道时触发
	peerConnection.OnTrack(func(track *webrtc.TrackRemote, receiver *webrtc.RTPReceiver) {
		fmt.Printf("Got remote track: %s\n", track.Label())
	})

	// 创建一个DataChannel对象
	dataChannel, err := peerConnection.CreateDataChannel("test", nil)
	if err != nil {
		panic(err)
	}

	// 注册onOpen函数，当DataChannel对象打开时触发
	dataChannel.OnOpen(func() {
		fmt.Printf("Data channel '%s'-'%d' open.\n", dataChannel.Label(), dataChannel.ID())
	})

	// 注册onMessage函数，当收到DataChannel对象发送的消息时触发
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Got message on data channel '%s': %s\n", dataChannel.Label(), string(msg.Data))
	})

	// 通过SignalServer获取远端的session描述SDP
	remoteSDP := getSessionDescriptionFromSignalServer()

	// 将远端的SDP设置到PeerConnection对象中
	err = peerConnection.SetRemoteDescription(remoteSDP)
	if err != nil {
		panic(err)
	}

	// 创建一个answer SDP，用于回复远端的offer SDP
	answerSDP, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// 将answer SDP设置到PeerConnection对象中
	err = peerConnection.SetLocalDescription(answerSDP)
	if err != nil {
		panic(err)
	}

	// 将answer SDP发送给远端
	sendAnswerSDPToSignalServer(answerSDP)

	// 等待连接关闭
	<-peerConnection.ConnectionState()
	fmt.Printf("PeerConnection closed\n")
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
