package main

import (
	"encoding/json"
	"github.com/hashicorp/go-retryablehttp"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"
)

type Result struct {
	Status string
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/process", process)
	http.ListenAndServe(":9090", nil)
}

func home(w http.ResponseWriter, r *http.Request)  {
	t := template.Must(template.ParseFiles(filepath.Join("templates", "home.html")))
	t.Execute(w, Result{})
}

func process(w http.ResponseWriter, r *http.Request)  {
	result := makeHttpCall("http://localhost:9091", r.FormValue("coupon"), r.FormValue("cc-number"))
	t := template.Must(template.ParseFiles(filepath.Join("templates", "home.html")))
	t.Execute(w, result)
}

func makeHttpCall(urlMicroService string, coupon string, ccNumber string) Result {
	values := url.Values{}
	values.Add("coupon", coupon)
	values.Add("ccNumber", ccNumber)

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
