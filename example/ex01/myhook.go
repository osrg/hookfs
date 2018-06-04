package main

import (
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	hookfs "github.com/osrg/hookfs/hookfs"
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

// implements hookfs.HookOnWrite
func (this *MyHook) PreWrite(path string, buf []byte, offset int64) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(3) {
		sleep := 3 * time.Second
		log.WithFields(log.Fields{
			"this":  this,
			"ctx":   ctx,
			"sleep": sleep,
		}).Info("MyPreWrite: sleeping")
		time.Sleep(sleep)
	}
	return nil, false, ctx
}

// implements hookfs.HookOnWrite
func (this *MyHook) PostWrite(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(70) {
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
		}).Info("MyPostWrite: returning ENOSPC")
		return syscall.ENOSPC, true
	}
	return nil, false
}

// implements hookfs.HookOnMkdir
func (this *MyHook) PreMkdir(path string, mode uint32) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(95) {
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
		}).Info("MyPreMkdir: returning EACCES")
		return syscall.EACCES, true, ctx
	}
	return nil, false, ctx
}

// implements hookfs.HookOnMkdir
func (this *MyHook) PostMkdir(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(5) {
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
		}).Info("MyPostMkdir: returning EPERM")
		return syscall.EPERM, true
	}
	return nil, false
}

// implements hookfs.HookOnRmdir
func (this *MyHook) PreRmdir(path string) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(30) {
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
		}).Info("MyPreRmdir: returning EACCES")
		return syscall.EACCES, true, ctx
	}
	return nil, false, ctx
}

// implements hookfs.HookOnRmdir
func (this *MyHook) PostRmdir(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(30) {
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
		}).Info("MyPostRmdir: returning EPERM")
		return syscall.EPERM, true
	}
	return nil, false
}

// implements hookfs.HookOnOpenDir
func (this *MyHook) PreOpenDir(path string) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(30) && path != "" {
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
		}).Info("MyPreOpenDir: returning EACCES")
		return syscall.EACCES, true, ctx
	}
	return nil, false, ctx
}

// implements hookfs.HookOnOpenDir
func (this *MyHook) PostOpenDir(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(30) && ctx.(MyHookContext).path != "" {
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
		}).Info("MyPostOpenDir: returning EPERM")
		return syscall.EPERM, true
	}
	return nil, false
}

// implements hookfs.HookOnFsync
func (this *MyHook) PreFsync(path string, flags uint32) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(90) && path != "" {
		sleep := 3 * time.Second
		log.WithFields(log.Fields{
			"this":  this,
			"ctx":   ctx,
			"sleep": sleep,
		}).Info("MyPreFsync: sleeping")
		time.Sleep(sleep)
	}
	return nil, false, ctx
}

// implements hookfs.HookOnFsync
func (this *MyHook) PostFsync(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(80) && ctx.(MyHookContext).path != "" {
		log.WithFields(log.Fields{
			"this": this,
			"ctx":  ctx,
		}).Info("MyPostFsync: returning EIO")
		return syscall.EIO, true
	}
	return nil, false
}
