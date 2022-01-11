package main

import (
	"context"
	"fmt"
	pb "my/productInfo.proto.pb"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

type server struct {
	products []*pb.Product
}

func (s *server) AddProduct(ctx context.Context, product *pb.Product) (productId *pb.ProductID, err error) {
	product.Id = uuid.New().String() //product 的 id 生成
	s.products = append(s.products, product)
	productId = &pb.ProductID{
		Value: product.Id,
	}
	return productId, status.New(codes.OK, "").Err()
}
func (s *server) GetProduct(ctx context.Context, Id *pb.ProductID) (*pb.Product, error) {
	for _, product := range s.products {
		if product.Id == Id.Value {
			return product, status.New(codes.OK, "").Err()
		}
	}
	return nil, status.Errorf(codes.NotFound, "数据库中不存在该 id 的商品", Id)
}
func main() {
	port := ":5001"
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("net service is starting at %s\n", port)
	s := grpc.NewServer()
	pb.RegisterProductInfoServer(s, &server{})
	if err = s.Serve(lis); err != nil {
		panic(fmt.Sprintf("tpc web service port:【%s】 launch fail %v", port, err))
	}
}
