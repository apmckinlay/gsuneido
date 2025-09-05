// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	threadName: #LibLocateList
	Start()
		{
		if Suneido.GetDefault('LibLocateList_Thread', false)
			return

		Suneido.LibLocateList_Thread = true
		SujsAdapter.CallOnRenderBackend('RegisterBeforeDisconnectFn', .StopThread)
		Thread(.Run)
		}

	StopThread()
		{
		Suneido.LibLocateList_Thread = false
		}

	Run()
		{
		// After the thread has been started, name it, and ensure that nothing else
		// has started a similar thread. Then start checking / updating the list
		Thread.Name(.threadName)
		Database.SessionId("LibLocateList")

		Suneido.LibLocateList_ForceRun = false
		nextRun = Date()
		forever
			{
			if Suneido.LibLocateList_Thread isnt true
				return

			if .allowUpdate?(nextRun)
				{
				Suneido.LibLocateList_ForceRun = false
				.updateIfNeeded()
				nextRun = Date().Plus(seconds: 15)
				}
			Thread.Sleep(2.SecondsInMs())
			}
		}

	allowUpdate?(nextRun)
		{
		return not TestRunner.RunningTests?() and
			(Suneido.LibLocateList_ForceRun or Date() > nextRun)

		}

	libLocateOb: #(libs: (), libNums: (), list: ())
	updateIfNeeded()
		{
		if false is libs = .libs()
			return
		info = Suneido.GetDefault(#LibLocate, .libLocateOb)
		libNums = .libNums(libs)

		libs.Add('Builtin', at: 0)
		if info.libs isnt libs or info.libNums isnt libNums or info.list.Empty?()
			Suneido.LibLocate = Object(:libs, :libNums,	list: .getList(libs, info.list))
		}

	libs()
		{
		try
			return Libraries().MergeUnion(LibraryTables()) // put Use'd Libraries first
		catch(unused, '*socket connection timeout')
			return false
		}

	// use the max num in each library to detect when records have been added
	libNums(libs)
		{
		x = Object()
		for lib in libs
			try
				x[lib] = QueryMax(lib, 'num', 0)
			catch (unused, "*nonexistent table")
				{}
		return x.Set_readonly()
		}

	getList(libs, oldList)
		{
		list = Object()
		padding = libs.Size().IntDigits()
		for (i = 0; i < libs.Size(); ++i)
			if .forceRunObserver({ .processLib(libs, i, list, padding) })
				return oldList
		return list.Sort!()
		}

	processLib(libs, i, list, padding)
		{
		lib = libs[i]
		li = i.Pad(padding)
		names = .getNames(i, lib)
		for x in names
			if .forceRunObserver({ .processName(x, lib, li, list) })
				return
		}

	getNames(i, lib)
		{
		try
			return i is 0  // should be Builtin
				? BuiltinNames()
				: QueryList(lib $ ' where group = -1', 'name')
		catch (unused, "*nonexistent table")
			return  #()
		}

	forceRunObserver(block)
		{
		forceRun = .forceRun?()
		if not forceRun
			block()
		return forceRun
		}

	forceRun?()
		{
		return Suneido.GetDefault(#LibLocateList_ForceRun, false)
		}

	processName(x, lib, li, list)
		{
		if '' isnt msg = .validRecord(x, lib)
			{
			.printError(msg)
			return
			}
		// convert e.g. My_Name to myname=My_Name:03
		// library is represented as index so it sorts properly
		// use '+' and '%' because they are < numbers & letters
		list.Add(x.Tr('_').Lower() $ '+' $ x $ '%' $ li)
		caps = x.Replace(`([A-Z])[A-Z]+`, `\1`).Tr('_a-z0-9?!')
		if caps.Size() > 1
			list.Add(caps.Lower() $ '+' $ x $ '%' $ li)
		}
	validRecord(x, lib)
		{
		return not String?(x)
			? "Error in LibLocateList.getList() - Invalid record: " $ lib $ ':' $ x
			: ""
		}
	printError(msg)
		{
		Print(msg)
		}

	GetMatches(prefix, max_matches = 30, justName = false)
		{
		info = Suneido.GetDefault(#LibLocate, .libLocateOb)
		return .getMatches(info, prefix, max_matches, justName)
		}

	getMatches(info, prefix, max_matches = 30, justName = false)
		{
		list = info.list
		prefixOrig = prefix.Tr('_')
		prefix = prefixOrig.Lower()
		from = list.BinarySearch(prefix)
		to = Min(from + max_matches, list.BinarySearch(prefix.RightTrim() $ '~'))
		// convert e.g. 'name=Name:03' to 'Name - lib' or 'Name - (lib)'
		matches = list[from .. to]
		matches.Map!({ it.AfterFirst('+') })
		matches.SortWith!(#Lower)
		if justName
			matches.Map!({ it.BeforeFirst('%') }).Unique!()
		else
			matches.Unique!().Map!({ it.Replace('%\d+$',
				{ ' - ' $ .indexToLib(info.libs, Number(it[1..])) }) })
		return .moveExactMatchesFront(matches, prefix, prefixOrig)
		}

	moveExactMatchesFront(matchList, prefix, prefixOrig)
		{
		exactMatch = false
		matches = Object()
		nonMatches = Object()
		for x in matchList
			{
			if x.BeforeFirst(' - ') is prefixOrig
				{
				exactMatch = x
				continue
				}

			xLower = x.Lower().Tr('_')
			if xLower.Prefix?(prefix $ ' - ') or xLower is prefix
				matches.Add(x)
			else
				nonMatches.Add(x)
			}
		matches.Append(nonMatches)
		return exactMatch is false ? matches : matches.Add(exactMatch, at: 0)
		}

	indexToLib(libs, i)
		{
		lib = libs[i]
		if not Libraries().Has?(lib) and lib isnt 'Builtin'
			lib = '(' $ lib $ ')'
		return lib
		}

	ForceRun()
		{
		Suneido.LibLocate = #(libs: (), libNums: (), list: ())
		Suneido.LibLocateList_ForceRun = true
		}
	}
