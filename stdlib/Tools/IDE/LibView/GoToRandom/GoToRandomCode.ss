// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
function (source)
	{
	do
		rec = RandomLibraryRecord()
		while rec.text.LineCount() < 20
	lib = rec.lib
	name = rec.name
	n = rec.text.LineCount()
	line = Random(n)
	libview = source.Ctrl
	path = LibHelp.NamePath(lib, name)
	libview.GotoPathLine(path, line)
	}
