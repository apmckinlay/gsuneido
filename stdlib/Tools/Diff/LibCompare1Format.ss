// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
DiffFormat
	{
	New(t1, t2, name, path = false)
		{
		super(@.diff(t1, t2, name, path))
		}
	diff(t1, t2, name, path)
		{
		if (path isnt false)
			q_string = " where name is '" $ name $ "' and path is " $ Display(path)
		else
			q_string = " where name is '" $ name $ "' and group is -1"
		lib_record1 = Query1(t1 $ q_string)
		lib_record2 = Query1(t2 $ q_string)
		Assert(lib_record1 isnt false and lib_record2 isnt false)
		return Object(lib_record1.text, lib_record2.text, name)
		}
	}
