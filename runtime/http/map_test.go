package http_test

import (
	"context"
	"math/rand"
	"net/http"
	"reflect"
	"testing"

	"github.com/amsokol/protobuf-rest/runtime"
	_http "github.com/amsokol/protobuf-rest/runtime/http"
)

func TestMap_Add(t *testing.T) {
	type args struct {
		method   string
		template string
		handler  _http.Handler
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"GET:/v1/articles",
			args{
				"GET",
				"/v1/articles",
				func(context.Context, http.ResponseWriter, *http.Request) {},
			},
			false,
		},
		{
			"GET:/ /",
			args{
				"GET",
				"/ /",
				func(context.Context, http.ResponseWriter, *http.Request) {},
			},
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := _http.NewMap()

			if err := m.Add(tt.args.method, tt.args.template, tt.args.handler); (err != nil) != tt.wantErr {
				t.Errorf("Map.Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMap_Match(t *testing.T) {
	paths := []struct {
		method   string
		template string
		handler  _http.Handler
	}{
		{
			"GET",
			"/v1/articles",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
	}

	type args struct {
		method  string
		urlPath string
	}

	tests := []struct {
		name  string
		args  args
		want  _http.Handler
		want1 runtime.Values
	}{
		{
			"GET:/v1/articles",
			args{
				"GET",
				"/v1/articles",
			},
			func(context.Context, http.ResponseWriter, *http.Request) {},
			runtime.Values{},
		},
		{
			"POST:/v1/articles",
			args{
				"POST",
				"/v1/articles",
			},
			nil,
			nil,
		},
		{
			"GET:/v1/article",
			args{
				"GET",
				"/v1/article",
			},
			nil,
			nil,
		},
	}

	m := _http.NewMap()
	for _, p := range paths {
		if err := m.Add(p.method, p.template, p.handler); err != nil {
			t.Fatal(err)
		}
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := m.Match(tt.args.method, tt.args.urlPath)

			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
				t.Errorf("Map.Match() got = %#v, want %#v", got, tt.want)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("Map.Match() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func BenchmarkPathMap_Match(b *testing.B) {
	paths := []struct {
		method   string
		template string
		handler  _http.Handler
	}{
		{
			"GET",
			"",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1/articles",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1/articles/{value}",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1/articles/{value}/data",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1/articles/{value=data/*}",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1/articles/{value=data1/*/*/*}",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1/articles/{value=data2/symbol/**}",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1/books/articles/{value=data/items/*}",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1/books/articles/{value=data/items/*}/symbol/{number}",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1/tables/*",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
		{
			"GET",
			"/v1/tables/**",
			func(context.Context, http.ResponseWriter, *http.Request) {},
		},
	}

	type args struct {
		method  string
		urlPath string
	}

	tests := []struct {
		name  string
		args  args
		want  _http.Handler
		want1 runtime.Values
	}{
		{
			"GET:/v1/articles",
			args{
				"GET",
				"/v1/articles",
			},
			func(context.Context, http.ResponseWriter, *http.Request) {},
			runtime.Values{},
		},
		{
			"GET:/v1/articles/12345",
			args{
				"GET",
				"/v1/articles/12345",
			},
			func(context.Context, http.ResponseWriter, *http.Request) {},
			runtime.Values{
				"value": "12345",
			},
		},
		{
			"GET:/v1/articles/data/12345",
			args{
				"GET",
				"/v1/articles/data/12345",
			},
			func(context.Context, http.ResponseWriter, *http.Request) {},
			runtime.Values{
				"value": "data/12345",
			},
		},
		{
			"GET:/v1/articles/data1/12345",
			args{
				"GET",
				"/v1/articles/data1/12345",
			},
			func(context.Context, http.ResponseWriter, *http.Request) {},
			runtime.Values{
				"value": "data1/12345",
			},
		},
		{
			"GET:/v1/articles/data2/symbol/some_data/12345",
			args{
				"GET",
				"/v1/articles/data2/symbol/some_data/12345",
			},
			func(context.Context, http.ResponseWriter, *http.Request) {},
			runtime.Values{
				"value": "data2/symbol/some_data/12345",
			},
		},
		{
			"GET:/v1/articles/data2/symbol",
			args{
				"GET",
				"/v1/articles/data2/symbol",
			},
			func(context.Context, http.ResponseWriter, *http.Request) {},
			runtime.Values{
				"value": "data2/symbol",
			},
		},
		{
			"GET:/v1/articles/data2/symbol/12345",
			args{
				"GET",
				"/v1/articles/data2/symbol/12345",
			},
			func(context.Context, http.ResponseWriter, *http.Request) {},
			runtime.Values{
				"value": "data2/symbol/12345",
			},
		},
	}

	m := _http.NewMap()
	for _, p := range paths {
		if err := m.Add(p.method, p.template, p.handler); err != nil {
			b.Fatal(err)
		}
	}

	l := len(tests)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tt := tests[rand.Intn(l)]

			// _, _ = m.Match(tt.args.method, tt.args.urlPath)

			/* */
			got, got1 := m.Match(tt.args.method, tt.args.urlPath)

			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
				b.Fatalf("Map.Match() got = %#v, want %#v", got, tt.want)
			}

			if !reflect.DeepEqual(got1, tt.want1) {
				b.Fatalf("Map.Match() got1 = %v, want %v", got1, tt.want1)
			}
			/* */
		}
	})
}
