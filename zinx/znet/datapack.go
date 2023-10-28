package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"myServerDemo/zinx/utils"
	"myServerDemo/zinx/ziface"
)

// 封包、拆包的具体模块
type DataPack struct{}

// 拆包 封包实例的一个初始化方法
func NewDataPack() *DataPack {
	return &DataPack{}
}

// 获取包的头的长度方法
func (dp *DataPack) GetHeadLen() uint32 {
	// DataLen uint32（4字节）+ Id uin32（4字节）
	return 8
}

// 封包方法
// datalen|msgId|data
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})

	// 将dataLen 写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
		return nil, err
	}

	// 将MsgId 写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// 将data数据 写进dataBuff中
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

// 拆包方法 （将包的Head信息读出来）之后再更具head信息里的data的长度，再进行一次读
func (dp *DataPack) Unpack([]byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader([]byte{})

	// 只解压head信息，得到dataLen和MsgID
	msg := &Message{}

	// 读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	// 读MsgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	// 判断dataLen是否超出里我们允许的最大包长度
	if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
		return nil, errors.New("too Large msg data recv!")
	}

	return msg, nil
}
