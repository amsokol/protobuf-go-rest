package runtime_test

import (
	"reflect"
	"testing"

	"github.com/amsokol/protobuf-rest/runtime"
)

func TestValues_New(t *testing.T) {
	type args struct {
		key   string
		value string
		add   bool
	}

	tests := []struct {
		name string
		vv   runtime.Values
		args args
		want runtime.Values
	}{
		{
			"add value",
			make(runtime.Values),
			args{
				"key1",
				"value1",
				false,
			},
			runtime.Values{
				"key1": "value1",
			},
		},
		{
			"replace value",
			runtime.Values{
				"key1": "value1",
			},
			args{
				"key1",
				"value2",
				false,
			},
			runtime.Values{
				"key1": "value2",
			},
		},
		{
			"add complex value (two parts)",
			runtime.Values{
				"key1": "value1",
			},
			args{
				"key1",
				"value2",
				true,
			},
			runtime.Values{
				"key1": "value1/value2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.vv.New(tt.args.key, tt.args.value, tt.args.add)
			if !reflect.DeepEqual(tt.vv, tt.want) {
				t.Errorf("Values.New() got = %v, want %v", tt.vv, tt.want)
			}
		})
	}
}
