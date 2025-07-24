// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Dump Libraries'
	CallClass(hwnd = 0)
		{
		if false isnt libs = OkCancel(Object(this), .Title, hwnd)
			.DumpLibraries(libs)
		}
	Controls()
		{
		ob = Object('Vert')
		ob.Add(Object('TwoList', LibraryTables(), name: 'dumpLibs'))
		ob.Add('Skip', #(HorzEqual
			(Button, 'All') Skip (Button, 'In Use') Skip (Button None)))
		return Object('Record', ob)
		}
	OK()
		{
		return .Data.FindControl('dumpLibs').GetNewList()
		}
	DumpLibraries(libs)
		{
		for lib in libs.Copy()
			Database.Dump(lib)
		}
	On_All()
		{
		.setList(LibraryTables())
		}
	On_In_Use()
		{
		.setList(Libraries())
		}
	On_None()
		{
		.setList(#())
		}
	setList(libs)
		{
		ctrl = .Data.FindControl('dumpLibs')
		if ctrl is false
			return

		ctrl.AllBack()

		if libs.Empty?()
			return

		ctrl.Set(libs.Join(','))
		}
	}