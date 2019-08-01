package southbound

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/jerryzhengj/tagchk/etc"
	log "github.com/jeanphorn/log4go"
	"io/ioutil"
	"net/http"
	"strings"
)

const(
	DRIECTION_IN  = "RIGHT2LEFT"  //方向 右 --> 左(入库)
	DRIECTION_OUT = "LEFT2RIGHT" //方向 左 --> 右(出库)
)

var deviceNoMap map[string]string

func InitDeviceNoMap(m map[string]string){
	deviceNoMap = m
}

func GetDeviceNo(key string)string{
	key,ok := deviceNoMap[key]
	if !ok {
		return ""
	}
	return key
}

func Send(deviceNo string, epc string, direction string){


	log.LOGGER("App").Info("deviceNo=%s,epc=%s,direction=%s",deviceNo, epc, direction)
    reqMsg := buildRFIDNotifyMsg(deviceNo, epc, direction)
    jsonStr, err := json.Marshal(reqMsg)
    if err != nil {
		log.LOGGER("App").Error("failed to translate RfidNotifyRequest to json:%s",err)
	}

	log.LOGGER("App").Info("send msg=%s",string(jsonStr))
	host := etc.Config.Upstream.Url
	retryNum := etc.Config.Upstream.RetryNum

	respCode,err := sendRFIDNotify(string(jsonStr),host)
	if err != nil ||  respCode != 0 {
		log.LOGGER("App").Error("send msg failed:%s",err)
		for i:=0; i < retryNum; i++ {
			respCode,err = sendRFIDNotify(string(jsonStr),host)
			if err != nil ||  respCode != 0{
				log.LOGGER("App").Error("retry %d send msg failed:%s",i+1,err)
			}else{
				log.LOGGER("App").Error("send msg success after retry %d times ",i+1)
				break
			}
		}
	}

	if respCode != 0 {
		log.LOGGER("App").Error("sed RfidNotifyRequest[%s] failed after retry %d times",string(jsonStr),retryNum)
	}

}

func sendRFIDNotify(msg string, host string) (int,error){
	resp, err := http.Post(host, "application/json;charset=utf-8", bytes.NewReader([]byte(msg)))
	defer func() {
		if resp != nil {
			if resp.Body != nil {
				_ = resp.Body.Close()
			}
		}
	}()
	//发送消息失败，则返回-1
	if err != nil {
		return -1, err
	}

	//读取响应消息失败，则返回-2
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil{
		return -2, err2
	}
	log.LOGGER("App").Info("response msg:%s",string(body))
	str := string(body)
	if strings.Contains(str,"success"){
		return 0, nil
	}else{
		return -1 ,errors.New("server error")
	}
	//noticeResp := RfidNotifyResponse{}
	//err3 :=json.Unmarshal(body, &noticeResp)
	//
	////解析效应消息失败，则返回-3
	//if err3 != nil {
	//	return -3, err3
	//}
	//
	//if noticeResp.Response.Header.Code == 0{
	//	return 0, nil
	//}else{
	//	return -4, errors.New("server error")  //服务端响应代码不为0，则执行失败
	//}
}

func buildRFIDNotifyMsg(deviceNo string, epc string, direction string)RfidNotifyRequest{
	reqMsg := RfidNotifyRequest{
		Id:"",
		MsgType: etc.Config.Upstream.MsgType,
	}

	reqMsg.Request.Header = ReqHeader{
		RequestId :"",
		ClientCode : etc.Config.Upstream.ClientCode,
		UserId : etc.Config.Upstream.UserId,
		UserKey : etc.Config.Upstream.UserKey,
		Version : etc.Config.Upstream.Version,
	}

	reqMsg.Request.Body = ReqBody{
		RfidCode : epc,
		PositionCode : deviceNo,
		Direction : direction,
	}

	return reqMsg

}
