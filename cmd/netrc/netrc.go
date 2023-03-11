package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/exograd/go-program"
	"github.com/galdor/go-netrc"
)

func main() {
	var c *program.Command

	p := program.NewProgram("example",
		"an example program for the go-netrc library")

	p.AddCommand("list", "list all netrc entries", cmdList)

	c = p.AddCommand("search", "search for netrc entries", cmdSearch)
	c.AddOption("m", "machine", "hostname", "", "a hostname to match")
	c.AddOption("p", "port", "number", "", "a port number to match")
	c.AddOption("l", "login", "login", "", "a login to match")
	c.AddOption("a", "account", "name", "", "an account name to match")

	p.ParseCommandLine()
	p.Run()
}

func cmdList(p *program.Program) {
	var entries netrc.Entries
	if err := entries.Load(netrc.DefaultPath()); err != nil {
		p.Fatal("cannot load netrc entries: %v", err)
	}

	printEntries(os.Stdout, entries)
}

func cmdSearch(p *program.Program) {
	search := netrc.Search{
		Machine: p.OptionValue("machine"),
		Login:   p.OptionValue("login"),
		Account: p.OptionValue("account"),
	}

	if p.IsOptionSet("port") {
		portString := p.OptionValue("port")
		i64, err := strconv.ParseInt(portString, 10, 64)
		if err != nil || i64 < 1 || i64 > 65535 {
			p.Fatal("invalid port number %q", portString)
		}

		search.Port = int(i64)
	}

	var entries netrc.Entries
	if err := entries.Load(netrc.DefaultPath()); err != nil {
		p.Fatal("cannot load netrc entries: %v", err)
	}

	printEntries(os.Stdout, entries.Search(search))
}

func printEntries(w io.Writer, entries netrc.Entries) {
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
		fmt.Fprintf(w, "%-*s  %s\n", keyWidth, keys[i], e.Password)
	}
}
