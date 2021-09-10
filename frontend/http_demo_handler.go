package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type HttpDemoResultPage struct {
	Stylesheet template.HTML
	CatBody    template.HTML
	Namespace  string
	Cluster    string
	Injection  string
	Url        string
}

var httpDemoResultPageTemplate string = `
<!DOCTYPE html>
<html lang="en">
	<head>
		<title>Frontend</title>
		<meta charset="utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1">
		{{.Stylesheet}}
	</head>
	<body>
		<div class="container">
			<div class="header clearfix navbar navbar-inverse">
				<div class="container">
				<h3>Frontend Sample App</h3>
				<h3>Cluster: {{.Cluster}}</h3>
				<h3>Namespace: {{.Namespace}} [istio-injection={{.Injection}}]</h3>
				<h3>Requested URL: {{.Url}} </h3>
				</div>
			</div>

			{{.CatBody}}
		</div>
	</body>
</html>
`

type HttpDemoHandler struct{}

func (h *HttpDemoHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	queryParams, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		panic(err)
	}
	destination := queryParams["url"][0]
	destination = "http://" + destination
	bearer := "Bearer " + queryParams["auth"][0]

	newReq, err := http.NewRequest("GET", destination, nil)
	newReq.Header.Set("Authorization", bearer)
	newReq.Header.Set("X-B3-Sampled", "1")
	newReq.Header.Set("traceID", req.Header.Get("traceID"))
	newReq.Header.Set("X-B3-Traceid", req.Header.Get("traceID"))

	httpClient := &http.Client{Timeout: time.Second * 5}
	getResp, err := httpClient.Do(newReq)

	if err != nil {
		resp.Header().Add("X-B3-Sampled", "1")
		resp.Header().Add("traceID", req.Header.Get("traceID"))
		resp.Header().Add("X-B3-Traceid", req.Header.Get("traceID"))
		template := template.Must(template.New("errorPageTemplate").Parse(errorPageTemplate))
		err = template.Execute(resp, ErrorPage{
			Stylesheet: stylesheet,
			Error:      err,
		})
		if err != nil {
			panic(err)
		}

		return
	}
	defer getResp.Body.Close()

	readBytes, err := ioutil.ReadAll(getResp.Body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		resp.Write([]byte(fmt.Sprintf("read body failed: %s", err)))
		return
	}

	theTemplate := template.Must(template.New("httpDemoResultPage").Parse(httpDemoResultPageTemplate))
	catBody := template.HTML(string(readBytes))
	err = theTemplate.Execute(resp, HttpDemoResultPage{
		Stylesheet: stylesheet,
		CatBody:    catBody,
		Namespace:  os.Getenv("NAMESPACE"),
		Cluster:    os.Getenv("CLUSTER"),
		Injection:  os.Getenv("INJECTION"),
		Url:        destination,
	})
	if err != nil {
		panic(err)
	}
}
