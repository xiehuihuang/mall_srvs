/**
 * @file: global.go
 * @time: 2022-11-05 10:44
 * @Author: jack
 * @Email: 793936517@qq.com
 * @desc:
 **/

package global

import (
	"gorm.io/gorm"
	"mall_srvs/inventory_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	NacosConfig  *config.NacosConfig
	//ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
