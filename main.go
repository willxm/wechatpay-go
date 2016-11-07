package main

import (
	"log"
	"net/http"
	"wxpay/pay"
)

func main() {
	http.HandleFunc("/unifiedorder", pay.UnifiedOrder)
	http.HandleFunc("/wxpaynotify", pay.WxpayCallback)
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
