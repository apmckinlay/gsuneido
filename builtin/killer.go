// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package builtin

import . "github.com/apmckinlay/gsuneido/core"

// killer is used by Defer, Delay, and SocketServer
type killer struct {
	ValueBase[*killer]
	kill func()
}

var _ Value = (*killer)(nil)

func (k *killer) Equal(other any) bool {
	return k == other
}

func (k *killer) SetConcurrent() {
	// need to allow this because of saving in concurrent places
	// still shouldn't be calling it from other threads
}

func (k *killer) Lookup(_ *Thread, method string) Value {
	return killerMethods[method]
}

var killerMethods = methods("killer")

var _ = method(killer_Kill, "()")

func killer_Kill(this Value) Value {
	this.(*killer).kill()
	return nil
}
