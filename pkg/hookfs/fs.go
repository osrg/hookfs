package hookfs

import (
	"fmt"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
	"github.com/hanwen/go-fuse/fuse/pathfs"
	log "github.com/sirupsen/logrus"
)

// HookFs is the object hooking the fs.
type HookFs struct {
	Original   string
	Mountpoint string
	FsName     string
	fs         pathfs.FileSystem
	hook       Hook
}

// NewHookFs creates a new HookFs object
func NewHookFs(original string, mountpoint string, hook Hook) (*HookFs, error) {
	log.WithFields(log.Fields{
		"original":   original,
		"mountpoint": mountpoint,
	}).Debug("Hooking a fs")

	loopbackfs := pathfs.NewLoopbackFileSystem(original)
	hookfs := &HookFs{
		Original:   original,
		Mountpoint: mountpoint,
		FsName:     "hookfs",
		fs:         loopbackfs,
		hook:       hook,
	}
	return hookfs, nil
}

// String implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) String() string {
	return fmt.Sprintf("HookFs{Original=%s, Mountpoint=%s, FsName=%s, Underlying fs=%s, hook=%s}",
		h.Original, h.Mountpoint, h.FsName, h.fs.String(), h.hook)
}

// SetDebug implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) SetDebug(debug bool) {
	h.fs.SetDebug(debug)
}

// GetAttr implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) GetAttr(name string, context *fuse.Context) (*fuse.Attr, fuse.Status) {
	return h.fs.GetAttr(name, context)
}

// Chmod implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Chmod(name string, mode uint32, context *fuse.Context) fuse.Status {
	return h.fs.Chmod(name, mode, context)
}

// Chown implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Chown(name string, uid uint32, gid uint32, context *fuse.Context) fuse.Status {
	return h.fs.Chown(name, uid, gid, context)
}

// Utimens implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Utimens(name string, Atime *time.Time, Mtime *time.Time, context *fuse.Context) fuse.Status {
	return h.fs.Utimens(name, Atime, Mtime, context)
}

// Truncate implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Truncate(name string, size uint64, context *fuse.Context) fuse.Status {
	return h.fs.Truncate(name, size, context)
}

// Access implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Access(name string, mode uint32, context *fuse.Context) fuse.Status {
	return h.fs.Access(name, mode, context)
}

// Link implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Link(oldName string, newName string, context *fuse.Context) fuse.Status {
	return h.fs.Link(oldName, newName, context)
}

// Mkdir implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Mkdir(name string, mode uint32, context *fuse.Context) fuse.Status {
	hook, hookEnabled := h.hook.(HookOnMkdir)
	var prehookErr, posthookErr error
	var prehooked, posthooked bool
	var prehookCtx HookContext

	if hookEnabled {
		prehooked, prehookCtx, prehookErr = hook.PreMkdir(name, mode)
		if prehooked {
			log.WithFields(log.Fields{
				"h":          h,
				"prehookErr": prehookErr,
				"prehookCtx": prehookCtx,
			}).Debug("Mkdir: Prehooked")
			if prehookErr == nil {
				log.WithFields(log.Fields{
					"h":          h,
					"prehookErr": prehookErr,
					"prehookCtx": prehookCtx,
				}).Fatal("Mkdir is prehooked, but did not returned an error. h is very strange.")
			}
			return fuse.ToStatus(prehookErr)
		}
	}

	lowerCode := h.fs.Mkdir(name, mode, context)
	if hookEnabled {
		posthooked, posthookErr = hook.PostMkdir(int32(lowerCode), prehookCtx)
		if posthooked {
			log.WithFields(log.Fields{
				"h":           h,
				"posthookErr": posthookErr,
			}).Debug("Mkdir: Posthooked")
			return fuse.ToStatus(posthookErr)
		}
	}

	return lowerCode
}

// Mknod implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Mknod(name string, mode uint32, dev uint32, context *fuse.Context) fuse.Status {
	return h.fs.Mknod(name, mode, dev, context)
}

// Rename implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Rename(oldName string, newName string, context *fuse.Context) fuse.Status {
	hook, hookEnabled := h.hook.(HookOnRename)
	if hookEnabled {
		preHooked, prehookCtx, prehookErr := hook.PreRename(oldName, newName)
		if preHooked {
			log.WithFields(log.Fields{
				"h":          h,
				"prehookErr": prehookErr,
				"prehookCtx": prehookCtx,
			}).Debug("Rename: Prehooked")

			if prehookErr != nil {
				return fuse.ToStatus(prehookErr)
			}
		}
	}

	status := h.fs.Rename(oldName, newName, context)

	if hookEnabled {
		postHooked, posthookCtx, posthookErr := hook.PostRename(oldName, newName)
		if postHooked {
			log.WithFields(log.Fields{
				"h":          h,
				"posthookErr": postHooked,
				"posthookCtx": posthookCtx,
			}).Debug("Rename: Posthooked")

			if posthookErr != nil {
				return fuse.ToStatus(posthookErr)
			}
		}
	}

	return status
}

// Rmdir implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Rmdir(name string, context *fuse.Context) fuse.Status {
	hook, hookEnabled := h.hook.(HookOnRmdir)
	var prehookErr, posthookErr error
	var prehooked, posthooked bool
	var prehookCtx HookContext

	if hookEnabled {
		prehooked, prehookCtx, prehookErr = hook.PreRmdir(name)
		if prehooked {
			log.WithFields(log.Fields{
				"h":          h,
				"prehookErr": prehookErr,
				"prehookCtx": prehookCtx,
			}).Debug("Rmdir: Prehooked")
			if prehookErr == nil {
				log.WithFields(log.Fields{
					"h":          h,
					"prehookErr": prehookErr,
					"prehookCtx": prehookCtx,
				}).Fatal("Rmdir is prehooked, but did not returned an error. h is very strange.")
			}
			return fuse.ToStatus(prehookErr)
		}
	}

	lowerCode := h.fs.Rmdir(name, context)
	if hookEnabled {
		posthooked, posthookErr = hook.PostRmdir(int32(lowerCode), prehookCtx)
		if posthooked {
			log.WithFields(log.Fields{
				"h":           h,
				"posthookErr": posthookErr,
			}).Debug("Mkdir: Posthooked")
			return fuse.ToStatus(posthookErr)
		}
	}

	return lowerCode
}

// Unlink implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Unlink(name string, context *fuse.Context) fuse.Status {
	return h.fs.Unlink(name, context)
}

// GetXAttr implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) GetXAttr(name string, attribute string, context *fuse.Context) ([]byte, fuse.Status) {
	return h.fs.GetXAttr(name, attribute, context)
}

// ListXAttr implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) ListXAttr(name string, context *fuse.Context) ([]string, fuse.Status) {
	return h.fs.ListXAttr(name, context)
}

// RemoveXAttr implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) RemoveXAttr(name string, attr string, context *fuse.Context) fuse.Status {
	return h.fs.RemoveXAttr(name, attr, context)
}

// SetXAttr implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) SetXAttr(name string, attr string, data []byte, flags int, context *fuse.Context) fuse.Status {
	return h.fs.SetXAttr(name, attr, data, flags, context)
}

// OnMount implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) OnMount(nodeFs *pathfs.PathNodeFs) {
	h.fs.OnMount(nodeFs)
	hook, hookEnabled := h.hook.(HookWithInit)
	if hookEnabled {
		err := hook.Init()
		if err != nil {
			log.Error(err)
			log.Warn("Disabling hook")
			h.hook = nil
		}
	}
}

// OnUnmount implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) OnUnmount() {
	h.fs.OnUnmount()
}

// Open implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Open(name string, flags uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	hook, hookEnabled := h.hook.(HookOnOpen)
	var prehookErr, posthookErr error
	var prehooked, posthooked bool
	var prehookCtx HookContext

	if hookEnabled {
		prehooked, prehookCtx, prehookErr = hook.PreOpen(name, flags)
		if prehooked {
			log.WithFields(log.Fields{
				"h":          h,
				"prehookErr": prehookErr,
				"prehookCtx": prehookCtx,
			}).Debug("Open: Prehooked")
			if prehookErr == nil {
				log.WithFields(log.Fields{
					"h":          h,
					"prehookErr": prehookErr,
					"prehookCtx": prehookCtx,
				}).Fatal("Open is prehooked, but did not returned an error. h is very strange.")
			}
			return nil, fuse.ToStatus(prehookErr)
		}
	}

	lowerFile, lowerCode := h.fs.Open(name, flags, context)
	hFile, hErr := newHookFile(lowerFile, name, h.hook)
	if hErr != nil {
		log.WithField("error", hErr).Panic("NewHookFile() should not cause an error")
	}

	if hookEnabled {
		posthooked, posthookErr = hook.PostOpen(int32(lowerCode), prehookCtx)
		if posthooked {
			log.WithFields(log.Fields{
				"h":           h,
				"posthookErr": posthookErr,
			}).Debug("Open: Posthooked")
			return hFile, fuse.ToStatus(posthookErr)
		}
	}

	return hFile, lowerCode
}

// Create implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Create(name string, flags uint32, mode uint32, context *fuse.Context) (nodefs.File, fuse.Status) {
	lowerFile, lowerCode := h.fs.Create(name, flags, mode, context)
	hFile, hErr := newHookFile(lowerFile, name, h.hook)
	if hErr != nil {
		log.WithField("error", hErr).Panic("NewHookFile() should not cause an error")
	}
	return hFile, lowerCode
}

// OpenDir implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) OpenDir(name string, context *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	hook, hookEnabled := h.hook.(HookOnOpenDir)
	var prehookErr, posthookErr error
	var prehooked, posthooked bool
	var prehookCtx HookContext

	if hookEnabled {
		prehooked, prehookCtx, prehookErr = hook.PreOpenDir(name)
		if prehooked {
			log.WithFields(log.Fields{
				"h":          h,
				"prehookErr": prehookErr,
				"prehookCtx": prehookCtx,
			}).Debug("OpenDir: Prehooked")
			if prehookErr == nil {
				log.WithFields(log.Fields{
					"h":          h,
					"prehookErr": prehookErr,
					"prehookCtx": prehookCtx,
				}).Fatal("OpenDir is prehooked, but did not returned an error. h is very strange.")
			}
			return nil, fuse.ToStatus(prehookErr)
		}
	}

	lowerEnts, lowerCode := h.fs.OpenDir(name, context)
	if hookEnabled {
		posthooked, posthookErr = hook.PostOpenDir(int32(lowerCode), prehookCtx)
		if posthooked {
			log.WithFields(log.Fields{
				"h":           h,
				"posthookErr": posthookErr,
			}).Debug("OpenDir: Posthooked")
			return lowerEnts, fuse.ToStatus(posthookErr)
		}
	}

	return lowerEnts, lowerCode
}

// Symlink implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Symlink(value string, linkName string, context *fuse.Context) fuse.Status {
	return h.fs.Symlink(value, linkName, context)
}

// Readlink implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) Readlink(name string, context *fuse.Context) (string, fuse.Status) {
	return h.fs.Readlink(name, context)
}

// StatFs implements hanwen/go-fuse/fuse/pathfs.FileSystem. You are not expected to call h manually.
func (h *HookFs) StatFs(name string) *fuse.StatfsOut {
	return h.fs.StatFs(name)
}

// Serve starts the server (blocking).
func (h *HookFs) Serve() error {
	server, err := newHookServer(h)
	if err != nil {
		return err
	}
	server.Serve()
	return nil
}

// Serve initiates the FUSE loop. Normally, callers should run Serve()
// and wait for it to exit, but tests will want to run this in a
// goroutine.
func (h *HookFs) ServeAsync() (*fuse.Server, error) {
	server, err := newHookServer(h)
	if err != nil {
		return nil, err
	}
	go func() {
		server.Serve()
	}()

	return server, nil
}
