package main

import (
	"fmt"

	galago ".."
)

// This is the main entry point for your REST API. You should use the
// galago.NewAppFromCLI() function if you wish to allow run-time
// configuration of the application. Otherwise, you can create it
// yourself. You can always customize the application generated by
// the CLI further until it is to your liking.
func main() {
	app := galago.NewAppFromCLI()

	/*
	   Set things like rate limits, custom serializers,
	   and global Middleware here.
	*/

	// You can add multiple controllers to your app
	app.AddController(DemoController())
	app.Listen()
}

// DemoController returns the Controller used for this Demo.
func DemoController() *galago.Controller {
	controller := galago.NewController()

	// Each route here takes a method, path and function. The function
	// will be executed for each request matching the specified
	// path and method.
	//
	// See each function applied to a route in this file for more
	// detailed information on what it does.

	// Note that in this example we use a route parameter called name.
	// This can later be retrieved using request.GetField("name")
	controller.AddRoute(
		galago.NewRoute("GET", "example/field/{name}", ExampleField),
	)

	controller.AddRoute(
		galago.NewRoute("GET", "example/header", ExampleHeader),
	)

	controller.AddRoute(
		galago.NewRoute("GET", "example/download", ExampleDownload),
	)

	controller.AddRoute(
		galago.NewRoute("GET", "example/data", ExampleData),
	)

	controller.AddRoute(
		galago.NewRoute("GET", "example/query", ExampleQuery),
	)

	// In this example we create a route and apply a middleware to it.
	// This middleware will set a header on the response as it is
	// leaving the framework.
	//
	// Middleware can be applied to a specific route, or globally to
	// the App itself. Here we are applying it to a route.
	controller.AddRoute(galago.NewRoute(
		"GET", "example/middleware", ExampleMiddleware,
	).AddMiddleware(galago.Middleware{
		After: func(response *galago.Response) {
			response.SetHeader("Test", "It's set")
		},
	}))

	// You can add middleware globally to a Controller
	controller.AddMiddleware(
		galago.Middleware{
			After: func(response *galago.Response) {
				response.SetHeader("Controller-Middleware", "Yes")
			},
		},
	)

	// You can also apply middleware to multiple routes in a
	// controller at once
	controller.AddMiddlewareFor(
		[]string{"example/data", "example/field/{name}"},
		galago.Middleware{
			After: func(response *galago.Response) {
				response.SetHeader("MultiRoute-Middleware", "Yes")
			},
		},
	)

	return controller
}

// ExampleField will take the last route parameter with the key
// name and return it to the user in a message.
//
// curl http://localhost:8080/example/header
func ExampleField(request galago.Request) *galago.Response {
	return galago.NewResponse(200, map[string]interface{}{
		"message": fmt.Sprintf(
			"Your name is %s", *request.GetField("name")),
	})
}

// ExampleHeader will return a response with the Test Header set.
//
// curl -vvvv http://localhost:8080/example/header
func ExampleHeader(request galago.Request) *galago.Response {
	return galago.NewResponse(200, map[string]interface{}{
		"success": true,
	}).SetHeader("Test", "It Worked!")
}

// ExampleDownload will operate as a normal file download by setting
// the Content-Disposition header. Try opening it in your browser.
//
// http://localhost:8080/example/download
func ExampleDownload(request galago.Request) *galago.Response {
	return galago.NewResponse(
		200, galago.DownloadSerializer().MakeRawData(
			"Hello, World!",
		),
	).MakeDownload("hello_world.txt")
}

// ExampleData will take the input data passed to the request and
// parse it, returning the value.
//
// It expects the input data to be in JSON format.
//
// { "client" : { "name" : "John Doe" } }
//
// curl -H'Content-Type: application/json' \
//      -d'{"client":{"name":"John Doe"}}' \
//      http://localhost:8080/example/data
func ExampleData(request galago.Request) *galago.Response {
	name := request.GetData("client.name")

	if name != nil {
		return galago.NewResponse(200, map[string]interface{}{
			"your_name_is": name.(string),
		})
	}

	return galago.NewResponse(400, map[string]interface{}{
		"error": "Missing `client.name` data key/value pair.",
	})
}

// ExampleQuery will look for the query parameter named "name" and
// return it to the user.
//
// curl 'http://localhost:8080/example/query?name=John%20Doe'
func ExampleQuery(request galago.Request) *galago.Response {
	name := request.GetQuery("name")

	if name != nil {
		return galago.NewResponse(200, map[string]interface{}{
			"your_name_is": *name,
		})
	}

	return galago.NewResponse(400, map[string]interface{}{
		"error": "Missing `name` query parameter.",
	})
}

// ExampleMiddleware will use middleware to apply a header to any
// request passed to this route, as it leaves the framework.
//
// curl -vvvv http://localhost:8080/example/middleware
func ExampleMiddleware(request galago.Request) *galago.Response {
	return galago.NewResponse(200, map[string]interface{}{
		"message": "The middleware should have added the Test header",
	})
}
