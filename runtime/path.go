package runtime

import (
	"fmt"
	"strings"
)

type Path = *Segment

/*
Path template syntax

Template = "/" Segments [ Verb ] ;
Segments = Segment { "/" Segment } ;
Segment  = "*" | "**" | LITERAL | Variable ;
Variable = "{" FieldPath [ "=" Segments ] "}" ;
FieldPath = IDENT { "." IDENT } ;
Verb     = ":" LITERAL ;

The syntax `*` matches a single URL path segment.
The syntax `**` matches zero or more URL path segments,
which must be the last part of the URL path except the `Verb`.

https://pkg.go.dev/google.golang.org/genproto/googleapis/api/annotations
*/
func NewPath(template string) (Path, error) {
	// normalize template
	t := strings.TrimSpace(template)
	t = strings.Trim(t, "/")
	t = strings.ToLower(t) + "/"

	var (
		b strings.Builder
		v bool
		p Path
		s *Segment
	)

	b.Grow(len(t))

	for i, c := range t {
		// end of segment
		if c == '/' && !v {
			// create new segment
			cs, ce, err := NewSegmentChain(b.String())
			if err != nil {
				return nil, fmt.Errorf("create new path segment chain '%s': %w", template, err)
			}

			if p == nil {
				// save root segment (=path)
				p = cs
			} else {
				// add new segment to chain
				s.Next = cs
			}

			s = ce

			// start new segment
			b.Reset()
			b.Grow(len(t) - i)

			continue
		}

		b.WriteRune(c)

		switch c {
		case '{':
			v = true
		case '}':
			v = false
		}
	}

	return p, nil
}
