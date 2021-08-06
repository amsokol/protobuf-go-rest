package main

import (
	"flag"
	"fmt"
	"os"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	_version = "v0.1.0"
)

func main() {
	// useInputFile := flag.String("input_file", "", "use file instead of stdin")
	showVersion := flag.Bool("version", false, "print the current version")

	flag.Parse()

	if *showVersion {
		fmt.Printf("protoc-gen-protobuf-rest %v\n", _version)

		return
	}

	/*
		// safe input data for DEBUG
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return
		}
		_ = os.WriteFile("./stdin.debug", data, 0666)
		return
	*/

	r, err := os.Open("../../stdin.debug")
	if err != nil {
		return
	}

	protogen.Options{
		ParamFunc: flag.CommandLine.Set,
		DebugFile: r,
	}.Run(func(gen *protogen.Plugin) error {
		gen.SupportedFeatures = uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
		for _, f := range gen.Files {
			if !f.Generate {
				continue
			}
			generateFile(gen, f)
		}

		return nil
	})
}
