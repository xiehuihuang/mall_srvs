/**
 * @file: main.go
 * @time: 2022-11-05 10:42
 * @Author: jack
 * @Email: 793936517@qq.com
 * @desc:
 **/

package main

import (
	"flag"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mall_srvs/inventory_srv/utils/register/consul"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mall_srvs/inventory_srv/global"
	"mall_srvs/inventory_srv/handler"
	"mall_srvs/inventory_srv/initialize"
	"mall_srvs/inventory_srv/proto"
	"mall_srvs/inventory_srv/utils"
	"net"
)

func main() {

	IP := flag.String("ip", "0.0.0.0", "ip地址")
	//Port := flag.Int("port", global.ServerConfig.Port, "端口号")
	Port := flag.Int("port", 50053, "端口号")
	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.ServerConfig)

	// 从global配置文件中user_srv服务的ip及端口号

	flag.Parse()
	zap.S().Info("ip: ", *IP)
	zap.S().Info("port: ", *Port)
	// 生产环境动态获取服务可用端口号
	if *Port == 0 {
		*Port, _ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterInventoryServer(server, &handler.InventoryServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	//lis, err := net.Listen("tcp", fmt.Sprintf("#{*IP}:#{*Port}"))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//启动服务
	go func() {
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()
	//服务注册
	register_client := consul.NewRegistryClient(global.ServerConfig.ConsulInfo.Host, global.ServerConfig.ConsulInfo.Port)
	serviceId := fmt.Sprintf("%s", uuid.NewV4())
	err = register_client.Register(global.ServerConfig.Host, *Port, global.ServerConfig.Name, global.ServerConfig.Tags, serviceId)
	if err != nil {
		zap.S().Panic("服务注册失败:", err.Error())
	}
	zap.S().Debugf("启动服务，端口: %d", *Port)
	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	if err = register_client.DeRegister(serviceId); err != nil {
		zap.S().Info("注销失败:", err.Error())
	} else {
		zap.S().Info("注销成功:")
	}
}
