package ziface

/**
 * 将客户端请求的链接信息和请求的数据包装到一个Request中
 */
type IRequest struct {
	// 得到当前链接
	GetConnection IConnection

	// 得到请求的消息数据
	GetData []byte
}
