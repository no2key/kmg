package netTester

import (
	"net"
	"time"

	"github.com/bronze1man/kmg/kmgNet"
	"github.com/bronze1man/kmg/kmgTime"
)

//client 只读取数据,不写入数据
func readOnly(listenerNewer ListenNewer, Dialer DirectDialer, debug bool) {
	kmgTime.MustNotTimeout(func() {
		toWrite := []byte("hello world")
		listener := listenAccept(listenerNewer, func(c net.Conn) {
			defer c.Close()
			for i := 0; i < 2; i++ {
				_, err := c.Write(toWrite)
				mustNotError(err)
				time.Sleep(time.Microsecond)
			}
		})
		defer listener.Close()
		conn1, err := Dialer()
		mustNotError(err)
		if debug {
			conn1 = kmgNet.NewDebugConn(conn1, "readOnly")
		}
		defer conn1.Close()
		for i := 0; i < 2; i++ {
			mustReadSame(conn1, toWrite)
			time.Sleep(time.Microsecond)
		}
		conn1.Close()

		listener.Close()
	}, time.Second)

}
