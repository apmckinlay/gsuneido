// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	// LEGACY-TYPECHECK-BINARY
	// the in-process checker is used as soon as the exe provides it, so a dev
	// switches over just by updating their exe
	// to retire the binary, delete the TypeCheckerLegacyBinary record - nothing
	// here needs editing, the name is looked up indirectly so its absence is
	// only reachable on an exe that has no TypeChecker
	transport()
		{
		return (.useBuiltin?)() ? TypeCheckerBuiltin : Global(#TypeCheckerLegacyBinary)
		}

	// Transport is recorded here rather than in transport() so it is written
	// once per process instead of on every call
	useBuiltin?: MemoizeSingle
		{
		Func()
			{
			builtin? = BuiltinNames().BinarySearch?(#TypeChecker)
			if not Suneido.Member?(#TypeCheckProperties)
				Suneido.TypeCheckProperties = Object()
			Suneido.TypeCheckProperties.Transport = builtin? ? #builtin : #binary
			return builtin?
			}
		}

	getProps(key, def = #())
		{
		if not Suneido.Member?(#TypeCheckProperties)
			Suneido.TypeCheckProperties = Object()

		if key is ""
			return Suneido.TypeCheckProperties

		return Suneido.TypeCheckProperties.GetDefault(key, def)
		}

	setProps(key, val)
		{
		if not Suneido.Member?(#TypeCheckProperties)
			Suneido.TypeCheckProperties = Object()

		Suneido.TypeCheckProperties[key] = val
		}

	// LEGACY-TYPECHECK-BINARY
	// false once there is no binary, which hides the picker
	BinaryPath()
		{
		return (.transport()).Path()
		}

	// LEGACY-TYPECHECK-BINARY
	SetBinaryPath(path)
		{
		(.transport()).SetPath(path)
		}

	BinaryExists?()
		{
		// never type check a live system
		if not (.liveSystem?)()
			return (.transport()).Available?()

		return false
		}

	liveSystem?: MemoizeSingle
		{
		Func()
			{
			return ServerEval("Thread.List").Any?({ it.Has?("scheduler extra process") })
			}
		}

	Policy()
		{
		.getProps("").GetInit(#Policy, TypeCheckerPolicy())
		}

	SetPolicy(policy)
		{
		Assert(Object?(policy))
		.setProps(#Policy, policy)
		}

	// we skip lineage checking in two particular cases if the lib is not loaded
	// or if we are passed in a function
	Run(className, method, policy = false, references? = true,
		skipLineageOrLibName = false, restartOnError? = true, src = false)
		{
		sources = .OrderedSrc(className, :skipLineageOrLibName, :src)
		return .Check(sources, method, policy, references?, :restartOnError?)
		}

	Check(orderedSrc, method, policy = false, references? = true, restartOnError? = true)
		{
		if policy is false
			policy = TypeCheckerPolicy()
		refs = references? ? .references(orderedSrc) : #()
		return (.transport()).Check(method, orderedSrc, refs, policy, :restartOnError?)
		}

	// LEGACY-TYPECHECK-BINARY
	Server()
		{
		return (.transport()).Start()
		}

	// LEGACY-TYPECHECK-BINARY
	StopServer()
		{
		(.transport()).Stop()
		}

	// if skipLineageOrLibName is false, then we build the lineage via TypeCheckerLineage
	// if it is a string then we assume it to be a valid lib name like stdlib, axonlib...
	// and do a direct db query for that record
	OrderedSrc(className, skipLineageOrLibName = false, src = false)
		{
		if String?(skipLineageOrLibName)
			return [Object(name: className,
					src: src isnt false
						? src
						: Query1(skipLineageOrLibName, name: className, group: -1).text)]

		chains = TypeCheckerLineage([className])
		chain = chains.GetDefault(className, false)
		if chain is false
			throw 'TypeChecker: "' $ className $ '" is not a loadable class'

		if src is false
			return chain.Copy()

		result = Object()
		for e in chain
			result.Add(e.name is className ? Object(name: e.name, :src) : e)
		return result
		}

	references(orderedSrc)
		{
		seen = Object()
		depNames = Object()
		for e in orderedSrc
			seen[e.name] = true
		for e in orderedSrc
			{
			refs = TypeCheckerRefs(e.src)
			for name in refs.constructed.Members()
				if not seen.Member?(name)
					depNames[name] = true

			for name in refs.called.Members()
				if not seen.Member?(name)
					depNames[name] = true
			}
		return .chainReferences(seen, depNames)
		}

	chainReferences(seen, depNames)
		{
		references = Object()
		chains = TypeCheckerLineage(depNames.Members())
		for depName in depNames.Members()
			{
			chain = chains.GetDefault(depName, false)
			if chain is false
				continue
			for e in chain // base->leaf; a shared base lands once
				if not seen.Member?(e.name)
					{
					references.Add(e)
					seen[e.name] = true
					}
			}
		return references
		}

	FormatDiagnostics(diagnostics, library = false)
		{
		if diagnostics is false or not Object?(diagnostics)
			return #(), #()

		// checker emits base-first, line-descending within each class;
		// reversing the whole list yields leaf-first, line-ascending
		errors = diagnostics.GetDefault(#errors, Object()).Reverse!()
		warnings = diagnostics.GetDefault(#warnings, Object()).Reverse!()
		errors.Map!({ .formatDiagnostic(#ERROR, it, library) })
		warnings.Map!({ .formatDiagnostic(#WARNING, it, library) })

		return errors, warnings
		}

	formatDiagnostic(kind, d, library)
		{
		if library is false
			// control's parseDiagnosticLine strips "KIND: " then reads Class.Method:Line
			return String(kind) $ ": " $ String(d.class) $ '.' $ String(d.method) $ ':' $
				String(d.line) $ ' ' $ String(d.msg)

		// lead with lib:Record:line for the goto regex; keep "KIND:" after it so
		// Addon_highlight_warnings still colors the line
		return String(library) $ ':' $ String(d.class) $ ':' $ String(d.line) $ ' ' $
			String(kind) $ ": " $ String(d.method) $ ' ' $ String(d.msg)
		}
	}
