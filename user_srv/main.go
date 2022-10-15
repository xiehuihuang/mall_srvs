/**
 * @file: main.go
 * @time: 2022-10-11 17:11
 * @Author: jack
 * @Email: 793936517@qq.com
 * @desc:
 **/

package main

import (
	"flag"
	"fmt"
	"github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"mall_srvs/user_srv/global"
	"mall_srvs/user_srv/handler"
	"mall_srvs/user_srv/initialize"
	"mall_srvs/user_srv/proto"
	"net"
)

func main() {
	//初始化
	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.ServerConfig)

	// 从global配置文件中user_srv服务的ip及端口号
	IP := flag.String("ip", global.ServerConfig.Host, "ip地址")
	Port := flag.Int("port", global.ServerConfig.Port, "端口号")

	flag.Parse()
	zap.S().Info("ip: ", *IP)
	zap.S().Info("port: ", *Port)
	// 生产环境动态获取服务可用端口号
	//if *Port == 0 {
	//	*Port, _ = utils.GetFreePort()
	//}
	//
	//zap.S().Info("port: ", *Port)

	server := grpc.NewServer()
	proto.RegisterUserServer(server, &handler.UserServer{})
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *IP, *Port))
	//lis, err := net.Listen("tcp", fmt.Sprintf("#{*IP}:#{*Port}"))
	if err != nil {
		panic("failed to listen:" + err.Error())
	}

	//注册服务健康检查
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())

	//服务注册
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", global.ServerConfig.ConsulInfo.Host,
		global.ServerConfig.ConsulInfo.Port)
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	//生成对应的检查对象
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", *IP, *Port), // user_srv grpc服务ip、port
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}
	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	registration.ID = global.ServerConfig.Name
	registration.Port = *Port
	registration.Tags = []string{"jack", "user_srv"}
	registration.Address = *IP // user_srv grpc服务ip
	registration.Check = check
	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic(err)
	}
	err = server.Serve(lis)
	if err != nil {
		panic("failed to start grpc:" + err.Error())
	}
}
