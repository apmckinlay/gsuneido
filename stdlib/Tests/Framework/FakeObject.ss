// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(@methods)
		{
		.methods = methods
		}
	Default(@args)
		{
		method = args.PopFirst()
		if not .methods.Member?(method)
			throw "FakeObject unexpected call: " $
				'.' $ method $ Display(args)[1..]
		m = .methods[method]
		return Function?(m) ? m(@args) : m
		}
	}