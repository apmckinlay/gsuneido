// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package typechecker

import (
	"fmt"
	"sync"
	"testing"
)

// TypeCheckerGui fans rows out across up to 16 threads, so Process runs
// concurrently in one process. Run under -race.
func TestProcessConcurrent(t *testing.T) {
	const threads = 16
	const iterations = 8
	args := benchArgs(1)
	refs := benchArgs(4)
	for i := range refs {
		refs[i].Name = fmt.Sprintf("Ref%d", i)
	}

	want, err := Process(Request{Method: "TypeInfer", Arguments: args, References: refs})
	if err != nil {
		t.Fatal(err)
	}
	wantInfo, ok := want.Results[0].(TypeInfo)
	if !ok {
		t.Fatalf("unexpected result type %T", want.Results[0])
	}

	var wg sync.WaitGroup
	errs := make(chan error, threads*iterations)
	for range threads {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range iterations {
				got, err := Process(Request{
					Method: "TypeInfer", Arguments: args, References: refs})
				if err != nil {
					errs <- err
					return
				}
				info, ok := got.Results[0].(TypeInfo)
				if !ok {
					errs <- fmt.Errorf("unexpected result type %T", got.Results[0])
					return
				}
				if len(info.Methods) != len(wantInfo.Methods) {
					errs <- fmt.Errorf("method count %d, want %d",
						len(info.Methods), len(wantInfo.Methods))
					return
				}
				for m, vars := range wantInfo.Methods {
					gv, ok := info.Methods[m]
					if !ok {
						errs <- fmt.Errorf("method %q missing", m)
						return
					}
					for k, v := range vars {
						if gv[k] != v {
							errs <- fmt.Errorf("%s.%s = %q, want %q", m, k, gv[k], v)
							return
						}
					}
				}
			}
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		t.Error(err)
	}
}
