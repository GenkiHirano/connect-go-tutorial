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
// 	err := connect.NewError(connect.CodeUnknown,errors.New("oh no!"),)
// 	err.Meta().Set("Greet-Version", "v1")
// 	return nil, err
// }

// func call() {
// 	_, err := greetv1connect.NewGreetServiceClient(http.DefaultClient,"https://api.acme.com",
// 	).Greet(context.Background(),connect.NewRequest(&greetv1.GreetRequest{}),)
// 	if connectErr := new(connect.Error); errors.As(err, &connectErr) {
// 		fmt.Println(err.Meta().Get("Greet-Version"))
// 	}
// }

// インターセプタ設定
func newInterCeptors() connect.Option {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// ここでヘッダをセットするなど色々処理を書ける
			req.Header().Set("hoge", "fuga")
			return next(ctx, req)
		})
	}
	return connect.WithInterceptors(connect.UnaryInterceptorFunc(interceptor))
}

func main() {
	greeter := &GreetServer{}
	mux := http.NewServeMux()
	interceptor := newInterCeptors()
	path, handler := greetv1connect.NewGreetServiceHandler(greeter, interceptor)
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
