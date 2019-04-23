package main

import (
	"flag"
	"fmt"
	"github.com/osrg/hookfs/pkg/example"
	"github.com/osrg/hookfs/pkg/hookfs"
	"math/rand"
	"os"
	"time"

	//hookfs "github.com/osrg/hookfs/hookfs"
	log "github.com/sirupsen/logrus"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "%s [OPTIONS] MOUNTPOINT ORIGINAL...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options\n")
		flag.PrintDefaults()
	}

	logLevel := flag.Int("log-level", 0, fmt.Sprintf("log level (%d..%d)", hookfs.LogLevelMin, hookfs.LogLevelMax))

	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(2)
	}

	mountpoint := flag.Arg(0)
	original := flag.Arg(1)
	hookfs.SetLogLevel(*logLevel)

	serve(original, mountpoint)
}

func serve(original string, mountpoint string) {
	fs, err := hookfs.NewHookFs(original, mountpoint, &example.MyHook{})
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("Serving %s", fs)
	log.Infof("Please run `fusermount -u %s` after using this, manually", mountpoint)
	if err = fs.Serve(); err != nil {
		log.Fatal(err)
	}
}
