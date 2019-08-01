package southbound


type RfidNotifyRequest struct {
	Id string   `json:"id"`
	MsgType string `json:"msgType"`
	Request RequestType `json:"request"`
}

type RequestType struct {
	Header ReqHeader `json:"header"`
	Body ReqBody  `json:"body"`
}

type ReqHeader struct {
	RequestId string  `json:"requestId"`
	ClientCode string  `json:"clientCode"`
	UserId string  `json:"userId"`
	UserKey string  `json:"userKey"`
	Version string  `json:"version"`
}
type ReqBody struct {
	RfidCode string  `json:"rfidCode"`
	PositionCode string  `json:"positionCode"`
	Direction string  `json:"direction"`
}

type RfidNotifyResponse struct {
	Id string  `json:"id"`
	MsgType string  `json:"msgType"`
	Response ResponseType  `json:"response"`
}

type ResponseType struct {
	Header RespHeader  `json:"header"`
}

type RespHeader struct {
	ResponseId string  `json:"responseId"`
	Code int  `json:"code"`
	Msg string   `json:"msg"`
}
