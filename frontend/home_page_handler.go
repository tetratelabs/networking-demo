package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"os"
)

type HomePage struct {
	Stylesheet  template.HTML
	Cachebuster int
	Namespace   string
	Cluster     string
	Injection   string
}

var homePageTemplate string = `
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
				</div>
			</div>

			<div class="jumbotron">
				<form action="/proxy/" method="get" class="form-inline">
					<div class="row">
					<h2>HTTP Test</h2>
						<div class=".col-md-4.col-md-offset-4">
							<div class="form-group">
								<label for="url">Backend HTTP URL&nbsp;&nbsp;&nbsp;&nbsp;</label>
								<input type="text" name="url" class="form-control" placeholder="productpage:9080">
							</div>
						</div>
						<div class=".col-md-4.col-md-offset-4">
							<div class="form-group">
								<label for="auth">Backend Auth Header</label>
								<input type="text" name="auth" class="form-control" placeholder="Authorization: Bearer">
							</div>
						</div>
						<div class=".col-md-4.col-md-offset-4">
							<input type="hidden" name="cachebuster" value="{{.Cachebuster}}">
							<button type="submit" class="btn btn-default">Submit</button>
						</div>
  					</div>
				</form>
			</div>

			<!-- <div class="jumbotron">
				<form action="/udp-test/" method="get" class="form-inline">
					<div class="row">
						<h2>UDP Test</h2>
						<div class=".col-md-4.col-md-offset-4">
				  		<div class="form-group">
								<label for="url">Backend UDP Server Address</label>
								<input type="text" name="url" class="form-control" placeholder="appName.apps.internal:9001">
							</div>
							<br>
				  		<div class="form-group">
								<label for="message">Message</label>
								<input type="text" name="message" class="form-control" placeholder="hello world">
							</div>
							<input type="hidden" name="cachebuster" value="{{.Cachebuster}}">
							<button type="submit" class="btn btn-default">Submit</button>
						</div>
  				</div>
				</form>
			</div> -->
		</div>
	</body>
</html>
`

type HomePageHandler struct{}

func (h *HomePageHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Add("X-B3-Sampled", "1")
	resp.Header().Add("traceID", req.Header.Get("traceID"))
	resp.Header().Add("X-B3-Traceid", req.Header.Get("traceID"))
	template := template.Must(template.New("homePage").Parse(homePageTemplate))
	err := template.Execute(resp, HomePage{
		Stylesheet:  stylesheet,
		Cachebuster: rand.Int(),
		Namespace:   os.Getenv("NAMESPACE"),
		Cluster:     os.Getenv("CLUSTER"),
		Injection:   os.Getenv("INJECTION"),
	})
	if err != nil {
		panic(err)
	}
}
