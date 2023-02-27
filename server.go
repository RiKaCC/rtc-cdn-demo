package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

func main() {
	http.HandleFunc("/rtc", handleConnect)
	log.Println("Listen on http://localhost:7880/rtc")
	err := http.ListenAndServe(":7880", nil)
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
	pc := createPeerConnection(rtcConfig)

	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		log.Printf("Got ICE Candidate %#v", candidate)
		if candidate == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := pc.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, candidate)
		} else if onICECandidateErr := signalCandidate() {

		}

	})
}

func signalCandidate(addr string, c *webrtc.ICECandidate) error {
	payload := []byte(c.ToJSON().Candidate)
	resp, err := http.Post(fmt.Sprintf("http://%s/candidate", addr), "application/json; charset=utf-8", bytes.NewReader(payload))
	if err != nil {
		return err
	}

	if err = resp.Body.Close(); err != nil {
		return err
	}

	return nil
}

type signalMessage struct {
	Type string `json:"type"`
	SDP  string `json:"sdp"`
}

func (p *peer) sendSignal(signalType string, sdp string) {
	log.Printf("sending signal %s", signalType)
	msg := signalMessage{
		Type: signalType,
		SDP:  sdp,
	}

	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("json signal message failed, %s", err.Error())
		return
	}

	if err = p.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		log.Printf("websocket send signal messge failed, %s", err.Error())
		return
	}
}
