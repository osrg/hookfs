package hookfs

import (
	"fmt"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type hookFile struct {
	file nodefs.File
	name string
	hook Hook
}

func newHookFile(file nodefs.File, name string, hook Hook) (*hookFile, error) {
	log.WithFields(log.Fields{
		"file": file,
		"name": name,
	}).Debug("Hooking a file")

	hookfile := &hookFile{
		file: file,
		name: name,
		hook: hook,
	}
	return hookfile, nil
}

// implements nodefs.File
func (h *hookFile) SetInode(inode *nodefs.Inode) {
	h.file.SetInode(inode)
}

// implements nodefs.File
func (h *hookFile) String() string {
	return fmt.Sprintf("HookFile{file=%s, name=%s}", h.file.String(), h.name)
}

// implements nodefs.File
func (h *hookFile) InnerFile() nodefs.File {
	return h.file.InnerFile()
}

// implements nodefs.File
func (h *hookFile) Read(dest []byte, off int64) (fuse.ReadResult, fuse.Status) {
	hook, hookEnabled := h.hook.(HookOnRead)
	var prehookBuf, posthookBuf []byte
	var prehookErr, posthookErr error
	var prehooked, posthooked bool
	var prehookCtx HookContext

	if hookEnabled {
		prehookBuf, prehookErr, prehooked, prehookCtx = hook.PreRead(h.name, int64(len(dest)), off)
		if prehooked {
			log.WithFields(log.Fields{
				"h": h,
				// "prehookBuf": prehookBuf,
				"prehookErr": prehookErr,
				"prehookCtx": prehookCtx,
			}).Debug("Read: Prehooked")
			return fuse.ReadResultData(prehookBuf), fuse.ToStatus(prehookErr)
		}
	}

	lowerRR, lowerCode := h.file.Read(dest, off)
	if hookEnabled {
		lowerRRBuf, lowerRRBufStatus := lowerRR.Bytes(make([]byte, lowerRR.Size()))
		if lowerRRBufStatus != fuse.OK {
			log.WithField("error", lowerRRBufStatus).Panic("lowerRR.Bytes() should not cause an error")
		}
		posthookBuf, posthookErr, posthooked = hook.PostRead(int32(lowerCode), lowerRRBuf, prehookCtx)
		if posthooked {
			if len(posthookBuf) != len(lowerRRBuf) {
				log.WithFields(log.Fields{
					"h": h,
					// "posthookBuf": posthookBuf,
					"posthookErr":    posthookErr,
					"posthookBufLen": len(posthookBuf),
					"lowerRRBufLen":  len(lowerRRBuf),
					"destLen":        len(dest),
				}).Warn("Read: Posthooked, but posthookBuf length != lowerrRRBuf length. You may get a strange behavior.")
			}

			log.WithFields(log.Fields{
				"h": h,
				// "posthookBuf": posthookBuf,
				"posthookErr": posthookErr,
			}).Debug("Read: Posthooked")
			return fuse.ReadResultData(posthookBuf), fuse.ToStatus(posthookErr)
		}
	}

	return lowerRR, lowerCode
}

// implements nodefs.File
func (h *hookFile) Write(data []byte, off int64) (uint32, fuse.Status) {
	hook, hookEnabled := h.hook.(HookOnWrite)
	var prehookErr, posthookErr error
	var prehooked, posthooked bool
	var prehookCtx HookContext

	if hookEnabled {
		prehookErr, prehooked, prehookCtx = hook.PreWrite(h.name, data, off)
		if prehooked {
			log.WithFields(log.Fields{
				"h":          h,
				"prehookErr": prehookErr,
				"prehookCtx": prehookCtx,
			}).Debug("Write: Prehooked")
			return 0, fuse.ToStatus(prehookErr)
		}
	}

	lowerWritten, lowerCode := h.file.Write(data, off)
	if hookEnabled {
		posthookErr, posthooked = hook.PostWrite(int32(lowerCode), prehookCtx)
		if posthooked {
			log.WithFields(log.Fields{
				"h":           h,
				"posthookErr": posthookErr,
			}).Debug("Write: Posthooked")
			return 0, fuse.ToStatus(posthookErr)
		}
	}

	return lowerWritten, lowerCode
}

// implements nodefs.File
func (h *hookFile) Flush() fuse.Status {
	return h.file.Flush()
}

// implements nodefs.File
func (h *hookFile) Release() {
	h.file.Release()
}

// implements nodefs.File
func (h *hookFile) Fsync(flags int) fuse.Status {
	hook, hookEnabled := h.hook.(HookOnFsync)
	var prehookErr, posthookErr error
	var prehooked, posthooked bool
	var prehookCtx HookContext

	if hookEnabled {
		prehookErr, prehooked, prehookCtx = hook.PreFsync(h.name, uint32(flags))
		if prehooked {
			log.WithFields(log.Fields{
				"h":          h,
				"prehookErr": prehookErr,
				"prehookCtx": prehookCtx,
			}).Debug("Fsync: Prehooked")
			return fuse.ToStatus(prehookErr)
		}
	}

	lowerCode := h.file.Fsync(flags)
	if hookEnabled {
		posthookErr, posthooked = hook.PostFsync(int32(lowerCode), prehookCtx)
		if posthooked {
			log.WithFields(log.Fields{
				"h":           h,
				"posthookErr": posthookErr,
			}).Debug("Fsync: Posthooked")
			return fuse.ToStatus(posthookErr)
		}
	}

	return lowerCode
}

// implements nodefs.File
func (h *hookFile) Truncate(size uint64) fuse.Status {
	return h.file.Truncate(size)
}

// implements nodefs.File
func (h *hookFile) GetAttr(out *fuse.Attr) fuse.Status {
	return h.file.GetAttr(out)
}

// implements nodefs.File
func (h *hookFile) Chown(uid uint32, gid uint32) fuse.Status {
	return h.file.Chown(uid, gid)
}

// implements nodefs.File
func (h *hookFile) Chmod(perms uint32) fuse.Status {
	return h.file.Chmod(perms)
}

// implements nodefs.File
func (h *hookFile) Utimens(atime *time.Time, mtime *time.Time) fuse.Status {
	return h.file.Utimens(atime, mtime)
}

// implements nodefs.File
func (h *hookFile) Allocate(off uint64, size uint64, mode uint32) fuse.Status {
	return h.file.Allocate(off, size, mode)
}

// implements nodefs.Flock
func (h *hookFile) Flock(flags int) fuse.Status {
	return h.file.Flock(flags)
}
