/*
	@introduce 分布式锁
	@author 陈开广
	@time 2017-08-24
*/
package zookeeper

import (
	"github.com/chrisho/mosquito/helper"
	"errors"
	"strings"
)

// 创建锁节点
func CreatLockNode(path string) error {
	// 过滤空格
	if path = helper.TrimStringSpace(path); path == "" {
		return errors.New("lock node path is empty")
	}
	// 添加斜杠
	if ! strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	// 创建锁目录
	lockPath := zkRootPath  + path
	if ok, _, _ := zkConn.Exists(lockPath); !ok {
		_, err := zkConn.Create(zkRootPath, nil, 0, getAcl())

		if err != nil {
			return err
		}
	}
	return nil
}

