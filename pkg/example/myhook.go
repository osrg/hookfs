package example

import (
	"math/rand"
	"syscall"
	"time"

	hookfs "github.com/qiffang/hookfs/pkg/hookfs"
	log "github.com/sirupsen/logrus"
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
func (h *MyHook) PreOpen(path string, flags uint32) (bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	if probab(5) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreOpen: returning EIO")
		return true, ctx, syscall.EIO
	}
	return false, ctx, nil
}

// PostOpen implements hookfs.HookOnOpen
func (h *MyHook) PostOpen(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	if probab(5) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostOpen: returning EPERM")
		return true, syscall.EPERM
	}
	return false, nil
}

// PreRead implements hookfs.HookOnRead
func (h *MyHook) PreRead(path string, length int64, offset int64) ([]byte, bool, hookfs.HookContext, error) {
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
	return nil, false, ctx, nil
}

// PostRead implements hookfs.HookOnRead
func (h *MyHook) PostRead(realRetCode int32, realBuf []byte, ctx hookfs.HookContext) ([]byte, bool, error) {
	if probab(70) {
		buf := []byte("Hello HookFS hooked Data!\n")
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
			"buf": buf,
		}).Info("MyPostRead: returning injected buffer")
		return buf, true, nil
	}
	return nil, false, nil
}

// PreWrite implements hookfs.HookOnWrite
func (h *MyHook) PreWrite(path string, buf []byte, offset int64) (bool, hookfs.HookContext, error) {
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
	return false, ctx, nil
}

// PostWrite implements hookfs.HookOnWrite
func (h *MyHook) PostWrite(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	if probab(70) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostWrite: returning ENOSPC")
		return true, syscall.ENOSPC
	}
	return false, nil
}

// PreMkdir implements hookfs.HookOnMkdir
func (h *MyHook) PreMkdir(path string, mode uint32) (bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	if probab(95) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreMkdir: returning EACCES")
		return true, ctx, syscall.EACCES
	}
	return false, ctx, nil
}

// PostMkdir implements hookfs.HookOnMkdir
func (h *MyHook) PostMkdir(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	if probab(5) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostMkdir: returning EPERM")
		return true, syscall.EPERM
	}
	return false, nil
}

// PreRmdir implements hookfs.HookOnRmdir
func (h *MyHook) PreRmdir(path string) (bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	if probab(30) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreRmdir: returning EACCES")
		return true, ctx, syscall.EACCES
	}
	return false, ctx, nil
}

// PostRmdir implements hookfs.HookOnRmdir
func (h *MyHook) PostRmdir(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	if probab(30) {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostRmdir: returning EPERM")
		return true, syscall.EPERM
	}
	return false, nil
}

// PreOpenDir implements hookfs.HookOnOpenDir
func (h *MyHook) PreOpenDir(path string) (bool, hookfs.HookContext, error) {
	ctx := MyHookContext{path: path}
	if probab(30) && path != "" {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPreOpenDir: returning EACCES")
		return true, ctx, syscall.EACCES
	}
	return false, ctx, nil
}

// PostOpenDir implements hookfs.HookOnOpenDir
func (h *MyHook) PostOpenDir(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	if probab(30) && ctx.(MyHookContext).path != "" {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostOpenDir: returning EPERM")
		return true, syscall.EPERM
	}
	return false, nil
}

// PreFsync implements hookfs.HookOnFsync
func (h *MyHook) PreFsync(path string, flags uint32) (bool, hookfs.HookContext, error) {
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
	return false, ctx, nil
}

// PostFsync implements hookfs.HookOnFsync
func (h *MyHook) PostFsync(realRetCode int32, ctx hookfs.HookContext) (bool, error) {
	if probab(80) && ctx.(MyHookContext).path != "" {
		log.WithFields(log.Fields{
			"h":   h,
			"ctx": ctx,
		}).Info("MyPostFsync: returning EIO")
		return true, syscall.EIO
	}
	return false, nil
}

func probab(percentage int) bool {
	return rand.Intn(99) < percentage
}
