package main

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	http.ListenAndServe(":9091", nil)
}

func home(w http.ResponseWriter, r *http.Request)  {
	coupon := r.PostFormValue("coupon")
	ccNumber := r.PostFormValue("ccNumber")
	resultCoupon := makeHttpCall("http://localhost:9092", coupon)
	result := Result{Status: "Negado"}

	if ccNumber == "1" {
		result.Status = "Aprovado"
	}

	if resultCoupon.Status == "Invalido" {
		result.Status = "Invalido cupom"
	}

	jsonData, err := json.Marshal(result)

	if err != nil {
		log.Fatal("Erro processing json")
	}

	fmt.Fprintf(w, string(jsonData))
}

func makeHttpCall(urlMicroService string, coupon string) Result {
	values := url.Values{}
	values.Add("coupon", coupon)

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 5
	res, err := retryClient.PostForm(urlMicroService, values);

	if err != nil {
		return Result{Status: "Servidor fora do ar!"}
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal("Result not valid")
	}

	result := Result{}
	json.Unmarshal(data, &result)

	return result
}
