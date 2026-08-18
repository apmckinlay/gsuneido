// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

func TypeInfer(name, src string, resolver ClassResolver) (*ClassObject, TypeEnv) {
	cls := NewClassObject(name, ParseClass(src))
	env := NewTypeEnv()
	pctx := NewPassCtx()

	lineage := cls.Lineage(resolver)
	parentReturns := map[string]DynType{}
	for _, c := range lineage {
		RunPipeline(c, env, pctx, parentReturns)
		parentReturns = env.SnapshotReturns()
	}
	return cls, env
}
