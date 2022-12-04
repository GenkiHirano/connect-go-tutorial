package main

import (
	"context"
	"errors"
	"example/gen/greet/v1/greetv1connect"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/durationpb"

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

// // connect-goでエラーをラップした実装例
// func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
// 	if err := ctx.Err(); err != nil {
// 		return nil, err // automatically coded correctly
// 	}
// 	if err := validateGreetRequest(req.Msg); err != nil {
// 		return nil, connect.NewError(connect.CodeInvalidArgument, err)
// 	}
// 	greeting, err := doGreetWork(ctx, req.Msg)
// 	if err != nil {
// 		return nil, connect.NewError(connect.CodeUnknown, err)
// 	}
// 	return connect.NewResponse(&greetv1.GreetResponse{
// 		Greeting: greeting,
// 	}), nil
// }

// Connectヘッダー (HTTPヘッダー)
// func (s *GreetServer) Greet(ctx context.Context, req *connect.Request[greetv1.GreetRequest]) (*connect.Response[greetv1.GreetResponse], error) {
// 	fmt.Println(req.Header().Get("Acme-Tenant-Id"))
// 	res := connect.NewResponse(&greetv1.GreetResponse{})
// 	res.Header().Set("Greet-Version", "v1")
// 	return res, nil
// }

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
}


func newTransientError() error {
	err := connect.NewError(
		connect.CodeUnavailable,
		errors.New("overloaded: back off and retry"),
	)
	retryInfo := &errdetails.RetryInfo{
		RetryDelay: durationpb.New(10*time.Second),
	}
	if detail, detailErr := connect.NewErrorDetail(retryInfo); detailErr == nil {
		err.AddDetail(detail)
	}
	return err
}
