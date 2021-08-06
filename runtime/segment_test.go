package runtime_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/amsokol/protobuf-rest/runtime"
)

func TestSegment_Match(t *testing.T) {
	ss := []string{
		"",                                     // 0
		"/",                                    // 1
		"/v1",                                  // 2
		"/v1/articles",                         // 3
		"/v1/articles/{value}",                 // 4
		"/v1/articles/{value}/data",            // 5
		"/v1/articles/{value=data/*}",          // 6
		"/v1/articles/{value=data1/*/*/*}",     // 7
		"/v1/articles/{value=data2/symbol/**}", // 8
		"/v1/books/articles/{value=data/items/*}",                 // 9
		"/v1/books/articles/{value=data/items/*}/symbol/{number}", // 10
		"/v1/tables/*",  // 11
		"/v1/tables/**", // 12
	}

	pp := []runtime.Path{}

	for _, s := range ss {
		p, err := runtime.NewPath(s)
		if err != nil {
			t.Fatal(err)
		}

		pp = append(pp, p)
	}

	type args struct {
		path string
	}

	tests := []struct {
		name  string
		args  args
		want  int
		want1 runtime.Values
	}{
		{
			"/v1/articles",
			args{
				"/v1/articles",
			},
			3,
			runtime.Values{},
		},
		{
			"/v1/articles/12345",
			args{
				"/v1/articles/12345",
			},
			4,
			runtime.Values{
				"value": "12345",
			},
		},
		{
			"/v1/articles/data/12345",
			args{
				"/v1/articles/data/12345",
			},
			6,
			runtime.Values{
				"value": "data/12345",
			},
		},
		{
			"/v1/articles/data1/12345",
			args{
				"/v1/articles/data1/12345",
			},
			7,
			runtime.Values{
				"value": "data1/12345",
			},
		},
		{
			"/v1/articles/data2/symbol/some_data/12345",
			args{
				"/v1/articles/data2/symbol/some_data/12345",
			},
			8,
			runtime.Values{
				"value": "data2/symbol/some_data/12345",
			},
		},
		{
			"/v1/articles/data2/symbol",
			args{
				"/v1/articles/data2/symbol",
			},
			8,
			runtime.Values{
				"value": "data2/symbol",
			},
		},
		{
			"/v1/articles/data2/symbol/12345",
			args{
				"/v1/articles/data2/symbol/12345",
			},
			8,
			runtime.Values{
				"value": "data2/symbol/12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				got1 runtime.Values
				got  = -1
			)

			for i, p := range pp {
				got1 = make(runtime.Values)
				if p.Match(strings.Split(strings.Trim(tt.args.path, "/"), "/"), got1) {
					got = i

					break
				}
			}

			if got != tt.want {
				t.Errorf("Segment.Match() = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("PathMap.Match() got1 = %#v, want %#v", got1, tt.want1)
			}
		})
	}
}
