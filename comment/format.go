package comment

import (
	"bufio"
	"io"
	"strings"
)

var comment string // where we build up our comment to format
var export = "export "

// An error is provided to panic if write fails.
type error_ struct {
	error
}

// write writes the string s to out using out.Write.  If write is unable
// to write all of s it will panic with an error (which is caught by Format).
func write(out *bufio.Writer, s string) {
	b := []byte(s)
	for len(b) > 0 {
		n, err := out.Write(b)
		if n == len(b) {
			return
		}
		if err != nil {
			panic(error_{err})
		}
		b = b[n:]
	}
}

// isspace returns true if r is a white space character.
func isspace(r byte) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

// dump formats any comment collected, targetting lines of target length and
// writes it to out.  If line is not empty it is also written out.  Any error in
// writing will cause a panic (which is caught by Format).
func dump(out, cout *bufio.Writer, prefix, line string, target int) {
	for len(comment)+1 > target {
		var j int
		for j = target - 1; j > 0; j-- {
			if isspace(comment[j]) {
				break
			}
		}
		if j == 0 {
			for j = target; j < len(comment); j++ {
				if isspace(comment[j]) {
					break
				}
			}
		}
		if j == len(comment) {
			break
		}
		for j--; isspace(comment[j]); j-- {
		}
		write(out, prefix+" "+comment[:j+1]+"\n")
		for j++; isspace(comment[j]); j++ {
		}
		comment = comment[j:]
	}
	if comment != "" {
		write(out, prefix+" "+comment+"\n")
		comment = ""
	}
	if strings.HasPrefix(line, prefix) {
		write(out, line)
	} else if line != "" {
		write(cout, line)
	}
}

// Format reads text from r, formats any leading block comment, and writes it to
// w.  Any text following the initial block comment (that is all text starting
// with the first line that does not start with delim) is copied verbatim.
// Block comments are formatted to not be more than target bytes long.  All
// trailing spaces are removed from the block comment.  Only lines with up to
// one space between delim and the first non-space character are formatted.  A
// line consisting of just delim or has more than a single space following delim
// are copied verbatim.
//
// If delim is empty then Format attempts to autodetect the comment string based
// on the first line (only # and // are recognized).
//
// Format is wrapped by the command cfmt.  It is designed to be used with vi.
// Place the cursor on the first line of a block comment and type
// ``!}cfmt<ENTER>'' to have that block comment formatted.
func Format(w io.Writer, r io.Reader, target int, delim string) (err error) {
	return Format2(w,w,r,target,delim)
}

// Format2 is like Format with the exception that two io.Writers are provided.
// The formatted comment is written to the first io.Writer and the remaining
// text (code) is written to the second.
func Format2(w, cw io.Writer, r io.Reader, target int, delim string) (err error) {
	out := bufio.NewWriter(w)
	cout := out
	if w != cw {
		cout = bufio.NewWriter(cw)
	}
	in := bufio.NewReader(r)
	defer func() {
		var err1 error
		if p := recover(); p != nil {
			var ok bool
			if err1, ok = p.(error_); !ok {
				panic(p)
			}
		}
		err2 := out.Flush()
		var err3 error
		if out != cout {
			err3 = cout.Flush()
		}
		if err == nil {
			if err1 != nil {
				err = err1
			} else if err2 != nil {
				err = err2
			} else if err3 != nil {
				err = err2
			}
		}
	}()

	var last *string
	var lerr error
	next := func() (string, error) {
		if last != nil {
			defer func() { last = nil }()
			return *last, lerr
		}
		return in.ReadString('\n')

	}

	line, lerr := next()
	last = &line // force next to return this line again
	i := 0
	switch {
	case len(line) == 0:
	case line[0] == ' ':
		for len(line) > i && line[i] == ' ' {
			target--
			i++
		}
	case line[0] == '\t':
		for len(line) > i && line[i] == '\t' {
			target -= 8
			i++
		}
	}
	if delim == "" {
		switch {
		case strings.HasPrefix(line[i:], "-- "):
			delim = "--"
		case strings.HasPrefix(line[i:], "//"):
			delim = "//"
		case strings.HasPrefix(line[i:], "#"):
			delim = "#"
		default:
			delim = "//"
		}
	}
	prefix := delim
	if strings.HasPrefix(line[i:], prefix) {
		target -= len(prefix)
		prefix = line[:i] + prefix
	}

	hasPrefix := func(s string) bool {
		if !strings.HasPrefix(s, prefix) {
			return false
		}
		return !strings.HasPrefix(s[len(prefix):], export)
	}

	plen := len(prefix)
	for {
		line, err := next()
		if (err != nil && err != io.EOF) || !hasPrefix(line) {
			dump(out, cout, prefix, line, target)
			if err != nil && err != io.EOF {
				return err
			}
			break
		}
		// strip trailing whitespace
		// after this line will be "//" or will have at least one
		// non-white space character.
		for isspace(line[len(line)-1]) {
			line = line[:len(line)-1]
		}

		// If we have an empty // line then dump the comments thus far
		// and start fresh
		if line == prefix {
			dump(out, cout, prefix, line+"\n", target)
		} else {
			// Trim the leading // and the first space (if any)
			if line[plen] == ' ' {
				line = line[plen+1:]
			} else {
				line = line[plen:]
			}

			// If we start with a tab then output what we got
			if line[0] == '\t' {
				dump(out, cout, prefix, prefix+line+"\n", target)
			} else if line[0] == ' ' {
				dump(out, cout, prefix, prefix+" "+line+"\n", target)
			} else {
				// two spaces after a .
				if len(comment) > 0 && comment[len(comment)-1] == '.' {
					comment += " "
				}
				if comment == "" {
					comment = line
				} else {
					comment += " " + line
				}
			}
		}
		if err == io.EOF {
			dump(out, cout, prefix, "", target)
			return nil
		}
	}
	_, err = io.Copy(cout, in)
	return err
}
