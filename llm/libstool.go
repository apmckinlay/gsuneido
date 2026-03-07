// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package llm

import (
	"context"

	"github.com/apmckinlay/gsuneido/core"
)

// libraries
var _ = addTool(toolSpec{
	name:        "suneido_libraries",
	description: "Get a list of the libraries currently in use in Suneido",
	handler: func(ctx context.Context, args map[string]any) (any, error) {
		libs := core.GetDbms().Libraries()
		return librariesOutput{Libraries: libs}, nil
	},
})

type librariesOutput struct {
	Libraries []string `json:"libraries" jsonschema:"List of libraries currently in use"`
}
