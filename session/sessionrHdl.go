package session

import (
	//"github.com/xtaci/smux"
	"errors"
	"fmt"
	"github.com/OpenIoTHub/utils/io"
	"github.com/OpenIoTHub/utils/models"
	"github.com/OpenIoTHub/utils/msg"
	"net"
	"time"
)

//:TODO 恢复的没有用，为什么会panic，为什么恢复没用
func PanicHandler() {
	fmt.Printf("panic 产生")

}

func sessionHdl(id string, sessionIn net.Listener) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("panic 恢复")
			fmt.Println(err)
			fmt.Println("结束一个explorer的访问")
		}
		if sessionIn != nil {
			err := sessionIn.Close()
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}()
	for {
		conn, err := sessionIn.Accept()
		if err != nil {
			fmt.Printf(err.Error())
			break
		}
		go sessionConnHdl(id, conn)
	}
}

//访问器的登录处理 conn : 访问器 stream ： 内网端
func sessionConnHdl(id string, conn net.Conn) {
	respOk := func() {
		err := msg.WriteMsg(conn, &models.CheckStatusResponse{
			Code:    0,
			Message: "",
		})
		if err != nil {
			fmt.Println(err.Error())
		}
	}
	respNotOk := func(err error) {
		err = msg.WriteMsg(conn, &models.CheckStatusResponse{
			Code:    1,
			Message: err.Error(),
		})
		if err != nil {
			fmt.Println(err.Error())
		}
		conn.Close()
	}
	var workConn net.Conn
	stream, err := GetStream(id)
	if err != nil {
		fmt.Println(err.Error())
		respNotOk(err)
		return
	}
	err = msg.WriteMsg(stream, &models.RequestNewWorkConn{
		Type:   "kcp",
		Config: "",
	})
	if err != nil {
		fmt.Println(err.Error())
		respNotOk(err)
		return
	}
	if _, ok := sessions[id]; ok {
		//超时返回错误
		select {
		case workConn = <-sessions[id].WorkConn:
			respOk()
		case <-time.After(time.Second * 3):
			respNotOk(errors.New("获取内网连接超时"))
			return
		}

	} else {
		fmt.Println(err.Error())
		respNotOk(err)
		return
	}

	go io.Join(workConn, conn)
}
