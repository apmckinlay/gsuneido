// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.cls)
		{
		}
	Default(@args)
		{
		args[0] = .cls[args[0]]
		.Eval(@args)
		}
	Getter_(m)
		{
		if .cls.Member?(m)
			return .cls[m]
		if .cls.Method?('Getter_' $ m)
			return .Eval(.cls['Getter_' $ m])
		}
	}