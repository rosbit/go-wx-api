package wxapi

import (
	"io"
	"github.com/rosbit/go-wx-api/auth"
	"github.com/rosbit/go-wx-api/msg"
	"github.com/rosbit/go-wx-api/log"
)

func InitWxAPI(workerNum int, logger io.Writer) {
	wxlog.SetLogger(logger)
	wxauth.StartAuthThreads(workerNum)
	wxmsg.StartWxMsgParsers(workerNum)
}
