/**
 * @file: inventory.go
 * @time: 2022-11-05 10:45
 * @Author: jack
 * @Email: 793936517@qq.com
 * @desc:
 **/

package handler

import (
	"context"
	"fmt"
	goredislib "github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm/clause"
	"mall_srvs/inventory_srv/global"
	"mall_srvs/inventory_srv/model"
	"mall_srvs/inventory_srv/proto"
	"sync"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
}

func (*InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	//设置库存， 如果我要更新库存
	var inv model.Inventory
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num

	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	if result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "没有库存信息")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

var m sync.Mutex

func (*InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1

	//这里的变量哪些可以放到global中， redis的配置是否应该在nacos中
	client := goredislib.NewClient(&goredislib.Options{
		Addr:     "127.0.0.1:6379",
		Password: "topsky",
	})
	pool := goredis.NewPool(client) // or, pool := redigo.NewPool(...)

	// Create an instance of redisync to be used to obtain a mutual exclusion
	// lock.
	rs := redsync.New(pool)

	tx := global.DB.Begin()
	//m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback() //回滚之前的操作
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}
		//for {
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		//zookeeper的分布式锁 -

		fmt.Println("开始获取锁")
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		fmt.Println("获取锁成功")
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
		fmt.Println("开始释放锁")
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "解放redis分布式锁异常")
		}
		fmt.Println("释放锁成功")
		// update inventory set stocks =stocks - 1, version=version+1 where goods= goods and version=version
		// 这种写法有瑕疵， 为什么？
		// 零值 对于int 类型来说，默认值是0，这种会被gorm给忽略掉
		//if result := tx.Model(&model.Inventory{}).Where("goods=? and version=?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
		//if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods=? and version=?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
		//	zap.S().Info("库存扣减失败")
		//} else {
		//	break
		//}
		//}
		//tx.Save(&inv)
	}
	tx.Commit() // 需要自己动手提交操作
	//m.Unlock()  // 释放锁
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) Sell2(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1
	tx := global.DB.Begin()
	// 使用Mutex互斥锁
	m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要自己动手提交操作
	m.Unlock()  // 释放锁
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) Sell3(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 使用mysql悲观锁 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1
	tx := global.DB.Begin()
	//m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks -= goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要自己动手提交操作
	//m.Unlock()  // 释放锁
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) Sell4(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1
	// 使用mysql乐观锁
	tx := global.DB.Begin()
	//m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		//if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
		//	tx.Rollback() //回滚之前的操作
		//	return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		//}
		for {
			if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
				tx.Rollback() //回滚之前的操作
				return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
			}
			//判断库存是否充足
			if inv.Stocks < goodInfo.Num {
				tx.Rollback() //回滚之前的操作
				return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
			}
			//扣减， 会出现数据不一致的问题 - 锁，分布式锁
			inv.Stocks -= goodInfo.Num
			// update inventory set stocks =stocks - 1, version=version+1 where goods= goods and version=version
			// 这种写法有瑕疵， 为什么？
			// 零值 对于int 类型来说，默认值是0，这种会被gorm给忽略掉
			//if result := tx.Model(&model.Inventory{}).Where("goods=? and version=?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
			if result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").Where("goods=? and version=?", goodInfo.GoodsId, inv.Version).Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1}); result.RowsAffected == 0 {
				zap.S().Info("库存扣减失败")
			} else {
				break
			}
		}
		//tx.Save(&inv)
	}
	tx.Commit() // 需要自己动手提交操作
	//m.Unlock()  // 释放锁
	return &emptypb.Empty{}, nil
}

func (*InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	//扣减库存， 本地事务 [1:10,  2:5, 3: 20]
	//数据库基本的一个应用场景：数据库事务
	//并发情况之下 可能会出现超卖 1
	tx := global.DB.Begin()
	m.Lock() //获取锁 这把锁有问题吗？  假设有10w的并发， 这里并不是请求的同一件商品  这个锁就没有问题了吗？
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() //回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "没有库存信息")
		}
		//扣减， 会出现数据不一致的问题 - 锁，分布式锁
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 需要自己动手提交操作
	m.Unlock()  // 释放锁
	return &emptypb.Empty{}, nil
}
