// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
function(lib)
	{
	Database("ensure " $ lib $ " (num, parent, group, name, text)
		key(num) key(name, group) index(parent, name)")
	}

