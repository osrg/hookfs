package hookfs

import (
	"fmt"
	"github.com/hanwen/go-fuse/fuse"
	"github.com/osrg/hookfs/pkg/utils"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func TestRenameHook_RenameHookRenameTest(t *testing.T) {
	//init action
	original := filepath.Join(string(filepath.Separator), "tmp", fmt.Sprintf("dev-%d", time.Now().Unix()))
	mountpoint := filepath.Join(string(filepath.Separator), "tmp", fmt.Sprintf("mountpoint-%d", time.Now().Unix()))
	newFuseServer(t, original, mountpoint)
	//remember to call unmount after you do not use it
	defer cleanUp(mountpoint, original)

	//normal logic
	log.Print(filepath.Join(mountpoint, "tsdb.txt"))
	_, err := os.Create(filepath.Join(mountpoint, "tsdb.txt"))

	if err != nil {
		log.Printf("%v",err)
	}

	//rename should be failed
	err = os.Rename(filepath.Join(mountpoint, "tsdb.txt"), filepath.Join(mountpoint, "tsdbNew.txt"))
	utils.NotOk(t, err)
	fmt.Println(err)
}

func newFuseServer(t *testing.T, original,mountpoint string)(*fuse.Server){
	createDirIfAbsent(original)
	createDirIfAbsent(mountpoint)
	fs, err :=  NewHookFs(original, mountpoint, &TestRenameHook{})
	utils.Ok(t, err)
	server, err := fs.ServeAsync()
	if err != nil {
		log.Fatalf("start server failed, %v", err)
	}
	utils.Ok(t, err)

	return server
}

func cleanUp(mountpoint string, original string) {
	syscall.Unmount(mountpoint, -1)

	os.RemoveAll(mountpoint)
	os.RemoveAll(original)
	fmt.Println("Done")
}

func createDirIfAbsent(name string) {
	_, err := os.Stat(name)
	if err != nil {
		os.Mkdir(name, os.ModePerm)
	}
}

type TestRenameHook struct{}

func (h *TestRenameHook) PreRename(oldPatgh string, newPath string) (hooked bool, ctx HookContext, err error) {
	fmt.Printf("Pre renamed file from %s to %s \n", oldPatgh, newPath)
	return true, nil,  syscall.EIO
}

func (h *TestRenameHook) PostRename(oldPatgh string, newPath string) (hooked bool, ctx HookContext, err error) {
	fmt.Printf("Post renamed file from %s to %s \n", oldPatgh, newPath)
	return false, nil, nil
}
