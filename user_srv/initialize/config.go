/**
 * @file: config.go
 * @time: 2022-10-12 18:04
 * @Author: jack
 * @Email: 793936517@qq.com
 * @desc:
 **/

package initialize

import (
	"fmt"

	"github.com/spf13/viper"
	"mall_srvs/user_srv/global"
)

func GetEnvInfo(env string) bool {
	viper.AutomaticEnv()
	return viper.GetBool(env)
	//刚才设置的环境变量 想要生效 我们必须得重启goland
}

func InitConfig() {
	//从配置文件中读取出对应的配置
	debug := GetEnvInfo("MALL_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("user_srv/%s-pro.yaml", configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("user_srv/%s-debug.yaml", configFilePrefix)
	}

	v := viper.New()
	//文件的路径如何设置
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	//这个对象如何在其他文件中使用 - 全局变量
	if err := v.Unmarshal(&global.ServerConfig); err != nil {
		panic(err)
	}
}
