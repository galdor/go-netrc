package netrc

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
)

type Entry struct {
	Machine  string
	Port     int
	Login    string
	Password string
	Account  string
}

type Entries []Entry

func DefaultPath() string {
	if path := os.Getenv("NETRC"); path != "" {
		return path
	}

	homePath, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return path.Join(homePath, ".netrc")
}

func (entries *Entries) Load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("cannot open %q: %w", filePath, err)
	}
	defer file.Close()

	var es Entries

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	lineNumber := 0

	for scanner.Scan() {
		lineNumber++

		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		var e Entry
		if err := e.Load(line); err != nil {
			return fmt.Errorf("line %d: cannot parse entry: %w",
				lineNumber, err)
		}

		es = append(es, e)
	}

	*entries = es

	return nil
}

func (e *Entry) Load(line string) error {
	data := []byte(line)

	skipSpace := func() {
		for len(data) > 0 {
			if data[0] != ' ' && data[0] != '\t' {
				break
			}

			data = data[1:]
		}
	}

	readToken := func() string {
		skipSpace()

		if len(data) == 0 {
			return ""
		}

		space := bytes.IndexAny(data, " \t")

		var tokenEnd int
		if space >= 0 {
			tokenEnd = space
		} else {
			tokenEnd = len(data)
		}

		token := data[:tokenEnd]
		data = data[tokenEnd:]

		return string(token)
	}

	for len(data) > 0 {
		token := readToken()

		missingValueErr := fmt.Errorf("missing value after token %q", token)

		switch token {
		case "machine":
			machine := readToken()
			if machine == "" {
				return missingValueErr
			}
			e.Machine = machine

		case "port":
			port := readToken()
			if port == "" {
				return missingValueErr
			}

			i64, err := strconv.ParseInt(port, 10, 64)
			if err != nil || i64 < 1 || i64 > 65535 {
				return fmt.Errorf("invalid port number %q", port)
			}

			e.Port = int(i64)

		case "login":
			login := readToken()
			if login == "" {
				return missingValueErr
			}
			e.Login = login

		case "password":
			password := readToken()
			if password == "" {
				return missingValueErr
			}
			e.Password = password

		case "account":
			account := readToken()
			if account == "" {
				return missingValueErr
			}
			e.Account = account

		default:
			return fmt.Errorf("invalid token %q", token)
		}
	}

	return nil
}
