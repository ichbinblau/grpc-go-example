// The server demonstrates how to consume and validate OAuth2 tokens provided by
// clients for each RPC.
package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/lixd/grpc-go-example/data"
	pb "github.com/lixd/grpc-go-example/features/proto/echo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// var (
// 	errMissingMetadata = status.Errorf(codes.InvalidArgument, "missing metadata")
// 	errInvalidToken    = status.Errorf(codes.Unauthenticated, "invalid token")
// )

var port = flag.Int("port", 50052, "the port to serve on")

func main() {
	flag.Parse()

	// cert, err := tls.LoadX509KeyPair(data.Path("x509/na-server.crt"), data.Path("x509/na-server.key"))
	// if err != nil {
	// 	log.Fatalf("failed to load key pair: %s", err)
	// }

	// s := grpc.NewServer(grpc.UnaryInterceptor(ensureValidToken), grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	// s := grpc.NewServer(grpc.UnaryInterceptor(myEnsureValidToken), grpc.Creds(credentials.NewServerTLSFromCert(&cert)))
	s := grpc.NewServer(grpc.Creds(LoadKeyPair()))

	pb.RegisterEchoServer(s, &ecServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	// lis, err := net.Listen("unix", "/tmp/test.sock")

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Serving gRPC on 0.0.0.0" + fmt.Sprintf(":%d", *port))
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func LoadKeyPair() credentials.TransportCredentials {
	certificate, err := tls.LoadX509KeyPair(data.Path("x509/na-server.crt"), data.Path("x509/na-server.key"))
	if err != nil {
		panic("failed to load server certification: " + err.Error())
	}

	data, err := os.ReadFile(data.Path("x509/root-ca.crt"))
	if err != nil {
		panic("failed to load CA file: " + err.Error())
	}

	capool := x509.NewCertPool()
	if !capool.AppendCertsFromPEM(data) {
		panic("can't add ca cert")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{certificate},
		ClientCAs:    capool,
	}
	return credentials.NewTLS(tlsConfig)
}

type ecServer struct {
	pb.UnimplementedEchoServer
}

func (s *ecServer) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	return &pb.EchoResponse{Message: req.Message}, nil
}

// valid 校验认证信息有效性。
// func valid(authorization []string) bool {
// 	if len(authorization) < 1 {
// 		return false
// 	}
// 	token := strings.TrimPrefix(authorization[0], "Bearer ")
// 	return token == "some-secret-token"
// }

// // ensureValidToken 用于校验 token 有效性的一元拦截器。
// func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	// 如果 token不存在或者无效，直接返回错误，否则就调用真正的RPC方法。
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return nil, errMissingMetadata
// 	}
// 	if !valid(md["authorization"]) {
// 		return nil, errInvalidToken
// 	}
// 	return handler(ctx, req)
// }

// // myEnsureValidToken 自定义 token 校验
// func myEnsureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
// 	// 如果返回err不为nil则说明token验证未通过
// 	err := authentication.IsValidAuth(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return handler(ctx, req)
// }

// func streamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
// 	// authentication (token verification)
// 	md, ok := metadata.FromIncomingContext(ss.Context())
// 	if !ok {
// 		return errMissingMetadata
// 	}
// 	if !valid(md["authorization"]) {
// 		return errInvalidToken
// 	}

// 	err := handler(srv, ss)
// 	return err
// }
