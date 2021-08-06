package runtime

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type Segment struct {
	Value string   // url/value
	Field string   // segment field name
	IsVal bool     // indicates that segment is part of the "{field=value}" pattern value
	Next  *Segment // next segment
}

func (s *Segment) Match(splittedPath []string, values Values) bool {
	switch {
	case strings.EqualFold(s.Value, "**"):
		return s.doDoubleStar(splittedPath, values)
	case strings.EqualFold(s.Value, "*"):
		return s.doStar(splittedPath, values)
	}

	if len(splittedPath) == 0 || !strings.EqualFold(s.Value, splittedPath[0]) {
		// not matched
		return false
	}

	// matched!

	if len(s.Field) > 0 {
		// this is field value
		values.New(s.Field, splittedPath[0], s.IsVal)
	}

	if s.Next == nil {
		// last segment in template
		return len(splittedPath) == 1
	}

	// there are more segments in template
	return s.Next.Match(splittedPath[1:], values)
}

// Match: **.
func (s *Segment) doDoubleStar(splittedPath []string, v Values) bool {
	if len(s.Field) > 0 {
		// this is field value
		switch l := len(splittedPath); l {
		case 0:
			// empty value - nothing to do
		case 1:
			v.New(s.Field, splittedPath[0], s.IsVal)
		default:
			// compose field value
			var (
				val strings.Builder
				i   int
			)

			// try to estimate buffer size to minimize allocations
			val.Grow(len(splittedPath[0]) * l)

			for i < l {
				if i > 0 {
					val.WriteRune('/')
				}

				val.WriteString(splittedPath[i])
				i++
			}

			v.New(s.Field, val.String(), s.IsVal)
		}
	}

	return true
}

// Match: *.
func (s *Segment) doStar(splittedPath []string, v Values) bool {
	if s.Next == nil {
		// last segment of template
		switch l := len(splittedPath); l {
		case 1:
			if len(s.Field) > 0 {
				// this is field value
				v.New(s.Field, splittedPath[0], s.IsVal)
			}

			fallthrough
		case 0:
			// empty value - nothing to do
			return true
		}
		// url has more segments - not matched
		return false
	}

	// there is next segment in template - try to move inside
	if len(splittedPath) > 0 {
		if len(s.Field) > 0 {
			// this is field value
			v.New(s.Field, splittedPath[0], s.IsVal)
		}
		// move inside
		return s.Next.Match(splittedPath[1:], v)
	}

	// move inside
	return s.Next.Match(splittedPath, v)
}

// 0 - matched string
// 1 - {"field"=value}
// 2 - {field="value"}
// 3 - {"field"}
// 4 - "url".
var _ReSegment = regexp.MustCompile(`^{(\S+?)=(\S+?)}$|^{(\S+?)}$|^([^{}\s]+)$`)

func NewSegmentChain(pathSegment string) (*Segment, *Segment, error) {
	if len(pathSegment) == 0 {
		// empty segment is OK, e.g. root "/"
		s := &Segment{}

		return s, s, nil
	}

	// extract segment parts
	mm := _ReSegment.FindSubmatch([]byte(pathSegment))
	if mm == nil || len(mm[0]) == 0 {
		return nil, nil, fmt.Errorf("%w: '%s'", ErrInvalidSegmentFormat, pathSegment)
	}

	switch {
	case len(mm[1]) > 0 && len(mm[2]) > 0:
		// {"field"="value"} pattern
		return newSegmentFieldValue(string(mm[1]), string(mm[2]))
	case len(mm[3]) > 0:
		// {"field"} pattern
		return newSegmentField(string(mm[3]))
	default:
		// "url" pattern
		return newSegmentURL(string(mm[4]))
	}
}

func newSegmentURL(pathSegment string) (*Segment, *Segment, error) {
	s := &Segment{Value: pathSegment}

	return s, s, nil
}

func newSegmentField(field string) (*Segment, *Segment, error) {
	s := &Segment{Value: "*", Field: field}

	return s, s, nil
}

func newSegmentFieldValue(field string, value string) (*Segment, *Segment, error) {
	v := strings.Trim(value, "/")
	if len(v) == 0 {
		return nil, nil, fmt.Errorf("%w: '%s'", ErrInvalidFieldValueFormat, value)
	}

	var cs, ce *Segment

	for i, p := range strings.Split(v, "/") {
		s := &Segment{Value: p, Field: field}
		if i == 0 {
			cs = s
		} else {
			s.IsVal = true
			ce.Next = s
		}

		ce = s
	}

	return cs, ce, nil
}

var (
	ErrInvalidSegmentFormat    = errors.New("invalid url segment format")
	ErrInvalidFieldValueFormat = errors.New("invalid format of field value template")
)
