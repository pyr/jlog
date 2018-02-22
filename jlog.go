package main

import (
	"bufio"
	"os"
	"flag"
	"fmt"
	"encoding/json"
	"syscall"
	"time"
	"unsafe"
)

var noColor = flag.Bool("c", false, "Disable colored output")
var showExtra = flag.Bool("X", false, "Show extra keys in JSON messages")

// UNIX test to see if we are talking to a human or a file
func IsStdoutATerminal() bool {
	var termios syscall.Termios
	_, _, err := syscall.Syscall6(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		syscall.TCGETS,
		uintptr(unsafe.Pointer(&termios)),
		0,
		0,
		0)
	return err == 0
}

type Color int

const (
	COLOR_GREEN Color = iota
	COLOR_YELLOW
	COLOR_RED
	COLOR_NONE
)

var COLOR_CODES = map[Color]string{
	COLOR_GREEN: "32",
	COLOR_YELLOW: "33",
	COLOR_RED: "31",
}

func Colorize(color Color, text string) string {
	if code, ok := COLOR_CODES[color]; ok {
		return fmt.Sprintf("[%sm%s[0m", code, text)
	}
	return text
}

func KeepKey(k string) bool {
	return (k != "message" &&
		k != "level" &&
		k != "logger" &&
		k != "exception" &&
		k != "logger_name" &&
		k != "stack_trace" &&
		k != "level_value" &&
		k[0] != '@')
}

func ShowLine(colorOutput bool, showExtra bool, line string) {

	var msg map[string]interface{}
	err := json.Unmarshal([]byte(line), &msg)
	if err != nil {
		return
	}

	fields, ok := msg["@fields"].(map[string]interface{})
	if ok {
		fields["@message"] = msg["@message"]
		fields["@timestamp"] = msg["@timestamp"]
		msg = fields
	}

	color := COLOR_NONE;
	level,_ := msg["level"].(string)
	if (colorOutput) {
		if (level == "INFO") {
			color = COLOR_GREEN
		} else if (level == "WARN") {
			color = COLOR_YELLOW
		} else if (level == "ERROR") {
			color = COLOR_RED
		}
	}

	if level != "" {
		level = fmt.Sprintf("%-5s ", level)
		level = Colorize(color, level)
	}
	tmstr, ok := msg["@timestamp"].(string)
	if ok {
		tm, err := time.Parse(time.RFC3339Nano, tmstr)
		if err == nil {
			tmstr = tm.Format(time.Stamp)
		} else {
			tmstr = ""
		}
	} else {
		tmstr = ""
	}

	src, ok  := msg["logger_name"].(string)
	if !ok {
		src, _ = msg["logger"].(string)
	}

	if src != "" {
		src = fmt.Sprintf(" %s -", src)
	}

	payload, ok := msg["message"].(string)
	if !ok {
		payload, _ = msg["@message"].(string)
	}
	payload = Colorize(color, payload)
	out := fmt.Sprintf("%s[%s]%s %s", level, tmstr, src, payload)
	if showExtra {
		mdc, ok := msg["mdc"]
		if ok {
			mdcmap, _ := mdc.(map[string]interface{})
			for k, v := range mdcmap {
				out = fmt.Sprintf("%s %s=%v", out, k, v)
			}
		} else {
			for k, v := range msg {
				if KeepKey(k) {
					out = fmt.Sprintf("%s %s=%v", out, k, v)
				}
			}
		}
	}
	out = fmt.Sprintf("%s\n", out)
	stacktrace, ok := msg["stack_trace"]
	if ok {
		stacktrace, _ = stacktrace.(string)
		out = fmt.Sprintf("%s\t%s", out, stacktrace)
	} else {
		submap, ok := msg["exception"].(map[string]interface{})
		if ok {
			stacktrace, st_ok := submap["stacktrace"]
			if (st_ok) {
				out = fmt.Sprintf("%s%s\n", out, stacktrace)
			}
		}
	}
	fmt.Print(out)
}

func main() {
	flag.Parse()

	args := flag.Args()
	var scanner *bufio.Scanner
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "too many arguments")
		os.Exit(1)
	} else if len(args) == 0 {
		scanner = bufio.NewScanner(os.Stdin)
	} else {
		file, err := os.Open(args[0])
		if err != nil {
			fmt.Fprintln(os.Stderr, "cannot open file: %v", err)
			os.Exit(1)
		}
		scanner = bufio.NewScanner(file)
	}
	colorOutput := IsStdoutATerminal() && !*noColor
	for scanner.Scan() {
		line := scanner.Text()
		ShowLine(colorOutput, *showExtra, line)
	}
	os.Exit(0)
}
