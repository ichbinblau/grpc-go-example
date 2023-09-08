// The client demonstrates how to supply an OAuth2 token for every RPC.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lixd/grpc-go-example/data"
	ecpb "github.com/lixd/grpc-go-example/features/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var addr = flag.String("addr", "localhost:50052", "the address to connect to")

// var addr = flag.String("addr", fmt.Sprintf("unix://%s", "/tmp/test.sock"), "the address to connect to")

func callUnaryEcho(client ecpb.EchoClient, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := client.UnaryEcho(ctx, &ecpb.EchoRequest{Message: message})
	if err != nil {
		log.Fatalf("client.UnaryEcho(_) = _, %v: ", err)
	}
	fmt.Println("UnaryEcho: ", resp.Message)
}

func main() {
	flag.Parse()

	// 构建一个 PerRPCCredentials。
	// perRPC := oauth.NewOauthAccess(fetchToken())
	//myAuth := authentication.NewMyAuth()
	// creds, err := credentials.NewClientTLSFromFile(data.Path("x509/nodeagent-ca.crt"), "")
	// if err != nil {
	// 	log.Fatalf("failed to load credentials: %v", err)
	// }

	// conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(perRPC))
	//conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(myAuth))
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(LoadKeyPair()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := ecpb.NewEchoClient(conn)

	callUnaryEcho(client, "hello world")
}

func LoadKeyPair() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair(data.Path("x509/aggr-server.crt"), data.Path("x509/aggr-server.key"))
	if err != nil {
		panic("Load client certification failed: " + err.Error())
	}

	ca, err := os.ReadFile(data.Path("x509/root-ca.crt"))
	if err != nil {
		panic("can't read ca file")
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(ca) {
		panic("invalid CA file")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      capool,
	}

	return credentials.NewTLS(tlsConfig)
}

// fetchToken 获取授权信息
// func fetchToken() *oauth2.Token {
// 	return &oauth2.Token{
// 		AccessToken: "some-secret-token",
// 	}
// }
