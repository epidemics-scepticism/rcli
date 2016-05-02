package main

import (
	"fmt"
	"os"
	"strings"

	"./ricochet"
	"golang.org/x/crypto/ssh/terminal"
)

var (
	orig *terminal.State = nil
	term *terminal.Terminal = nil
	promptEnd string = "> "
	curHost string = ""
	inputMap = map[string]func(args string) {
		"add": cmdAdd,
		"remove": cmdRemove,
		"list": cmdList,
		"blacklist": cmdBlacklist,
		"msg": cmdMsg,
		"help": cmdHelp,
	}
)

func split(s, d string) (string,string) {
	r := strings.SplitN(s, d, 2)
	if len(r) == 2 {
		return r[0], r[1]
	} else if len(r) == 1 {
		return r[0], ""
	} else {
		return "", ""
	}
}

func inputHandler(arg string) {
	if len(arg) > 1 && arg[0] == '/' {
		cmd, args := split(arg[1:], " ")
		if f, ok := inputMap[cmd]; ok {
			f(args)
		} else {
			termPrint("unknown command, try /help")
		}
	} else {
		cmdMsg(curHost + " " + arg)
	}
	termUpdate()
}

func cmdAdd(args string) {
	host, msg := split(args, " ")
	req := &ricochet.ContactRequest{
		host,
		*nick,
		msg,
	}
	if e := r.e.AddContact(host, req); e != nil {
		s := fmt.Sprintf("%v", e)
		termPrint(s)
		return
	}
	r.Contacts = append(r.Contacts, host)
}

func cmdRemove(args string) {
	host, _ := split(args, " ")
	if e := r.e.RemoveContact(host); e != nil {
		s := fmt.Sprintf("%v", e)
		termPrint(s)
		return
	}
	for k := range r.Contacts { // wtf r u doin?
		if r.Contacts[k] == host {
			r.Contacts = append(r.Contacts[:k-1], r.Contacts[k:]...)
			return
		}
	}
}

func cmdList(args string) {
	for _, v := range r.Contacts {
		termPrint(v)
	}
}

func cmdBlacklist(args string) {
	host, opt := split(args, " ")
	opt = strings.ToLower(opt)
	if strings.Contains(opt, "true") || strings.Contains(opt, "yes") {
		if e := r.e.BlacklistContact(host, true); e != nil {
			s := fmt.Sprintf("%v", e)
			termPrint(s)
			return
		}
	} else if strings.Contains(opt, "false") || strings.Contains(opt, "no") {
		if e := r.e.BlacklistContact(host, false); e != nil {
			s := fmt.Sprintf("%v", e)
			termPrint(s)
			return
		}
	} else {
		termPrint("invalid option, try /help")
	}
}

func cmdMsg(args string) {
	host, msg := split(args, " ")
	if curHost != host {
		curHost = host
	}
	if e := r.e.SendMsg(host, msg); e != nil {
		s := fmt.Sprintf("%v", e)
		termPrint(s)
		return
	}
}

func cmdHelp(args string) {
	termPrint("/add ricochet:dogdogdogdogdogd.onion hello, is this dog?")
	termPrint("/remove ricochet:notdognotdognotd.onion")
	termPrint("/list")
	termPrint("/msg ricochet:dogdogdogdogdogd.onion hello, this is dog")
	termPrint("/blacklist ricochet:abusiveassholeab.onion true")
	termPrint("/help")
}

func termPrint(line string) {
	if term != nil {
		term.Write([]byte(line + "\n"))
		termUpdate()
	} else {
		fmt.Fprintln(os.Stderr, line)
	}
}

func termInit() {
	var e error
	term = terminal.NewTerminal(os.Stdin, promptEnd)
	orig, e = terminal.MakeRaw(0)
	if e != nil {
		s := fmt.Sprintf("%v", e)
		termPrint(s)
	}
}

func termEnd() {
	terminal.Restore(0, orig)
}

func termUpdate() {
	if term == nil {
		return
	}
	w, h, e := terminal.GetSize(0)
	if e != nil {
		s := fmt.Sprintf("%v", e)
		termPrint(s)
		return
	}
	term.SetPrompt(curHost + promptEnd)
	term.SetSize(w, h)
}
