package main

import (
	log "github.com/Sirupsen/logrus"
	hookfs "github.com/osrg/hookfs/hookfs"
	"syscall"
	"time"
)

// implements hookfs.HookContext
type MyHookContext struct {
	path string
}

// implements hookfs.Hook
type MyHook struct{}

// implements hookfs.HookWithInit
func (this *MyHook) Init() error {
	log.WithFields(log.Fields{
		"this": this,
	}).Info("MyInit: initializing")
	return nil
}

// implements hookfs.HookOnOpen
func (this *MyHook) PreOpen(path string, flags uint32) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(5) {
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
		}).Info("MyPreOpen: returning EIO")
		return syscall.EIO, true, ctx
	}
	return nil, false, ctx
}

// implements hookfs.HookOnOpen
func (this *MyHook) PostOpen(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(5) {
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
		}).Info("MyPostOpen: returning EPERM")
		return syscall.EPERM, true
	}
	return nil, false
}

// implements hookfs.HookOnRead
func (this *MyHook) PreRead(path string, length int64, offset int64) ([]byte, error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(3) {
		sleep := 3 * time.Second
		log.WithFields(log.Fields{
			"this":  this,
			"ctx":   ctx,
			"sleep": sleep,
		}).Info("MyPreRead: sleeping")
		time.Sleep(sleep)
	}
	return nil, nil, false, ctx
}

// implements hookfs.HookOnRead
func (this *MyHook) PostRead(realRetCode int32, realBuf []byte, ctx hookfs.HookContext) ([]byte, error, bool) {
	if probab(70) {
		buf := []byte("Hello HookFS hooked Data!\n")
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
			"buf":  buf,
		}).Info("MyPostRead: returning injected buffer")
		return buf, nil, true
	}
	return nil, nil, false
}
