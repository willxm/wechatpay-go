package pay

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"wxpay/tool"
)

type UnifyOrderReq struct {
	Appid            string `xml:"appid"`
	Mch_id           string `xml:"mch_id"`
	Nonce_str        string `xml:"nonce_str"`
	Sign             string `xml:"sign"`
	Body             string `xml:"body"`
	Out_trade_no     string `xml:"out_trade_no"`
	Total_fee        int    `xml:"total_fee"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
	Notify_url       string `xml:"notify_url"`
	Trade_type       string `xml:"trade_type"`
}

type UnifyOrderResp struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
	Appid       string `xml:"appid"`
	Mch_id      string `xml:"mch_id"`
	Nonce_str   string `xml:"nonce_str"`
	Sign        string `xml:"sign"`
	Result_code string `xml:"result_code"`
	Prepay_id   string `xml:"prepay_id"`
	Trade_type  string `xml:"trade_type"`
}

func UnifiedOrder(w http.ResponseWriter, r *http.Request) {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("fail to read body: ", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("app request，HTTP Body: ", string(body))
	var clientReq UnifyOrderReq
	err = json.Unmarshal(body, &clientReq)
	if err != nil {
		fmt.Println("decode http body error: ", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	//test data

	// clientReq.Appid = "app_id"
	// clientReq.Mch_id = "mch_id"
	// clientReq.Nonce_str = "nonce_str"
	// clientReq.Body = "body"
	// clientReq.Out_trade_no = "out_trade_no"
	// clientReq.Total_fee = 10
	// clientReq.Spbill_create_ip = "spbill_create_ip"
	// clientReq.Notify_url = "notify_url"
	// clientReq.Trade_type = "APP"

	var m map[string]interface{}
	m = make(map[string]interface{}, 0)
	m["appid"] = clientReq.Appid
	m["mach_id"] = clientReq.Mch_id
	m["nonce_str"] = clientReq.Nonce_str
	m["body"] = clientReq.Body
	m["out_trade_no"] = clientReq.Out_trade_no
	m["total_fee"] = clientReq.Total_fee
	m["spbill_create_ip"] = clientReq.Spbill_create_ip
	m["notify_url"] = clientReq.Notify_url
	m["trade_type"] = clientReq.Trade_type
	clientReq.Sign = tool.WxpayCalcSign(m, "wxpay_api_key")

	//xml encoding
	bytesReq, err := xml.Marshal(clientReq)
	if err != nil {
		fmt.Println("xml encoding fail: ", err)
		return
	}
	strReq := string(bytesReq)
	strReq = strings.Replace(strReq, "UnifyOrderReq", "xml", -1)
	fmt.Println(strReq)
	bytesReq = []byte(strReq)

	reqURL := "https://api.mch.weixin.qq.com/pay/unifiedorder"
	req, err := http.NewRequest("POST", reqURL, bytes.NewReader(bytesReq))
	if err != nil {
		fmt.Println("new http request fail: ", err)
		return
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println("wxpay api send fail", err)
		return
	}
	bytesResp, _ := ioutil.ReadAll(resp.Body)

	var clinetResp UnifyOrderResp
	_err := xml.Unmarshal(bytesResp, &clinetResp)
	if _err != nil {
		fmt.Println("decode http body error: ", _err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	fmt.Println("wxpay api return data: ", string(bytesResp))
	fmt.Println("encoding data: ", clinetResp)
	jsonBytes, err := json.Marshal(clinetResp)
	if err != nil {
		fmt.Println("encode json error: ", err)
	}
	w.Write(jsonBytes)
}
