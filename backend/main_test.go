package main

import (
	"sync"
	"testing"
)

func TestOrderPriority(t *testing.T) {
	m := &OrderManager{
		vipQueue:    []*Order{},
		normalQueue: []*Order{},
		completed:   []*Order{},
		processing:  make(map[int]*Order),
		robots:      make(map[int]*Robot),
		robotOrder:  []int{},
		nextRobotID: 1,
	}
	m.cond = sync.NewCond(&m.mu)

	// 添加普通订单和VIP订单
	normal1 := m.AddNormalOrder()
	vip1 := m.AddVipOrder()
	normal2 := m.AddNormalOrder()
	vip2 := m.AddVipOrder()

	// 获取订单，应该先VIP后普通，且VIP内部FIFO，普通内部FIFO
	// 模拟机器人获取
	order := m.GetOrder(1, nil)
	if order.ID != vip1.ID {
		t.Errorf("期望 VIP1，得到 %s", order.ID)
	}
	order = m.GetOrder(1, nil)
	if order.ID != vip2.ID {
		t.Errorf("期望 VIP2，得到 %s", order.ID)
	}
	order = m.GetOrder(1, nil)
	if order.ID != normal1.ID {
		t.Errorf("期望 Normal1，得到 %s", order.ID)
	}
	order = m.GetOrder(1, nil)
	if order.ID != normal2.ID {
		t.Errorf("期望 Normal2，得到 %s", order.ID)
	}
}

func TestReturnOrder(t *testing.T) {
	m := &OrderManager{
		vipQueue:    []*Order{},
		normalQueue: []*Order{},
		completed:   []*Order{},
		processing:  make(map[int]*Order),
		robots:      make(map[int]*Robot),
		robotOrder:  []int{},
		nextRobotID: 1,
	}
	m.cond = sync.NewCond(&m.mu)

	order := m.AddNormalOrder()
	m.GetOrder(1, nil) // 取走订单
	m.ReturnOrder(1)   // 归还
	// 检查队列头部是否是刚才的订单
	m.mu.Lock()
	if len(m.normalQueue) != 1 || m.normalQueue[0].ID != order.ID {
		t.Error("归还失败，订单未正确回到队列头部")
	}
	m.mu.Unlock()
}
