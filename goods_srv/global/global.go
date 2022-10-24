/**
 * @file: global.go
 * @time: 2022-10-22 15:02
 * @Author: jack
 * @Email: 793936517@qq.com
 * @desc:
 **/

package global

import (
	"gorm.io/gorm"
	"mall_srvs/goods_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  *config.NacosConfig
	//ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
