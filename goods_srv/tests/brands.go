/**
 * @file: brands.go
 * @time: 2022-10-22 15:27
 * @Author: jack
 * @Email: 793936517@qq.com
 * @desc:
 **/

package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc"

	"mall_srvs/goods_srv/proto"
)

var brandClient proto.GoodsClient
var conn *grpc.ClientConn

func Init() {
	var err error
	conn, err = grpc.Dial("127.0.0.1:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	brandClient = proto.NewGoodsClient(conn)
}

func TestGetBrandList() {
	rsp, err := brandClient.BrandList(context.Background(), &proto.BrandFilterRequest{})
	if err != nil {
		panic(err)
	}
	fmt.Println(rsp.Total)
	for _, brand := range rsp.Data {
		fmt.Println(brand.Name)
	}
}

func main() {
	Init()
	TestGetBrandList()

	conn.Close()
}
