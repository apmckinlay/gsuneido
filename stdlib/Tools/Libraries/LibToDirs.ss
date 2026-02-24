// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
// export a library to files with the same folder structure
// starts by creating a directory in dest with the same name as the library
class
	{
	CallClass(lib, dest = '.')
		{
		if DirExists?(path = Paths.Combine(dest, lib))
			throw "already exists: " $ Display(path)
		tm = TreeModel(lib)
		.export(tm, 0, dest, lib)
		}
	export(tm, parent, dest, name)
		{
		dest = Paths.Combine(dest, name)
		EnsureDir(dest)
		for c in tm.Children(parent)
			if c.group
				.export(tm, c.num, dest, c.name)
			else if not c.name.Suffix?("_alpha") and not c.name.Suffix?("_trial")
				PutFile(Paths.Combine(dest, c.name.Tr('?', 'Q')) $ ".ss", c.text)
		}
	}