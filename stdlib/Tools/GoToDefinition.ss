// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
function (name, lib, line)
	{
	if LibraryTables().Has?(lib)
		GotoLibView(name, libs: Object(lib), :line)
	else if BookTables().Has?(lib) // book
		{
		pos = name.FindLast('/')
		path = name[.. pos]
		name = name[pos + 1 ..]

		if false isnt page = Query1(lib, :path, :name)
			OpenBook(lib, lib $ page.path $ '/' $ page.name, bookedit?:)
		}
	}