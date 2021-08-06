package runtime

import "strings"

type Values map[string]string

func (vv Values) New(key string, value string, add bool) {
	if add {
		if e, ok := vv[key]; ok {
			value = strings.Join([]string{e, value}, "/")
		}
	}

	vv[key] = value
}
