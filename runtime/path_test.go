package runtime_test

import (
	"reflect"
	"testing"

	"github.com/amsokol/protobuf-rest/runtime"
)

func TestNewPath(t *testing.T) {
	type args struct {
		template string
	}

	tests := []struct {
		name    string
		args    args
		want    runtime.Path
		wantErr bool
	}{
		{
			"",
			args{
				template: "",
			},
			&runtime.Segment{},
			false,
		},
		{
			"/",
			args{
				template: "/",
			},
			&runtime.Segment{},
			false,
		},
		{
			"/v1",
			args{
				template: "/v1",
			},
			&runtime.Segment{
				Value: "v1",
			},
			false,
		},
		{
			"/v1/articles",
			args{
				template: "/v1/articles",
			},
			&runtime.Segment{
				Value: "v1",
				Next: &runtime.Segment{
					Value: "articles",
				},
			},
			false,
		},
		{
			"/v1/articles/{value}",
			args{
				template: "/v1/articles/{value}",
			},
			&runtime.Segment{
				Value: "v1",
				Next: &runtime.Segment{
					Value: "articles",
					Next: &runtime.Segment{
						Field: "value",
						Value: "*",
					},
				},
			},
			false,
		},
		{
			"/v1/articles/{value}/data",
			args{
				template: "/v1/articles/{value}/data",
			},
			&runtime.Segment{
				Value: "v1",
				Next: &runtime.Segment{
					Value: "articles",
					Next: &runtime.Segment{
						Field: "value",
						Value: "*",
						Next: &runtime.Segment{
							Value: "data",
						},
					},
				},
			},
			false,
		},
		{
			"/v1/articles/{value=data/*}",
			args{
				template: "/v1/articles/{value=data/*}",
			},
			&runtime.Segment{
				Value: "v1",
				Next: &runtime.Segment{
					Value: "articles",
					Next: &runtime.Segment{
						Field: "value",
						Value: "data",
						Next: &runtime.Segment{
							Field: "value",
							Value: "*",
							IsVal: true,
						},
					},
				},
			},
			false,
		},
		{
			"/v1/books/articles/{value=data/items/*}",
			args{
				template: "/v1/books/articles/{value=data/items/*}",
			},
			&runtime.Segment{
				Value: "v1",
				Next: &runtime.Segment{
					Value: "books",
					Next: &runtime.Segment{
						Value: "articles",
						Next: &runtime.Segment{
							Field: "value",
							Value: "data",
							Next: &runtime.Segment{
								Field: "value",
								Value: "items",
								IsVal: true,
								Next: &runtime.Segment{
									Field: "value",
									Value: "*",
									IsVal: true,
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"/v1/books/articles/{value=data/items/*}/symbol/{number}",
			args{
				template: "/v1/books/articles/{value=data/items/*}/symbol/{number}",
			},
			&runtime.Segment{
				Value: "v1",
				Next: &runtime.Segment{
					Value: "books",
					Next: &runtime.Segment{
						Value: "articles",
						Next: &runtime.Segment{
							Field: "value",
							Value: "data",
							Next: &runtime.Segment{
								Field: "value",
								Value: "items",
								IsVal: true,
								Next: &runtime.Segment{
									Field: "value",
									Value: "*",
									IsVal: true,
									Next: &runtime.Segment{
										Value: "symbol",
										Next: &runtime.Segment{
											Field: "number",
											Value: "*",
										},
									},
								},
							},
						},
					},
				},
			},
			false,
		},
		{
			"/v1/tables/*",
			args{
				template: "/v1/tables/*",
			},
			&runtime.Segment{
				Value: "v1",
				Next: &runtime.Segment{
					Value: "tables",
					Next: &runtime.Segment{
						Value: "*",
					},
				},
			},
			false,
		},
		{
			"/v1/tables/**",
			args{
				template: "/v1/tables/**",
			},
			&runtime.Segment{
				Value: "v1",
				Next: &runtime.Segment{
					Value: "tables",
					Next: &runtime.Segment{
						Value: "**",
					},
				},
			},
			false,
		},
		{
			"/ /",
			args{
				template: "/ /",
			},
			nil,
			true,
		},
		{
			"/v1/articles/{value=/}",
			args{
				template: "/v1/articles/{value=/}",
			},
			nil,
			true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := runtime.NewPath(tt.args.template)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewPath() error = %v, wantErr %v", err, tt.wantErr)

				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPath() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
