// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Singleton
	{
	New()
		{
		.spysMap = Object().Set_default(Object())
		.spysIdMap = Object()
		}

	Register(spy)
		{
		key = spy.Lib $ ":" $ spy.Name
		Assert(.spysMap[key].HasIf?({ it.Paths is spy.Paths }) is: false,
			msg: "Duplicate Spy on " $ key $ Opt('.', spy.Paths.Join('.')))
		.spysMap[key].Add(spy)
		.spysIdMap[spy.Id] = spy
		.updateOverride(key)
		}

	updateOverride(key)
		{
		lib = key.BeforeFirst(':')
		name =  key.AfterFirst(':')
		spys = .spysMap[key]

		LibraryOverride(lib, name)
		if spys.Empty?()
			return

		rec = Query1(lib, :name, group: -1)
		Assert(rec isnt: false, msg: "SpyManager can't find library record - " $ name)
		source = rec.text

		code = .useTdop?(spys)
			? .buildCodeWithTdop(source, spys)
			: .buildCodeWithClassHelp(source, spys)

		LibraryOverride(lib, name, code)
		}

	useTdop?(spys)
		{
		// need Tdop to handle nested class
		return spys.HasIf?()
			{ |spy|
			spy.Paths.Size() > 1
			}
		}

	buildCodeWithTdop(source, spys)
		{
		tdop = Tdop(source)
		astWriteMgr = AstWriteManager(source, tdop)
		writer = astWriteMgr.GetNewWriter()
		spys.Each()
			{ |spy|
			funcNode = .findFuncNode(tdop, spy.Paths)
			writer.Add(funcNode[.statementListOffset], 0, .buildInsert(spy))
			}
		return writer.ToString()
		}

	statementListOffset: 5
	memberListOffset: 4
	memberContentOffset: 2
	findFuncNode(tdop, paths)
		{
		curNode = tdop
		for path in paths
			{
			Assert(curNode.Match(TDOPTOKEN.CLASSDEF))
			find? = false
			for member in curNode[.memberListOffset].Children
				{
				// second condition is for matching private method
				if member[0].Value isnt path and not path.Suffix?('_' $ member[0].Value)
					continue
				curNode = member[.memberContentOffset]
				find? = true
				break
				}
			if find? is false
				throw "SpyOn cannot find specified method - " $ paths.Join(".")
			}
		Assert(curNode.Match(TDOPTOKEN.FUNCTIONDEF))
		return curNode
		}

	buildInsert(spy)
		{
		"\r\nres = SpyManager().Spy(" $ Display(spy.Id) $
			", " $ .buildParams(spy.Params) $ ")\r\n" $
		"if res.action is 'return' { return res.value }\r\n" $
		"if res.action is 'throw' { throw res.value }\r\n"
		}

	buildParams(params)
		{
		ClassHelp.RetrieveParamsList(params, list = Object())
		return 'Object(' $ list.Map({ ':' $ it }).Join(', ') $ ')'
		}

	buildCodeWithClassHelp(source, spys)
		{
		if spys[0].Method? is false
			{
			return .buildFuncWithSpy(source, spys[0])
			}
		methodRanges = ClassHelp.MethodRanges(source)
		replaces = Object()
		spys.Each()
			{ |spy|
			.buildReplaceMethods(spy, methodRanges, replaces, source)
			}
		replaces.Sort!(By(#from))
		code = ''
		prev = 0
		replaces.Each()
			{
			code $= source[prev..it.from] $ it.text
			prev = it.to
			}
		return code $ source[prev..]
		}

	buildReplaceMethods(spy, methodRanges, replaces, source)
		{
		method = spy.Paths[0]
		if false is methodRange = methodRanges.FindOne({
			it.name is method or method.Suffix?(it.name)
			})
			{
			code = '
				' $ method $ spy.Params $ '
				{' $ .buildInsert(spy) $ '
				super.' $ method $ spy.Params $ '
				}'
			last = methodRanges.Last()
			replaces.Add(Object(from: last.to, to: last.to, text: code))
			}
		else
			{
			methodText = source[methodRange.from..methodRange.to]
			replaces.Add(Object(from: methodRange.from, to: methodRange.to,
				text: .buildFuncWithSpy(methodText, spy)))
			}
		}

	buildFuncWithSpy(text, spy)
		{
		paramsEnd = ClassHelp.RetrieveParamsList(text, Object())
		statementsStart = text.Find('{', paramsEnd)
		return text[..statementsStart + 1] $
			.buildInsert(spy) $
			text[statementsStart + 2..]
		}

	Spy(id, locals)
		{
		Assert(.spysIdMap hasMember: id)
		return (.spysIdMap[id])(locals)
		}

	RemoveAll()
		{
		for key in .spysMap.Members()
			{
			lib = key.BeforeFirst(':')
			name =  key.AfterFirst(':')
			LibraryOverride(lib, name)
			}
		.spysMap.Delete(all:)
		.spysIdMap.Delete(all:)
		}

	RemoveOne(spy)
		{
		key = spy.Lib $ ":" $ spy.Name
		.spysMap[key].Remove(spy)
		.spysIdMap.Delete(spy.Id)
		.updateOverride(key)
		}

	Reset()
		{
		// disable Reset because we don't want ResetCaches / Singleton.ResetAll
		// to wipe out the spys (tests may call ResetCaches)
		}
	}