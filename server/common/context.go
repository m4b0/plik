package common

import (
	"fmt"
	"github.com/root-gg/context"
	"github.com/root-gg/logger"
	"net/http"
	"strings"
)

var rootLogger *logger.Logger = logger.NewLogger()
var rootContext *PlikContext = newRootContext()

func RootContext() *PlikContext {
	return rootContext
}

func Log() *logger.Logger {
	return rootLogger
}

type PlikContext struct {
	*context.Context
	*logger.Logger
}

func newRootContext() (ctx *PlikContext) {
	ctx = new(PlikContext)
	ctx.Context = context.NewContext("ROOT")
	ctx.Logger = rootLogger
	return
}

func NewPlikContext(name string, req *http.Request) (ctx *PlikContext) {
	ctx = new(PlikContext)
	ctx.Context = rootContext.Context.Fork(name).AutoDetach()
	ctx.Logger = rootContext.Logger.Copy()

	// TODO X-FORWARDED-FOR
	remoteAddr := strings.Split(req.RemoteAddr, ":")
	if len(remoteAddr) > 0 {
		ctx.Set("RemoteIp", remoteAddr[0])
	}

	ctx.UpdateLoggerPrefix("")
	return
}

func (ctx *PlikContext) Fork(name string) (fork *PlikContext) {
	fork = new(PlikContext)
	fork.Context = ctx.Context.Fork(name)
	fork.Logger = ctx.Logger.Copy()
	return fork
}

func (ctx *PlikContext) SetUpload(uploadId string) *PlikContext {
	ctx.Set("UploadId", uploadId)
	ctx.UpdateLoggerPrefix("")
	return ctx
}

func (ctx *PlikContext) SetFile(fileName string) *PlikContext {
	ctx.Set("FileName", fileName)
	ctx.UpdateLoggerPrefix("")
	return ctx
}

func (ctx *PlikContext) UpdateLoggerPrefix(prefix string) {
	str := ""
	if ip, ok := ctx.Get("RemoteIp"); ok {
		str += fmt.Sprintf("[%s]", ip)
	}
	if uploadId, ok := ctx.Get("UploadId"); ok {
		str += fmt.Sprintf("[%s]", uploadId)
	}
	if fileName, ok := ctx.Get("FileName"); ok {
		str += fmt.Sprintf("[%s]", fileName)
	}
	ctx.SetPrefix(str + prefix)
}