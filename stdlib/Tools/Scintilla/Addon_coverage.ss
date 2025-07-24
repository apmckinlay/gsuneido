// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	marginId: 6
	New(@args)
		{
		super(@args)
		.sub = PubSub.Subscribe('LibraryRecordChange', .reset)
		}
	Init()
		{
		.covered = Object()
		.covered.Marker = .MarkerIdx(level: .coveredLevel)
		.covered.Indicator = .IndicatorIdx(level: .coveredLevel)
		.nonCovered = Object()
		.nonCovered.Marker = .MarkerIdx(level: .nonCoveredLevel)
		.nonCovered.Indicator = .IndicatorIdx(level: .nonCoveredLevel)

		SendMessage(.Hwnd, SCI.SETMARGINS, .marginId + 1, 0)
		.SetMarginTypeN(.marginId, SC.MARGIN_RTEXT)

		defaultFore = .GetSchemeColor('defaultFore')
		.coveredStyle = SC.STYLE_LASTPREDEFINED + 1
		.uncoveredStyle = SC.STYLE_LASTPREDEFINED + 2
		.DefineStyle(.coveredStyle, defaultFore,
			back: .dark ? RGB(71,87,45) : RGB(225,240,225) /*= dark green */)
		.DefineStyle(.uncoveredStyle, defaultFore,
			back: .dark ? RGB(123,0,11) : RGB(255,225,225) /*= red */)
		}

	getter_dark()
		{
		return .dark = .GetSchemeColor('defaultBack') < 0xdddddd
		}

	ContextMenu()
		{
		return #('Remove Coverage Markers\tF6')
		}

	On_Remove_Coverage_Markers()
		{
		.reset()
		.UpdateUI()
		.MarkersChanged() // update overview bar
		}

	Modified(scn /*unused*/)
		{
		if not .SetMethodModifying?()
			.reset()
		}

	coveredLevel: 		85
	nonCoveredLevel: 	86
	Styling()
		{
		dark = .GetSchemeColor('defaultBack') < 0xdddddd
		box = dark ? INDIC.HIDDEN : INDIC.FULLBOX
		mark = dark ? SC.MARK_BACKGROUND : SC.MARK_VLINE
		return [[level: .coveredLevel,
				marker: [mark, back: dark ? RGB(71,87,45) : CLR.green], /*= dark green */
				indicator: [box, fore: CLR.green, back: CLR.green]],
			[level: .nonCoveredLevel,
				marker: [mark, back: dark ? RGB(123,0,11) : CLR.red], /*= dark red */
				indicator: [box, fore: CLR.red, back: CLR.red]]]
		}

	On_BeforeAllTests()
		{
		.Send('Save')
		.Send('ForceEntabAll')
		.reset()
		CoverageEnable(true)
		for item in .Send('GetTabsPaths', all?:, skipFolder?:)
			{
			lib = item.BeforeFirst('/')
			name = item.AfterLast('/')
			.startCoverageAndTest(lib, name)
			}
		}

	startCoverageAndTest(lib, name)
		{
		if name.Suffix?('Test')
			{
			.startCoverage(lib, name)
			if false isnt recordName = .getNonTest(name)
				.startCoverage(lib, recordName)
			}
		else
			{
			.startCoverage(lib, recordName = name)
			if false isnt recordTestName = .getTestName(name)
				.startCoverage(lib, recordTestName)
			if false isnt recordTestName = .getTestName2(name)
				.startCoverage(lib, recordTestName)
			}
		}

	startCoverage(lib, name)
		{
		id = lib $ ':' $ name
		if Suneido.coverage.Member?(id)
			return

		if .skip?(lib, name)
			return

		Unload(name)
		Global(name).StartCoverage(count:)
		Suneido.coverage[id] = Object().Set_default(0)
		}

	skip?(lib, name)
		{

		if name is '' // for stdlib:Test record
			return true

		tables = Libraries()
		if false is libPos = tables.Find(lib)
			return true

		for lib in tables[libPos+1..]
			if false isnt Query1(lib, :name, group: -1)
				return true

		if .verifyName(name) is false
			return true

		return false
		}

	getNonTest(name)
		{
		if not name.Suffix?('Test')
			return false
		nonTest = name.RemoveSuffix("Test").RemoveSuffix('_')
		if false isnt found = .verifyName(nonTest)
			return found
		return .verifyName(nonTest $ '?')
		}

	verifyName(name)
		{
		if false is fn = .global(name)
			return false
		if not Function?(fn) and not Class?(fn)
			return false
		if Type(fn).Has?("Builtin")
			return false
		return name
		}

	global(name)
		{
		try
			{
			fn = Global(name)
			return fn
			}
		catch (unused, "can't find|error loading|Assert FAILED: Global invalid")
			{
			return false
			}
		}

	getTestName(name)
		{
		test = name $ '_Test'
		if false isnt .global(name)
			return test
		return false
		}

	getTestName2(name)
		{
		test = name $ 'Test'
		if false isnt .global(test)
			return test
		return false
		}

	On_AfterAllTests()
		{
		names = Object()
		for id in .coverage.Members()
			{
			name = id.AfterFirst(':')
			c = .coverage[id] = Global(name).StopCoverage()
			for m in c.Members().Sort!()
				.coverageSorted[id].Add(Object(m, c[m]))
			names.AddUnique(name)
			}
		CoverageEnable(false)
		// unload has to be after coverage is turned off,
		// classes like Strings could be loaded again right away
		for name in names
			Unload(name)
		.UpdateUI()
		}

	getter_coverage()
		{
		if not Suneido.Member?('coverage')
			Suneido.coverage = Object().Set_default(Object().Set_default(0))
		return Suneido.coverage
		}

	getter_coverageSorted()
		{
		if not Suneido.Member?('coverageSorted')
			Suneido.coverageSorted = Object().Set_default(Object())
		return Suneido.coverageSorted
		}

	reset()
		{
		Suneido.coverage = Object().Set_default(Object().Set_default(0))
		Suneido.coverageSorted = Object().Set_default(Object())
		}

	UpdateUI()
		{
		.clear()
		if false is m = .getId()
			return

		if not .coverageSorted.Member?(m)
			{
			.SetMarginWidthN(.marginId, 0)
			return
			}

		.SetMarginWidthN(.marginId, ScaleWithDpiFactor(48/*=width*/))
		.SetMarginBackN(.marginId, CLR.red /*.GetSchemeColor('defaultBack')*/)
		code = .Get()
		pre = 0
		preLine = -1
		for s in .coverageSorted[m]
			{
			coverStart = .rollBackBlanks(code, s[0], pre)
			end = code.Find('\n', coverStart)
			line = .LineFromPosition(coverStart)
			if preLine is line
				continue

			txt = .getMarginText(s[1])
			style = s[1] is 0 ? .nonCovered : .covered
			.MarkerAdd(line, style.Marker)
			.SetIndicator(style.Indicator, coverStart, end - coverStart)
			.MarginSetStyle(line, s[1] is 0 ? .uncoveredStyle : .coveredStyle)
			SendMessageTextIn(.Hwnd, SCI.MARGINSETTEXT, line, txt)
			pre = end
			preLine = line
			}
		}

	clear()
		{
		.ClearIndicator(.covered.Indicator)
		.MarkerDeleteAll(.covered.Marker)
		.ClearIndicator(.nonCovered.Indicator)
		.MarkerDeleteAll(.nonCovered.Marker)
		}

	getId()
		{
		try
			{
			name = .Send('CurrentName')
			lib = .Send('CurrentTable')
			}
		catch(unused, '*socket connection timeout')
			return false
		return lib $ ':' $ name
		}

	getMarginText(c)
		{
		return c >= 65535 ? '>=64k' : String(c) /*= max 65535 */
		}

	rollBackBlanks(code, start, pre)
		{
		for (; start >= 0 and start > pre; start--)
			if code[start-1] not in ('\t', ' ')
				break
		return start
		}

	Destroy()
		{
		.sub.Unsubscribe()
		}
	}
