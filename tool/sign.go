package tool

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
)

func WxpayCalcSign(mReq map[string]interface{}, key string) (sign string) {
	fmt.Println("wechat pay sign calc, API KEY:", key)
	// sort parameter name by ASCII
	sorted_keys := make([]string, 0)
	for k, _ := range mReq {
		sorted_keys = append(sorted_keys, k)
	}
	sort.Strings(sorted_keys)

	//use & connect parameter
	var signStrings string
	for _, k := range sorted_keys {
		fmt.Printf("k=%v, v=%v\n", k, mReq[k])
		value := fmt.Sprintf("%v", mReq[k])
		if value != "" {
			signStrings = signStrings + k + "=" + value + "&"
		}
	}

	//add api_key
	if key != "" {
		signStrings = signStrings + "key=" + key
	}

	//md5 sign and to upper
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(signStrings))
	cipherStr := md5Ctx.Sum(nil)
	upperSign := strings.ToUpper(hex.EncodeToString(cipherStr))
	return upperSign
}

func WxpayVerifySign(needVerifyM map[string]interface{}, sign string) bool {
	signCalc := WxpayCalcSign(needVerifyM, "API_KEY")
	if sign == signCalc {
		fmt.Sprintln("verify success")
		return true
	}
	fmt.Println("verify fail")
	return false
}
