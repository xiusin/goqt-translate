package request

import (
	"net/http"
	"sync"

	"github.com/mozillazg/request"
)

// RequestPool 请求池
var RequestPool = sync.Pool{
	New: func() interface{} {
		return request.NewRequest(new(http.Client))
	},
}

// Get 发起Get请求
func Get(url string) ([]byte, error) {
	req := RequestPool.Get().(*request.Request)
	defer RequestPool.Put(req)
	resp, err := req.Get(url)
	if err != nil {
		return nil, err
	}
	content, err := resp.Content()
	defer resp.Body.Close() // 关闭body 否则内存溢出
	return content, nil
}

// InArray 数组查找
func InArray(idx int, arr []int) bool {
	for _, i := range arr {
		if i == idx {
			return true
		}
	}
	return false
}
