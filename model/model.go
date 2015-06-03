package model

type Peer struct {
	Name string `json:"name,omitempty" binding:"required"`
	Url  string `json:"url,omitempty" binding:"required"`
}

type Msg struct {
	Sender    Peer   `json:"sender,omitempty" binding:"required"`
	Recipient Peer   `json:"recipient,omitempty" binding:"required"`
	MsgText   string `json:"msgText,omitempty" binding:"required"`
}

type Response struct {
	Status   bool   `json:"status"`
	ErrorMsg string `json:"errorMsg"`
	Peers    []Peer `json:"peers"`
}

type Event struct {
	Type         EventType `json:"type"`
	Msg          Msg       `json:"msg"`
	Peer         Peer      `json:"peer"`
	Peers        []Peer    `json:"peers"`
	ErrorMsg     string    `json:"errorMsg"`
	CompletedMsg string    `json:"completedMsg"`
}

type EventType int

const (
	EventTypeUnknown          EventType = 0
	EventTypeRegister                   = 1
	EventTypeSendMsg                    = 2
	EventTypeUnRegister                 = 3
	EventTypeNewPeersReceived           = 20
	EventTypeMsgReceived                = 21
	EventTypeError                      = 30
	EventTypeCompleted                  = 31
)
