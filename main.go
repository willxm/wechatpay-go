package main

import (
	"log"
	"net/http"

	"github.com/willxm/wechatpay-go/pay"
)

func main() {
	http.HandleFunc("/unifiedorder", pay.UnifiedOrder)
	http.HandleFunc("/wxpaynotify", pay.WxpayCallback)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
