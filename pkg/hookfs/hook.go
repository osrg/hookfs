package hookfs

// Hook is the base interface for user-written hooks.
//
// You have to implement HookXXX (e.g. HookOnOpen, HookOnRead, HookOnWrite, ..) interfaces.
type Hook interface{}

// HookContext is the context objects for interaction between prehooks and posthooks.
type HookContext interface{}

// HookWithInit is called on mount. This also implements Hook.
type HookWithInit interface {
	Init() (err error)
}

// HookOnOpen is called on open. This also implements Hook.
type HookOnOpen interface {
	// if hooked is true, the real open() would not be called
	PreOpen(path string, flags uint32) (hooked bool, ctx HookContext, err error)
	PostOpen(realRetCode int32, prehookCtx HookContext) (hooked bool, err error)
}

// HookOnRead is called on read. This also implements Hook.
type HookOnRead interface {
	// if hooked is true, the real read() would not be called
	PreRead(path string, length int64, offset int64) (buf []byte, hooked bool, ctx HookContext, err error)
	PostRead(realRetCode int32, realBuf []byte, prehookCtx HookContext) (buf []byte, hooked bool, err error)
}

// HookOnWrite is called on write. This also implements Hook.
type HookOnWrite interface {
	// if hooked is true, the real write() would not be called
	PreWrite(path string, buf []byte, offset int64) (hooked bool, ctx HookContext, err error)
	PostWrite(realRetCode int32, prehookCtx HookContext) (hooked bool, err error)
}

// HookOnMkdir is called on mkdir. This also implements Hook.
type HookOnMkdir interface {
	// if hooked is true, the real mkdir() would not be called
	PreMkdir(path string, mode uint32) (hooked bool, ctx HookContext, err error)
	PostMkdir(realRetCode int32, prehookCtx HookContext) (hooked bool, err error)
}

// HookOnRmdir is called on rmdir. This also implements Hook.
type HookOnRmdir interface {
	// if hooked is true, the real rmdir() would not be called
	PreRmdir(path string) (hooked bool, ctx HookContext, err error)
	PostRmdir(realRetCode int32, prehookCtx HookContext) (hooked bool, err error)
}

// HookOnOpenDir is called on opendir. This also implements Hook.
type HookOnOpenDir interface {
	// if hooked is true, the real opendir() would not be called
	PreOpenDir(path string) (hooked bool, ctx HookContext, err error)
	PostOpenDir(realRetCode int32, prehookCtx HookContext) (hooked bool, err error)
}

// HookOnFsync is called on fsync. This also implements Hook.
type HookOnFsync interface {
	// if hooked is true, the real fsync() would not be called
	PreFsync(path string, flags uint32) (hooked bool, ctx HookContext, err error)
	PostFsync(realRetCode int32, prehookCtx HookContext) (hooked bool, err error)
}

// HookOnRename is called on rename
type HookOnRename interface {
	PreRename(oldPatgh string, newPath string) (hooked bool, ctx HookContext, err error)
	PostRename(oldPatgh string, newPath string) (hooked bool, ctx HookContext, err error)
}