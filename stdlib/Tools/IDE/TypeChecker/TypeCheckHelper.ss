// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	binaryPath: "Z:/Software/TypeChecker/suneidotypes.exe"

	getProps(key, def = #())
		{
		if not Suneido.Member?(#TypeCheckProperties)
			Suneido.TypeCheckProperties = Object()

		if key is ''
			return Suneido.TypeCheckProperties

		return Suneido.TypeCheckProperties.GetDefault(key, def)
		}

	setProps(key, val)
		{
		if not Suneido.Member?(#TypeCheckProperties)
			Suneido.TypeCheckProperties = Object()

		Suneido.TypeCheckProperties[key] = val
		}

	BinaryPath()
		{
		return .getProps('').GetInit(#BinaryPath, .binaryPath)
		}

	// this is where the etatest1 continuous quality checks server will find binary
	ServerBinaryPath()
		{
		return "//server/d//Software/TypeChecker/suneidotypes.exe"
		}

	SetBinaryPath(path)
		{
		if path is .BinaryPath()
			return

		.setProps(#BinaryPath, path)
		.StopServer()
		}

	BinaryExists?()
		{
		if not FileExists?(.BinaryPath())
			{
			SuneidoLog("suneidotypes binary does not exist at " $ .BinaryPath())
			return false
			}
		return true
		}

	Policy()
		{
		.getProps('').GetInit(#Policy, TypeCheckerPolicy())
		}

	SetPolicy(policy)
		{
		Assert(Object?(policy))
		.setProps(#Policy, policy)
		}

	// if no long running server then spawn
	// else reuse existing server
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
		request = Object(:method, arguments: orderedSrc, references: refs, config: policy)
		request= Json.Encode(request)
		return .send(request, :restartOnError?)
		}

	send(request, restartOnError? = true)
		{
		try
			result = .callCheck(.Server(), request)
		catch (e) //respawn and retry
			{
			if not restartOnError?
				throw e
			.StopServer()
			result = .callCheck(.Server(), request)
			}
		return Json.Decode(result.GetDefault(#content, ''))
		}

	// only throws on transport failure not on a non 2XX http code
	callCheck(port, request)
		{
		return Http('POST', 'http://127.0.0.1:' $ port $ '/check', request,
			header: Object('Content-Type': 'application/json'), timeout: 120)
		}

	Server()
		{
		if .getProps('').Member?(#Server)
			return .getProps(#Server).port

		if not .BinaryExists?()
			throw "Binary unavailable"

		if not Sys.Windows?()
			Spawn(P.WAIT, 'chmod', '+x', .BinaryPath())
		rp = RunPiped(.BinaryPath() $ ' -serve')
		.setProps(#Server, Object(:rp, port: .readReadyPort(rp)))
		return .getProps(#Server).port
		}

	StopServer()
		{
		if not .getProps('').Member?(#Server)
			return

		server = .getProps(#Server)
		Suneido.TypeCheckProperties.Delete(#Server)

		if Number?(server.GetDefault(#port, false))
			try
				Http.Post('http://127.0.0.1:' $ server.port $ '/shutdown', '', timeout: 5)

		try
			server.rp.Close()
		}

	readReadyPort(rp)
		{
		line = rp.Readline() // binary prints "READY port=NNNN" once listening
		if not String?(line) or not line.Has?('READY port=')
			throw 'type checker did not report READY: ' $ Display(line)

		return Number(line.AfterFirst('=').Trim())
		}

	// if skipLineageOrLibName is false, then we build the lineage via TypeCheckerLineage
	// if it is a string then we assume it to be a valid lib name like stdlib, axonlib...
	// and do a direct db query for that record
	OrderedSrc(className, skipLineageOrLibName = false, src = false)
		{
		if String?(skipLineageOrLibName)
			return Object(Object(name: className,
				src: src isnt false
					? src
					: Query1(skipLineageOrLibName, name: className, group: -1).text))

		chains = TypeCheckerLineage(Object(className))
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
		errors.Map!({ .formatDiagnostic('ERROR', it, library) })
		warnings.Map!({ .formatDiagnostic('WARNING', it, library) })

		return errors, warnings
		}

	formatDiagnostic(kind, d, library)
		{
		if library is false
			// control's parseDiagnosticLine strips "KIND: " then reads Class.Method:Line
			return String(kind) $ ": " $ String(d.class) $ "." $ String(d.method) $
				":" $ String(d.line) $ " " $ String(d.msg)

		// lead with lib:Record:line for the goto regex; keep "KIND:" after it so
		// Addon_highlight_warnings still colors the line
		return String(library) $ ":" $ String(d.class) $ ":" $ String(d.line) $
			" " $ String(kind) $ ": " $ String(d.method) $ " " $ String(d.msg)
		}
	}
