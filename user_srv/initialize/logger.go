/**
 * @file: logger.go
 * @time: 2022-10-12 18:04
 * @Author: xieHuiHuang
 * @Email: 793936517@qq.com
 * @desc:
 **/

package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
