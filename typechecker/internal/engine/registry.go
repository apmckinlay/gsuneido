// Copyright Suneido Software Corp. All rights reserved.
// Governed by the MIT license found in the LICENSE file.

package engine

type RefSource struct {
	Name string
	Src  string
}

type Registries struct {
	Returns   map[string]map[string]DynType
	Seeds     map[string]map[string]DynType
	Members   map[string]map[string]DynType
	Summaries map[string]map[string]ReturnSummary
	Bases     map[string]string
	Sigs      map[string]map[string]*Signature
}

func (r Registries) Seed(env *TypeEnv) {
	env.Classes = r.Returns
	env.PreCtorClasses = r.Seeds
	env.ClassSummaries = r.Summaries
	env.ClassMembers = r.Members
	env.ClassBases = r.Bases
	env.ClassMethodSigs = r.Sigs
}

// refs must arrive base-first; a ref that panics is skipped silently
func BuildReferenceRegistry(refs []RefSource) Registries {
	r := Registries{
		Returns:   make(map[string]map[string]DynType, len(refs)),
		Seeds:     make(map[string]map[string]DynType, len(refs)),
		Members:   make(map[string]map[string]DynType, len(refs)),
		Summaries: make(map[string]map[string]ReturnSummary, len(refs)),
		Bases:     make(map[string]string, len(refs)),
		Sigs:      make(map[string]map[string]*Signature, len(refs)),
	}
	pctx := NewPassCtx()
	type refState struct {
		cls *ClassObject
		env TypeEnv
	}
	states := make([]refState, 0, len(refs))
	for _, ref := range refs {
		func() {
			defer func() { _ = recover() }()
			cls := NewClassObject(ref.Name, ParseClass(ref.Src))
			env := NewTypeEnv()
			r.Seed(&env)
			var parentReturns map[string]DynType
			if cls.Base != "" {
				parentReturns = r.Returns[cls.Base]
			}
			RunPipeline(cls, env, pctx, parentReturns)
			r.Sigs[ref.Name] = ClassMethodSignatures(cls)

			env = env.WithClass(cls, r.Sigs[ref.Name])
			states = append(states, refState{cls, env})
			if cls.IsFunction {
				r.Returns[ref.Name] = map[string]DynType{
					"CallClass": env.LookupReturn(ref.Name),
				}
				if s := env.SnapshotSummaries(); len(s) > 0 {
					if sum, ok := s[ref.Name]; ok {
						r.Summaries[ref.Name] = map[string]ReturnSummary{
							"CallClass": sum,
						}
					}
				}
				return
			}
			r.Returns[ref.Name] = env.SnapshotReturns()
			r.Members[ref.Name] = ClassStaticMemberTypes(cls)
			r.Bases[ref.Name] = cls.Base
			if seed := env.SnapshotPreCtorReturns(); len(seed) > 0 {
				r.Seeds[ref.Name] = seed
			}
			if s := env.SnapshotSummaries(); len(s) > 0 {
				r.Summaries[ref.Name] = s
			}
		}()
	}

	conflicts := func() int {
		n := 0
		for i := range states {
			n += len(states[i].env.ReqConflicts)
		}
		return n
	}
	for range maxFixpointPasses {
		changed := false
		before := conflicts()
		for i := range states {
			s := &states[i]
			func() {
				defer func() { _ = recover() }()
				CallsiteResolutionPass(s.cls, s.env, pctx)
				if RequirementPass(s.cls, s.env, pctx) {
					changed = true
				}
			}()
		}
		if conflicts() > before {
			for _, sigs := range r.Sigs {
				clearInferred(sigs)
			}
		}
		if !changed {
			return r
		}
	}
	return r
}
