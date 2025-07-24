// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
TestObserverText
	{
	New(file)
		{
		.f = File(file, 'w')
		}
	Output(s)
		{
		.f.Writeline(s)
		.f.Flush()
		}
	After(@args)
		{
		super.After(@args)
		.f.Close()
		}
	}