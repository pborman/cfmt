// cfmt reformats its input such that the initial block comments will fit in 80
// columns.  All input following the first non-block comment line is copied to
// output without change.
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/pborman/cfmt/comment"
	"github.com/pborman/options"
)

func main() {
	var flags = struct {
		Delim  string `getopt:"--delim=STRING -d comment delimiter (defaults to autodetect)"`
		Spell  bool   `getopt:"--spell -s run aspell on the comment"`
		Target int    `getopt:"--target=N -n target line length of comments"`
	}{
		Target: 80,
	}
	options.RegisterAndParse(&flags)
	var text, code bytes.Buffer
	err := comment.Format2(&text, &code, os.Stdin, flags.Target, flags.Delim)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if _, err := os.Stdout.Write(text.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	if flags.Spell {
		var text2 bytes.Buffer
		cmd := exec.Command("aspell", "list")
		cmd.Stdin = &text
		cmd.Stdout = &text2
		cmd.Stderr = &text2
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "aspell: %v\n", err)
		}
		words := unique(strings.Fields(text2.String()))
		if len(words) > 0 {
			fmt.Println("-- start of misspelled words ---")
			fmt.Println(strings.Join(words, "\n"))
			fmt.Println("-- end of misspelled words ---")
		}
	}
	if _, err := io.Copy(os.Stdout, &code); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func unique(words []string) []string {
	sort.Strings(words)
	h := 0
	for x, w := range words {
		if x == 0 || words[h] == w {
			continue
		}
		h++
		words[h] = w
	}
	return words[:h+1]
}
