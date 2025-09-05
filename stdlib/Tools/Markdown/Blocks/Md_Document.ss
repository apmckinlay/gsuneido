// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_ContainerBlock
	{
	New()
		{
		}

	Continue(line)
		{
		return line
		}

	GetOpenBlockItems()
		{
		list = Object()
		next = this
		do
			{
			list.Add(next)
			}
		while false isnt next = next.NextOpenBlock()
		return list
		}

	Finish()
		{
		.Close()
		}
	}