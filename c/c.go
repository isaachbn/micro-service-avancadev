package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io/ioutil"
	"log"
	"net/http"
)

type Coupon struct {
	Code string
}

type Coupons struct {
	Coupon []Coupon
}

func (receiver Coupons) Check(code string) string  {
	for _, item := range receiver.Coupon {
		if code == item.Code {
			return "Valido"
		}
	}

	return "Invalido"
}

type Result struct {
	Status string
}

var coupons Coupons
var codes []string;

// Busca códigos de promoção no serviço d
func init() {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	res, err := retryClient.Get("http://localhost:9093")

	if err == nil {
		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(data, &codes)

		for _, code := range codes{
			coupons.Coupon = append(coupons.Coupon, Coupon{
				Code: code,
			})
		}
	}
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9092", nil)
}


func home(w http.ResponseWriter, r *http.Request)  {
	coupon := r.PostFormValue("coupon")
	valid := coupons.Check(coupon)
	result := Result{Status: valid}
	jsoResult, err := json.Marshal(result)

	if err != nil {
		log.Fatal("Erro processing json")
	}

	fmt.Fprintf(w, string(jsoResult))
}