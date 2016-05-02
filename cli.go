package main

import (
	"fmt"
	"os"
//	"strings"
	"flag"
//	"crypto/rsa"
//	"crypto/rand"
//	"encoding/json"

//	"github.com/yawning/bulb"
//	"github.com/yawning/ricochet"
//	"golang.org/x/crypto/ssh/terminal"
)

var (
	cPort = flag.String("cport", os.Getenv("TOR_CONTROL_PORT"), "Tor Control Port (Defaults to 9151 if no envvar)")
	cHost = flag.String("chost", os.Getenv("TOR_CONTROL_HOST"), "Tor Control Host (Defaults to 127.0.0.1 if no envvar)")
	nick = flag.String("nick", os.Getenv("USER"), "Nickname to announce to peers")
	r *RCLIConfig
)

func getArgs() {
	flag.Parse()
	if len(*cPort) == 0 {
		*cPort = "9151"
	}
	if len(*cHost) == 0 {
		*cHost = "127.0.0.1"
	}
}

func main() {
	getArgs()
	r = &RCLIConfig{}
	end, e := r.Load()
	if e != nil {
		s := fmt.Sprintf("%v", e)
		termPrint(s)
		return
	}
	defer r.Save()
	go func() {
		for evt := range end.EventChan {
			s := fmt.Sprintf("%v", evt)
			termPrint(s)
		}
	}()
	termInit()
	defer termEnd()
	for {
		if l, e := term.ReadLine(); e != nil {
			s := fmt.Sprintf("%v", e)
			termPrint(s)
			return
		} else {
			inputHandler(l)
		}
	}
}
