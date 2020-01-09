package main

import (
	"fmt"
	"github.com/tomseago/go-eltee/api"
	"os"
	"os/user"
	"sort"
	"strconv"
	"unicode"

	"github.com/eyethereal/uniline"
	"github.com/mgutz/ansi"
)

type localContext struct {
	c         *api.AutoClient
	stateName string
    cpListener *cpListening
}

type helpEntry struct {
	syntax string
	short  string
	man    string
}

var errColor = ansi.ColorFunc("red+h")
var successColor = ansi.ColorFunc("blue")

var aRed = ansi.ColorCode("red+h")
var aReset = ansi.ColorCode("reset")

var commands = make(map[string]func(*localContext, []string))
var help = make(map[string]*helpEntry)

func NewLocalContext() *localContext {

	return &localContext{
		c: api.NewAutoClient(":3434"),
	}
}

func (lc *localContext) Run() {

	// It doesn't work right now for the prompt to have ansi color codes in it. Will have to work on that I guess.
	prompt := ":> "
	scanner := uniline.DefaultScanner()

	historyFile := ".ltc.history"
	user, _ := user.Current()
	if user != nil {
		historyFile = user.HomeDir + "/" + historyFile
	}
	scanner.LoadHistory(historyFile)

	for scanner.Scan(prompt) {
		line := scanner.Text()
		if len(line) > 0 {
			// scanner.AddToHistory(line)
			scanner.SaveHistory(historyFile)
			lc.handleLine(scanner, line)

			// tokenAuth := ""
			// if lc.userToken != nil {
			//     tokenAuth = lc.userToken.Authority
			// }
			prompt = fmt.Sprintf("%s> ", lc.stateName)
		}
	}
}

const (
	s_arg_done = iota
	s_in_arg
	s_in_quote
)

func parseLine(line string) []string {

	args := make([]string, 0)

	state := s_arg_done

	accum := make([]rune, 0, len(line))

	var quoteChar rune

	for x, c := range line {
		_ = x
		// fmt.Fprintf(os.Stderr, ansi.Color("x=%d accum=%s\n", "magenta"), x, string(accum))

		switch state {
		case s_arg_done:
			switch {
			case unicode.IsSpace(c):
				// Stay in this state
				continue

			case c == '\'' || c == '"' || c == '`':
				accum = append(accum, c)
				quoteChar = c
				state = s_in_quote

			default:
				accum = append(accum, c)
				state = s_in_arg

			}

		case s_in_arg:
			if unicode.IsSpace(c) {
				args = append(args, string(accum))
				accum = accum[:0]
				state = s_arg_done
			} else {
				accum = append(accum, c)
			}

		case s_in_quote:
			prevChar := ' '
			if len(accum) > 0 {
				prevChar = accum[len(accum)-1]
			}

			if c == quoteChar {
				accum = append(accum, c)
				if prevChar != '\\' {
					// End of the quoted part, so we convert that, but
					// it might not be the total end of the arg
					unquoted, e := strconv.Unquote(string(accum))
					_ = e
					// fmt.Fprintf(os.Stderr, "To unqote %v\n", string(accum))
					// fmt.Fprintf(os.Stderr, "Result = %v\n", unquoted)
					// fmt.Fprintf(os.Stderr, "e = %v\n", e)
					accum = accum[:0]
					accum = append(accum, ([]rune(unquoted))...)
					state = s_in_arg
				}
			} else if prevChar == '\\' {
				// fmt.Fprintf(os.Stderr, "len=%d prevChar=%v  c=%v\n", len(accum), prevChar, c)
				// We know c isn't the quote char - so it if it is an alternate
				// quote char, that's not allowed
				if c == '\'' || c == '"' || c == '`' {
					// These are not allowed, so nuke prev char before
					// appending this char
					accum = accum[:len(accum)-1]
				}
				accum = append(accum, c)
			} else {
				accum = append(accum, c)
			}
		}
	}

	if len(accum) > 0 {
		if state == s_in_quote {
			if accum[len(accum)-1] != quoteChar {
				accum = append(accum, quoteChar)
			}
			unquoted, _ := strconv.Unquote(string(accum))
			accum = []rune(unquoted)
		}
		args = append(args, string(accum))
	}

	return args
}

func (lc *localContext) handleLine(scanner *uniline.Scanner, line string) {

	args := parseLine(line)

	if len(args) == 0 {
		return
	}

	// Add it to history even if we don't know the command
	scanner.AddToHistory(line)

	handler := commands[args[0]]
	if handler == nil {
		fmt.Fprintf(os.Stderr, "unrecognized command '%s'\n", args[0])
		return
	}

	handler(lc, args[1:])
}

var cCmd = ansi.ColorFunc("red")
var cShort = ansi.ColorFunc("green")
var cSyntax = ansi.ColorFunc("")

func shortHelp(name string, entry *helpEntry) {
	fmt.Fprintf(os.Stderr, "%s - %s\n", cCmd(name), cShort(entry.short))
	// fmt.Fprintf(os.Stderr, "\t"+cSyntax("%s %s\n")+"\n", name, entry.syntax)
}

func longHelp(lc *localContext, args []string) {

	fmt.Fprintln(os.Stderr)

	if len(args) == 0 {
		// Show all available commands
		names := make([]string, 0, len(help))
		for k, _ := range help {
			names = append(names, k)
		}
		sort.Strings(names)

		for _, name := range names {
			entry := help[name]
			shortHelp(name, entry)
		}
		return
	}

	for _, v := range args {
		entry := help[v]
		if entry != nil {
			shortHelp(v, entry)
			fmt.Fprintf(os.Stderr, "\t"+cSyntax("%s %s\n")+"\n", v, entry.syntax)

			if len(entry.man) > 0 {
				fmt.Fprintln(os.Stderr)
				fmt.Fprintln(os.Stderr, entry.man)
				fmt.Fprintln(os.Stderr)
				continue
			}
		} else {
			fmt.Fprintf(os.Stderr, "no help for '%s'\n", v)
		}
	}
}

func printErr(msg string, err error) {
	fmt.Println(aRed, msg, ": ", err, aReset)
}

func failedTo(msg string, err error) bool {
	if err == nil {
		return false
	}

	printErr(msg, err)
	return true
}

func init() {
	commands["help"] = longHelp
	commands["?"] = longHelp

	commands["args"] = func(lc *localContext, args []string) {
		for ix, a := range args {
			fmt.Printf("%v: '%v'\n", ix, a)
		}
	}

	commands["exit"] = func(lc *localContext, args []string) {
		os.Exit(0)
	}

	commands["quit"] = func(lc *localContext, args []string) {
		os.Exit(0)
	}
}
