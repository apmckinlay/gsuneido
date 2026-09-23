// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
// BuiltDate > 20260819
class
	{
	Policy()
		{
		.getProps("").GetInit(#Policy, TypeCheckerPolicy())
		}

	SetPolicy(policy)
		{
		Assert(Object?(policy))
		.setProps(#Policy, policy)
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

	Run(className, method, policy = false, references? = true,
		skipLineageOrLibName = false, src = false)
		{
		sources = .OrderedSrc(className, :skipLineageOrLibName, :src)
		return .Check(sources, method, policy, references?)
		}

	// the builtin checker only exists in exes built after this date
	TypeCheckerAvailable?()
		{
		return BuiltDate() > #20260819
		}

	Check(orderedSrc, method, policy = false, references? = true)
		{
		if not .TypeCheckerAvailable?()
			return Object(diagnostics: Object(errors: #(), warnings: #()), result: false)
		if policy is false
			policy = TypeCheckerPolicy()
		refs = references? ? .references(orderedSrc) : #()
		TypeCheckerSignatures()
		if method is TypeCheckerMethods.Infer
			return TypeChecker.Infer(orderedSrc, refs, policy)
		if method is TypeCheckerMethods.Annotate
			return TypeChecker.Annotate(orderedSrc, refs, policy)
		throw "TypeChecker: unknown method: " $ Display(method)
		}

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

	FormatDiagnostics(diagnostics, library = false, record = false)
		{
		if diagnostics is false or not Object?(diagnostics)
			return #(), #()

		// checker emits base-first, line-descending within each class;
		// reversing the whole list yields leaf-first, line-ascending
		errors = .forRecord(diagnostics.GetDefault(#errors, #()), record).Reverse!()
		warnings = .forRecord(diagnostics.GetDefault(#warnings, #()), record).Reverse!()
		errors.Map!({ .formatDiagnostic(#ERROR, it, library) })
		warnings.Map!({ .formatDiagnostic(#WARNING, it, library) })

		return errors, warnings
		}

	// compare the class exactly - filtering the formatted lines on a
	// library $ ':' $ record prefix also matches record__webgui and record_Test
	// copying is the point of the false branch too: Reverse!/Map! below would
	// otherwise reverse the caller's diagnostics and replace them with strings
	forRecord(diags, record)
		{
		return record is false ? diags.Copy() : diags.Filter({ it.class is record })
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
