package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/exograd/go-program"
	"github.com/galdor/go-netrc"
)

func main() {
	p := program.NewProgram("example",
		"an example program for the go-netrc library")

	p.AddCommand("list", "list netrc entries", cmdList)

	p.ParseCommandLine()
	p.Run()
}

func cmdList(p *program.Program) {
	var entries netrc.Entries
	if err := entries.Load(netrc.DefaultPath()); err != nil {
		p.Fatal("cannot load netrc entries: %v", err)
	}

	keys := make([]string, len(entries))
	keyWidth := 0

	for i, e := range entries {
		var buf bytes.Buffer

		if e.Login != "" {
			buf.WriteString(e.Login)
			buf.WriteByte('@')
		}

		buf.WriteString(e.Machine)

		if e.Port != 0 {
			buf.WriteByte(':')
			buf.WriteString(strconv.Itoa(e.Port))
		}

		keys[i] = buf.String()

		if len(keys[i]) > keyWidth {
			keyWidth = len(keys[i])
		}
	}

	for i, e := range entries {
		fmt.Printf("%-*s  %s\n", keyWidth, keys[i], e.Password)
	}
}
