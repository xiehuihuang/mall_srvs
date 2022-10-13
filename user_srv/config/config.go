/**
 * @file: config.go
 * @time: 2022-10-12 18:02
 * @Author: xieHuiHuang
 * @Email: 793936517@qq.com
 * @desc:
 **/

package config

type MysqlConfig struct {
	Host     string `mapstructure:"host" json:"host"`
	Port     int    `mapstructure:"port" json:"port"`
	Name     string `mapstructure:"db" json:"db"`
	User     string `mapstructure:"user" json:"user"`
	Password string `mapstructure:"password" json:"password"`
}

type ServerConfig struct {
	MysqlInfo MysqlConfig `mapstructure:"mysql" json:"mysql"`
}
