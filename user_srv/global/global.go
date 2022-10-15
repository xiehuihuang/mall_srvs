/**
 * @file: global.go
 * @time: 2022-10-12 18:02
 * @Author: jack
 * @Email: 793936517@qq.com
 * @desc:
 **/

package global

import (
	"gorm.io/gorm"
	"mall_srvs/user_srv/config"
)

var (
	DB           *gorm.DB
	ServerConfig config.ServerConfig
	//ServerConfig *config.ServerConfig = &config.ServerConfig{}
)
