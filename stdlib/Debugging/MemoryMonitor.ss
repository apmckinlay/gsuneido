// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(title = false)
		{
		Window(this, exStyle: WS_EX.TOPMOST, :title)
		}
	New()
		{
		.timer = SetTimer(.WndProc.Hwnd, id: 100, ms: 1.SecondsInMs(), f: .onTimer)
		.inclres = .FindControl(#inclres)
		.inclsi = .FindControl(#inclsi)
		.arena = .FindControl(#arena)
		.dbsize = .FindControl(#dbsize)
		.tempdest = .FindControl(#tempdest)
		.transactions = .FindControl(#transactions)
		.fonts = .FindControl(#fonts)
		.suneidomems = .FindControl(#suneidomems)
		.connections = .FindControl(#connections)
		.heapob = .FindControl('heapobject')
		.threads = .FindControl(#threads)
		.resources = .FindControl(#resources)
		.suneidoinfo = .FindControl(#suneidoinfo)
		}
	Controls: (Border (Vert
		(Horz (Heading 'Memory Monitor') Skip
			(CheckBox 'Resources', name: 'inclres') Skip
			(CheckBox 'SuneidoInfo', name: 'inclsi'))
		Skip
		(Horz
			(Vert (Pair (Static 'Heap Size') (Static "999 mb" name: 'arena'))
				(Pair (Static 'Database Size') (Static "9999 mb" name: 'dbsize'))
				(Pair (Static Transactions) (Static "9999" name: 'transactions'))
				(Pair (Static TempDest) (Static "9999" name: 'tempdest'))
				(Pair (Static HwndFonts) (Static "999" name: 'fonts'))
				(Pair (Static SuneidoMembers) (Static "999" name: 'suneidomems'))
				(Pair (Static Connections) (Static "999" name: 'connections'))
				(Pair (Static HeapObjectSize) (Static "999" name: 'heapobject'))
				)
			Skip
			(Vert name: 'resources')
			(Vert name: 'suneidoinfo')
			)
		Skip
		(Editor height: 10, width: 50, name: 'threads')
		(Button "Inspect Suneido")
		), xstretch: 0, ystretch: 0)
	onTimer(@unused)
		{
		.arena.Set(ReadableSize(MemoryArena()))
		.dbsize.Set(ReadableSize(Database.CurrentSize()))
		.tempdest.Set(Database.TempDest())
		.transactions.Set(Database.Transactions().Size())
		.fonts.Set(Suneido.GetDefault('HwndFonts', #()).Size())
		.suneidomems.Set(Suneido.Members().Size())
		.connections.Set(Sys.Connections().Size())
		size = 'unknown'
		try size = ReadableSize(Suneido.GoMetric("/memory/classes/heap/objects:bytes"))
		.heapob.Set(size)
		.threads.Set(Thread.List().Join('\n'))

		if .inclres.Get() is true
			.updateResourceCounts()
		if .inclsi.Get() is true
			.updateSuneidoInfo()

		return 0
		}
	updateResourceCounts()
		{
		counts = ResourceCounts()
		.resources.RemoveAll()
		.resources.Append(#(Static 'ResourceCounts'))
		for m in counts.Members().Sort!()
			.resources.Append(
				Object('Pair', Object('Static', m), Object('Static', String(counts[m]))))
		}
	updateSuneidoInfo()
		{
		.suneidoinfo.RemoveAll()
		.suneidoinfo.Append(#(Pair, #(Static 'Suneido.Info'), #(Static '')))
		info = Suneido.Info().Copy().Remove('build_info').Remove('built')
		for m in info.Sort!()
			.suneidoinfo.Append(
				Object('Pair', Object('Static', m),
					Object('Static', String(Suneido.Info(m)))))
		}
	On_Inspect_Suneido()
		{
		Inspect(Suneido, title: 'Client Suneido')
		}
	Destroy()
		{
		KillTimer(.WndProc.Hwnd, .timer)
		ClearCallback(.onTimer)
		super.Destroy()
		}
	}
