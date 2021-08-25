package main

import "html/template"

type ErrorPage struct {
	Stylesheet template.HTML
	Error      error
}

var errorPageTemplate string = `
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
				</div>
			</div>

			<div class="jumbotron">
				<p><img width="100" height="100" src="https://raw.githubusercontent.com/adamzwickey/cf-networking-examples/master/frontend/err.png" /></p>
				<p class="lead">request failed: {{.Error}}</p>
			</div>
		</div>
	</body>
</html>
`
