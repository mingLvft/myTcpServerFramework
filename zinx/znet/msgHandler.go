package znet

import (
	"fmt"
	"myServerDemo/zinx/ziface"
	"strconv"
)

/*
消息处理模块的实现
*/
type MsgHandle struct {
	// 存放每个MsgID 所对应的处理方法
	Apis map[uint32]ziface.IRouter
}

// 初始化/创建MsgHandle方法
func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]ziface.IRouter),
	}
}

// 调度/执行对应的Router消息处理方法
func (mh *MsgHandle) DoMsgHandler(request ziface.IRequest) {
	// 1 从Request中找到MsgID
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID = ", request.GetMsgID(), " is NOT FOUND! Need Register!")
	}
	// 2 根据MsgID 调度对应的router业务即可
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

// 为消息添加具体的处理逻辑
func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	// 1 判断 当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// id已经注册了
		panic("repeat api, msgID = " + strconv.Itoa(int(msgID)))
	}
	// 2 添加msg与APi的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgID = ", msgID, " success")
}
