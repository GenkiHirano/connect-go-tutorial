package main

import (
	"context"
	"example/gen/greet/v1/greetv1connect"
	"fmt"
	"log"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	greetv1 "example/gen/greet/v1"
)

type GreetServer struct{}

func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
	log.Println("Request headers: ", req.Header())
	res := connect.NewResponse(&greetv1.GreetResponse{
		Greeting: fmt.Sprintf("Hello, %s!", req.Msg.Name),
	})
	res.Header().Set("Greet-Version", "v1")
	return res, nil
}

func main() {
	greeter := &GreetServer{}
	mux := http.NewServeMux()
	path, handler := greetv1connect.NewGreetServiceHandler(greeter)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// h2cを使うことで、TLSなしでHTTP/2を提供できるようにする
		h2c.NewHandler(mux, &http2.Server{}),
	)

	// api := http.NewServeMux()
	// api.Handle(greetv1connect.NewGreetServiceHandler(&greetServer{}))

	// mux := http.NewServeMux()
	// mux.Handle("/", newHTMLHandler())
	// mux.Handle("/grpc/", http.StripPrefix("/grpc", api))
	// http.ListenAndServe(":http", mux)
}
