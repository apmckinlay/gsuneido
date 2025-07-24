// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Controls(@unused)
		{
		return #(Horz
			(Static Query)
			Skip
			(Editor name: Source))
		}

	Source(source, unused)
		{
		return Object(query: source.Get(), exclude: #())
		}
	}