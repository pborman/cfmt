// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package comment

import (
	"bytes"
	"testing"
)

// Note that the first byte of in and out are ignored.  This makes it
// easier to write our input with `
var tests = []struct {
	in    string
	out   string
	delim string
}{
	{`
`, `
`,
		"//",
	},
	{`
text
`, `
text
`,
		"//",
	},
	{`
text`, `
text`,
		"//",
	},
	{`
// comment
`, `
// comment
`,
		"//",
	},
	{`
// comment`, `
// comment
`,
		"//",
	},
	{`
// 456789 bcdef 123
`, `
// 456789 bcdef
// 123
`,
		"//",
	},
	{`
// 456789 bcdef     123
`, `
// 456789 bcdef
// 123
`,
		"//",
	},
	{`
// 456789 bcdef0 123
`, `
// 456789 bcdef0
// 123
`,
		"//",
	},
	{`
// 456789 bcdef0      123
`, `
// 456789 bcdef0
// 123
`,
		"//",
	},
	{`
// 456789 bcdef01 123
`, `
// 456789
// bcdef01 123
`,
		"//",
	},
	{`
// 456789 bcdef01  123
`, `
// 456789
// bcdef01  123
`,
		"//",
	},
	{`
// 456789abcdef01234  123
`, `
// 456789abcdef01234
// 123
`,
		"//",
	},
	{`
//  456789abcdef01234  123
`, `
//  456789abcdef01234  123
`,
		"//",
	},
	{`
//  456789 bcdef01234  123
// line 1
// line 2
`, `
//  456789 bcdef01234  123
// line 1 line 2
`,
		"//",
	},
	{`
two
lines
`, `
two
lines
`,
		"//",
	},
	// This one has a space at the end of each line
	{`
two 
lines 
`, `
two 
lines 
`,
		"//",
	},
	{`
// two
lines
`, `
// two
lines
`,
		"//",
	},
	// This one has a space at the end of each input line
	{`
// two 
lines 
`, `
// two
lines 
`,
		"//",
	},
	{`
// two
// lines
`, `
// two lines
`,
		"//",
	},
	// This one has a space at the end of each input line
	{`
// two 
// lines 
`, `
// two lines
`,
		"//",
	},
	{`
// line1
line2
`, `
// line1
line2
`,
		"//",
	},
	{`
// line1
// line2
`, `
// line1 line2
`,
		"//",
	},

	{`
// line1
//  line2
//  line3
// line4
line5
`, `
// line1
//  line2
//  line3
// line4
line5
`,
		"//",
	},

	{`
// line1
// line2
// line3
// line4
line5
`, `
// line1 line2
// line3 line4
line5
`,
		"//",
	},
	{`
// line1
//line2
//line3
// line4
line5
`, `
// line1 line2
// line3 line4
line5
`,
		"//",
	},
	{`
// line1.
// line2
`, `
// line1.  line2
`,
		"//",
	},
	{`
  // ln1.
  // ln2
`, `
  // ln1.  ln2
`,
		"//",
	},
	{`
  // lin1
  // lin2
  // lin3
  // lin4
`, `
  // lin1 lin2
  // lin3 lin4
`,
		"//",
	},
	{`
  // line1
  // lin2
  // lin3
  // lin4
`, `
  // line1 lin2
  // lin3 lin4
`,
		"//",
	},
	{`
  // line1
  // line2
  // lin3
  // lin4
`, `
  // line1 line2
  // lin3 lin4
`,
		"//",
	},
	{`
  // line1.
  // line2
`, `
  // line1.
  // line2
`,
		"//",
	},
	{`
// line0
  // ln1
  // ln2
`, `
// line0
  // ln1
  // ln2
`,
		"//",
	},
	{`
	// line 0
	// line 1
`, `
	// line
	// 0
	// line
	// 1
`,
		"//",
	},
	{`
	// line-0
	// line-1
`, `
	// line-0
	// line-1
`,
		"//",
	},
	{`
	// 0
	// 1
	// 2
`, `
	// 0 1 2
`,
		"//",
	},
	{`
// a
// cmt
`,
		`
// a cmt
`,
		"",
	},
	{`
	// a
	// cmt
`,
		`
	// a cmt
`,
		"",
	},
	{`
# a
# cmt
`,
		`
# a cmt
`,
		"#",
	},
	{`
	# a
	# cmt
`,
		`
	# a cmt
`,
		"#",
	},
	{`
# a
# cmt
`,
		`
# a cmt
`,
		"",
	},
	{`
	# a
	# cmt
`,
		`
	# a cmt
`,
		"",
	},
	{`
	// a
	// b
	//exp c
`,
		`
	// a b
	//exp c
`,
		"",
	},
}

func TestFormat(t *testing.T) {
	export = "exp" // keep lines short in testing
	defer func() {
		if x := recover(); x != nil {
			t.Fatalf("%v\n", x)
		}
	}()
	for n, tt := range tests {
		tt.in = tt.in[1:]
		tt.out = tt.out[1:]
		outb := &bytes.Buffer{}
		Format(outb, bytes.NewBufferString(tt.in), 16, tt.delim)
		out := outb.String()

		if out != tt.out {
			t.Errorf("#%d: in=%q delim=%q out=%q want %q\nin:\n%s\nout:\n%s\nwant:\n%s", n, tt.in, tt.delim, out, tt.out, tt.in, out, tt.out)
		}
	}
}
