// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
function (source)
	{
	do
		rec = RandomLibraryRecord()
		while rec.text.LineCount() < 20 /*= minLines */
	libview = source.Ctrl
	path = LibHelp.NamePath(rec.lib, rec.name)
	libview.GotoPathLine(path, Random(rec.text.LineCount()))
	}
