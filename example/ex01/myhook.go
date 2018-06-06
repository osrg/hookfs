package main

import (
	"syscall"
	"time"

	log "github.com/Sirupsen/logrus"
	hookfs "github.com/osrg/hookfs/hookfs"
)

// MyHookContext implements hookfs.HookContext
type MyHookContext struct {
	path string
}

// MyHook implements hookfs.Hook
type MyHook struct{}

// Init implements hookfs.HookWithInit
func (h *MyHook) Init() error {
	log.WithFields(log.Fields{
		"h": h,
	}).Info("MyInit: initializing")
	return nil
}

// PreOpen implements hookfs.HookOnOpen
func (h *MyHook) PreOpen(path string, flags uint32) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(5) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreOpen: returning EIO")
		return syscall.EIO, true, ctx
	}
	return nil, false, ctx
}

// PostOpen implements hookfs.HookOnOpen
func (h *MyHook) PostOpen(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(5) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostOpen: returning EPERM")
		return syscall.EPERM, true
	}
	return nil, false
}

// PreRead implements hookfs.HookOnRead
func (h *MyHook) PreRead(path string, length int64, offset int64) ([]byte, error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(3) {
		sleep := 3 * time.Second
		log.WithFields(log.Fields{
			"h":     h,
			"ctx":   ctx,
			"sleep": sleep,
		}).Info("MyPreRead: sleeping")
		time.Sleep(sleep)
	}
	return nil, nil, false, ctx
}

// PostRead implements hookfs.HookOnRead
func (h *MyHook) PostRead(realRetCode int32, realBuf []byte, ctx hookfs.HookContext) ([]byte, error, bool) {
	if probab(70) {
		buf := []byte("Hello HookFS hooked Data!\n")
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
			"buf": buf,
		}).Info("MyPostRead: returning injected buffer")
		return buf, nil, true
	}
	return nil, nil, false
}

// PreWrite implements hookfs.HookOnWrite
func (h *MyHook) PreWrite(path string, buf []byte, offset int64) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(3) {
		sleep := 3 * time.Second
		log.WithFields(log.Fields{
			"h":     h,
			"ctx":   ctx,
			"sleep": sleep,
		}).Info("MyPreWrite: sleeping")
		time.Sleep(sleep)
	}
	return nil, false, ctx
}

// PostWrite implements hookfs.HookOnWrite
func (h *MyHook) PostWrite(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(70) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostWrite: returning ENOSPC")
		return syscall.ENOSPC, true
	}
	return nil, false
}

// PreMkdir implements hookfs.HookOnMkdir
func (h *MyHook) PreMkdir(path string, mode uint32) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(95) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreMkdir: returning EACCES")
		return syscall.EACCES, true, ctx
	}
	return nil, false, ctx
}

// PostMkdir implements hookfs.HookOnMkdir
func (h *MyHook) PostMkdir(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(5) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostMkdir: returning EPERM")
		return syscall.EPERM, true
	}
	return nil, false
}

// PreRmdir implements hookfs.HookOnRmdir
func (h *MyHook) PreRmdir(path string) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(30) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreRmdir: returning EACCES")
		return syscall.EACCES, true, ctx
	}
	return nil, false, ctx
}

// PostRmdir implements hookfs.HookOnRmdir
func (h *MyHook) PostRmdir(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(30) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostRmdir: returning EPERM")
		return syscall.EPERM, true
	}
	return nil, false
}

// PreOpenDir implements hookfs.HookOnOpenDir
func (h *MyHook) PreOpenDir(path string) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(30) && path != "" {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreOpenDir: returning EACCES")
		return syscall.EACCES, true, ctx
	}
	return nil, false, ctx
}

// PostOpenDir implements hookfs.HookOnOpenDir
func (h *MyHook) PostOpenDir(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(30) && ctx.(MyHookContext).path != "" {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostOpenDir: returning EPERM")
		return syscall.EPERM, true
	}
	return nil, false
}

// PreFsync implements hookfs.HookOnFsync
func (h *MyHook) PreFsync(path string, flags uint32) (error, bool, hookfs.HookContext) {
	ctx := MyHookContext{path: path}
	if probab(90) && path != "" {
		sleep := 3 * time.Second
		log.WithFields(log.Fields{
			"h":     h,
			"ctx":   ctx,
			"sleep": sleep,
		}).Info("MyPreFsync: sleeping")
		time.Sleep(sleep)
	}
	return nil, false, ctx
}

// PostFsync implements hookfs.HookOnFsync
func (h *MyHook) PostFsync(realRetCode int32, ctx hookfs.HookContext) (error, bool) {
	if probab(80) && ctx.(MyHookContext).path != "" {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostFsync: returning EIO")
		return syscall.EIO, true
	}
	return nil, false
}
