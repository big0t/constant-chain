package wire

import (
	"github.com/constant-money/constant-chain/common"
	"time"
)

type IHAVEMsg struct {
	peerID            string
	resourceType      string
	availableResource string
}

type IWANTMsg struct {
	peerID          string
	resourceType    string
	requestResource string
}

type ResMsg struct {
	status       string
	resourceType string
	content      string
}

type Resource struct {
	resourceType      string
	availableResource string
}

type BlockResourceInfo struct {
	from uint64
	to   uint64
	hash *common.Hash
}

type syncProtocol struct {
	PeerID           string
	outIHAVE         chan IHAVEMsg
	outIWANT         chan IWANTMsg
	inIWANT          chan IWANTMsg
	outRes           chan ResMsg
	localResources   map[string]*Resource
	networkResources map[string]map[string]*Resource
	subscribedTopics map[string]bool
}

var syncEngine *syncProtocol = nil

func GetSyncEngine(peerID string) *syncProtocol {
	if syncEngine == nil {
		sE := new(syncProtocol)
		sE.PeerID = peerID
		sE.outIHAVE = make(chan IHAVEMsg)
		sE.outIWANT = make(chan IWANTMsg)

		sE.inIWANT = make(chan IWANTMsg) // somebody else request our resource
		sE.outRes = make(chan ResMsg)    // response to the request
		sE.broadcast()
		sE.requestSubscribeTopic()
		return sE
	}
	return syncEngine
}

func (s *syncProtocol) Subscribe(messageType string) {
	//subscribe to a message type
	s.subscribedTopics[messageType] = true
}

func (s *syncProtocol) UnSubscribe(messageType string) {
	//unsubscribe to a message type
	delete(s.subscribedTopics, messageType)
}

func (s *syncProtocol) Serve() {
	//serve i want request
	go func() {
		for {
			reqMsg := <-s.inIWANT
			if s.localResources[reqMsg.resourceType] != nil {
				//return request message
				//s.outRes <-
			}
		}
	}()

}

func (s *syncProtocol) UpdateResource(resourceType string, content ...string) {
	//update available resource
}

func (s *syncProtocol) broadcast() {
	//periodically broadcast
	t := time.Tick(time.Second * 1000)
	for {
		<-t
		for _, v := range s.localResources {
			newMsg := IHAVEMsg{s.PeerID, v.resourceType, v.availableResource}
			s.outIHAVE <- newMsg
		}
	}
}

func (s *syncProtocol) requestSubscribeTopic() {
	//manage to request from a subscribed topic
	for k, _ := range s.subscribedTopics {
		if s.networkResources[k] != nil {
			//for peerID, availableResource = range s.networkResources[k] {
			//get available resource miss from local resource -> outIWANT
			//}
		}
	}
}

func parseContent(resourceType string, content string) interface{} {
	return nil
}
