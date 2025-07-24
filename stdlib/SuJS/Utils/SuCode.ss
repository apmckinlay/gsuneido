// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
Singleton
	{
	New()
		{
		.ensure()
		.checkUpdate()
		LibUnload.AddObserver(#SuCode, .checkUpdate)
		}

	Getter_Libraries()
		{
		libs = Libraries().Reverse!().Filter({ not it.Suffix?(#webgui) })
		if Instance?(this)
			.Libraries = libs
		return libs
		}

	ensure()
		{
		Database('ensure su_bundle_records (name, lib) key (name, lib)')
		Database('ensure su_code_bundle (part, code, code_built) key (part)')
		}

	checkUpdate(name = false)
		{
		update = name isnt false
			? not QueryEmpty?('su_bundle_records', :name) or .overrides.Has?(name)
			: .checkDependencies()
		if update is true
			.clear()
		return update
		}

	// There is one known case that will not be detected. It is when there are multiple
	// records modified and only the records with smaller lib_modified dates are restored.
	// In this case, neither of the dependency list and the max lib_modified/lib_committed
	// are changed. If SuCode (Singleton) has been initialized, the LibUnload
	// observer will catch this case, so it should be rare.
	checkDependencies()
		{
		if false is buildInfo = .getCurrentBundleBuildInfo()
			return false

		newRecs = .CollectDependencies().recs.Sort!(By(#name, #lib))
		i = 0
		lib_modified = lib_committed = Date.Begin()
		.getBundleRecs()
			{ |rec|
			newRec = newRecs[i++]
			if rec.lib isnt newRec.lib or
				rec.name isnt newRec.name or
				buildInfo.lib_committed < newRec.lib_committed or
				buildInfo.lib_modified < newRec.lib_modified
				return true

			lib_modified = Max(lib_modified, newRec.lib_modified)
			lib_committed = Max(lib_committed, newRec.lib_committed)
			}
		return lib_modified isnt buildInfo.lib_modified or
			lib_committed isnt buildInfo.lib_committed
		}

	getCurrentBundleBuildInfo()
		{
		if false is codeBundle = QueryFirst('su_code_bundle remove code sort part')
			return false
		return codeBundle.code_built
		}

	getRec(name, lib = false)
		{
		mappedName = .overrides.GetDefault(name, name)
		rec = false
		if lib isnt false
			{
			rec = .queryRec(lib, mappedName)
			rec.lib = lib
			}
		else
			for lib in .Libraries
				if false isnt rec = .queryRec(lib, mappedName)
					{
					rec.lib = lib
					break
					}
		if rec isnt false
			rec.name = name
		return rec
		}

	queryRec(lib, name)
		{
		return LibraryTags.GetRecord(name, lib)
		}

	Add(lib, name)
		{
		SuneidoLog.Once('INFO: SuCode.Add is called, but it should not be',
			params: [:lib, :name], calls:)
		QueryOutputIfNew('su_bundle_records', [:lib, :name])
		.clear()
		}

	clear()
		{
		QueryDo('delete su_code_bundle')
		if Instance?(this)
			.Delete(#CodeBundle)
		}

	Getter_CodeBundle()
		{
		return .Synchronized()
			{
			if false is code = .getBundle()
				code = .buildBundle()
			.CodeBundle = code
			}
		}

	getBundle()
		{
		code = false
		QueryApply('su_code_bundle sort part')
			{
			if code is false
				code = it
			else
				code.code $= it.code
			}
		return code
		}

	codePartSize: 500000 // 500K
	buildBundle()
		{
		toDelete = Object()
		toAdd = Object()
		code = '"use strict";
window.suCodeBundle = {
'
		buildInfo = Object(
			lib_committed: Date.Begin(),
			lib_modified: Date.Begin(),
			hash: ''
			)

		cmp = By(#name, #lib)
		newRecs = .CollectDependencies().recs.Sort!(cmp)
		oldRecs = .getBundleRecs()
		oldRecs.Add([name: 'xxx', lib: 'xxx'])
		for (i = j = 0; i < newRecs.Size(); i++)
			{
			while cmp(oldRecs[j], newRecs[i])
				toDelete.Add(oldRecs[j++])

			if oldRecs[j] isnt newRecs[i]
				toAdd.Add(newRecs[i])
			else
				j++

			rec = .getRec(newRecs[i].name, lib: newRecs[i].lib)
			code $= .buildCode(rec)
			buildInfo.lib_committed = Max(buildInfo.lib_committed, rec.lib_committed)
			buildInfo.lib_modified = Max(buildInfo.lib_modified, rec.lib_modified)
			}

		while j < oldRecs.Size() - 1
			toDelete.Add(oldRecs[j++])

		code $= '};'
		buildInfo.hash = Md5(code).ToHex()
		.update(code, buildInfo, toDelete, toAdd)
		return [:code, code_built: buildInfo]
		}

	BuildCodeBundle(extraSeeds = #())
		{
		code = '"use strict";
window.suCodeBundle = {
'
		collects = .CollectDependencies(extraSeeds)
		recs = collects.recs
		lib_committed = Date.Begin()
		lib_modified = Date.Begin()

		for i in ..recs.Size()
			{
			rec = .getRec(recs[i].name, lib: recs[i].lib)
			code $= .buildCode(rec)
			lib_committed = Max(lib_committed, rec.lib_committed)
			lib_modified = Max(lib_modified, rec.lib_modified)
			}
		code $= '};'
		return Object(:code, :lib_committed, :lib_modified)
		}

	getBundleRecs(block = false)
		{
		if block is false
			return QueryAll('su_bundle_records sort name, lib')

		QueryApply('su_bundle_records sort name, lib', :block)
		}

	buildCode(rec)
		{
		fn = JsTranslate(rec.text, rec.name, rec.lib)
		return Display(rec.name) $ ': function () { return ' $ fn $ '; },\r\n'
		}

	update(code, buildInfo, toDelete, toAdd)
		{
		RetryTransaction()
			{ |t|
			toDelete.Each({ t.QueryDo('delete su_bundle_records
				where lib is ' $ Display(it.lib) $ ' and name is ' $ Display(it.name)) })
			toAdd.Each({ t.QueryOutput('su_bundle_records', it) })
			t.QueryDo('delete su_code_bundle')
			part = 0
			for (i = 0; i < code.Size(); i += .codePartSize)
				{
				t.QueryOutput('su_code_bundle',
					[code: code[i::.codePartSize], code_built: buildInfo, :part])
				part++
				}
			}
		}

	overrides: #(
		SuneidoLog: SuOverride_SuneidoLog,
		TranslateLanguage: SuOverride_TranslateLanguage,
		ErrorLog: SuOverride_ErrorLog,
		Matcher_is: SuOverride_Matcher_is,

		Test: SuOverride_Test
		TestRunner: SuOverride_TestRunner
		)

	seeds: #(
		// extended type methods
		'HTMLElements', 'Elements', 'Nodes', 'Documents', 'Windows',
		'Dates', 'Numbers', 'Objects', 'Records', 'Sequences', 'Strings'
		// records referenced in strings
		'PrintCancelDoc', 'PrintEndDoc', 'PrintEndPage', 'PrintManager',
		'DoTaskWithPauseClient',
		"SuCanvasArc", "SuCanvasEllipse", "SuCanvasGroup", "SuCanvasImage",
			"SuCanvasLine", "SuCanvasRect", "SuCanvasRoundRect", "SuCanvasText",
		"SuDrawClickTracker", "SuDrawLineTracker", "SuDrawRectTracker",
			"SuDrawSelectTracker"
		"SuClipboardPasteString", "SuClipboardWriteHtml", "SuClipboardWriteString",
		"SuFlashWindow", "SuInitClient", "SuClearFocus", "SuJsExecute", "SuOpenBook",
		"SuSetGuiFont", "SuShutdown", "SuTaskbarUpdate",
		"SuDebugger"
		)
	skips: #('Print') // builtin in suneido.js
	CollectDependencies(extraSeeds = #())
		{
		recs = Object()
		builtins = Object()
		deps = Object()
		stack = Object()
		for seed in .initSeeds(extraSeeds)
			{
			stack.Add(seed)
			deps[seed.name] = #()
			}
		while stack.Size() > 0
			{
			cur = stack.PopLast()
			if String?(cur)
				builtins.Add(cur)
			else
				{
				recs.Add(cur.Project(#name, #lib, #lib_committed, #lib_modified))
				if .skips.Has?(cur.name)
					continue

				refs = Object().Set_default(0)
				ast = Suneido.Parse(cur.text)
				Qc_globalRefs.Traverse(ast, refs, #(), includeConstant?:,
					extraFn: .extraScan)
				refNames = refs.Members().Filter({ not deps.Member?(it) })
				if refNames.NotEmpty?()
					{
					paths = deps[cur.name].Copy().Add(cur.name)
					for name in refNames
						{
						refRec = .getRec(name)
						stack.Add(refRec is false ? name : refRec)
						deps[name] = paths
						}
					}
				}
			}
		return [:recs, :deps, :builtins]
		}

	initSeeds(extraSeeds)
		{
		seeds = Object()
		for lib in .Libraries.Copy().Reverse!()
			QueryApply(lib $ ' where group is -1 and name.Suffix?("Component")')
				{
				seeds[it.name] = Object(name: it.name, :lib,
					lib_committed: it.lib_committed,
					lib_modified: it.lib_modified,
					text: it.text)
				}
		for seed in .seeds
			seeds[seed] = .getRec(seed)
		for seed in extraSeeds
			seeds[seed] = .getRec(seed)
		return seeds.Values()
		}

	extraScan(ast, refs)
		{
		if ast.type isnt #Call or ast.func.type isnt #Ident or ast.func.name isnt 'Assert'
			return
		if ast.size < 2 or ast[1].name in (false, #msg)
			return
		++refs['Matcher_' $ ast[1].name]
		}

	Reset()
		{
		LibUnload.RemoveObserver(#SuCode)
		super.Reset()
		}
	}
