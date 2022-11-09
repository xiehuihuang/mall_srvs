/**
 * @file: logger.go
 * @time: 2022-11-05 10:49
 * @Author: jack
 * @Email: 793936517@qq.com
 * @desc:
 **/

package initialize

import "go.uber.org/zap"

func InitLogger() {
	logger, _ := zap.NewDevelopment()
	zap.ReplaceGlobals(logger)
}
