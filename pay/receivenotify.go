package pay

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/willxm/wechatpay-go/tool"
)

type WxpayNotifyResp struct {
	Return_code    string `xml:"return_code"`
	Return_msg     string `xml:"return_msg"`
	Appid          string `xml:"appid"`
	Mch_id         string `xml:"mac_id"`
	Device_info    string `xml:"device_info"`
	Nonce_str      string `xml:"nonce_str"`
	Sign           string `xml:"sign"`
	Result_code    string `xml:"result_code"`
	Err_code       string `xml:"err_code"`
	Err_code_des   string `xml:"err_code_des"`
	Openid         string `xml:"openid"`
	Trade_type     string `xml:"trade_type"`
	Bank_type      string `xml:""bank_type`
	Total_fee      int    `xml:"total_fee"`
	Fee_type       string `xml:"fee_type"`
	Cash_fee       int    `xml:"cash_fee"`
	Cash_fee_type  string `xml:"cash_fee_type"`
	Coupon_fee     int    `xml:"coupon_fee"`
	Coupon_count   int    `xml:"coupan_count"`
	Transaction_id string `xml:transaction_id`
	Out_trade_no   string `xml:"out_trade_no"`
	Attach         string `xml:"attach"`
	Time_end       string `xml:"time_end"`
}

type WxpayNotifyReq struct {
	Return_code string `xml:"return_code"`
	Return_msg  string `xml:"return_msg"`
}

func WxpayCallback(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("fail to read body: ", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println("wxpay notify，HTTP Body: ", string(body))
	var nrp WxpayNotifyResp
	err = xml.Unmarshal(body, &nrp)
	if err != nil {
		fmt.Println("decode http body error: ", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var respMap map[string]interface{}
	respMap = make(map[string]interface{}, 0)

	respMap["return_code"] = nrp.Return_code
	respMap["return_msg"] = nrp.Return_msg
	respMap["appid"] = nrp.Appid
	respMap["mch_id"] = nrp.Mch_id
	respMap["device_info"] = nrp.Device_info
	respMap["nonce_str"] = nrp.Nonce_str
	respMap["result_code"] = nrp.Result_code
	respMap["err_code"] = nrp.Err_code
	respMap["err_code_des"] = nrp.Err_code_des
	respMap["openid"] = nrp.Openid
	respMap["trade_type"] = nrp.Trade_type
	respMap["bank_type"] = nrp.Bank_type
	respMap["total_fee"] = nrp.Total_fee
	respMap["fee_type"] = nrp.Fee_type
	respMap["cash_fee"] = nrp.Cash_fee
	respMap["cash_fee_type"] = nrp.Cash_fee_type
	respMap["coupon_fee"] = nrp.Coupon_fee
	respMap["coupon_count"] = nrp.Coupon_count
	respMap["transacation_id"] = nrp.Transaction_id
	respMap["out_trade_no"] = nrp.Out_trade_no
	respMap["attach"] = nrp.Attach
	respMap["time_end"] = nrp.Time_end

	var nrq WxpayNotifyReq

	if tool.WxpayVerifySign(respMap, nrp.Sign) {
		//TODO save to database and etc...
		nrq.Return_code = "SUCCESS"
		nrq.Return_msg = "OK"
	} else {
		nrq.Return_code = "FAIL"
		nrq.Return_msg = "failed to verify sign"
	}

	bytes, err := xml.Marshal(nrq)
	strReq := strings.Replace(string(bytes), "WxpayNotifyReq", "xml", -1)
	if err != nil {
		fmt.Println("encoding xml error: ", err)
		http.Error(w.(http.ResponseWriter), http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	w.(http.ResponseWriter).WriteHeader(http.StatusOK)
	fmt.Fprint(w.(http.ResponseWriter), strReq)

}
