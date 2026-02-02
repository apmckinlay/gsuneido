// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Md_ContainerBlock
	{
	New()
		{
		.linkDefs = Object()
		}

	Continue(line, start)
		{
		return line, start
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

	AddLinkDefs(label, dest, title)
		{
		if .linkDefs.Member?(label)
			return // ignore
		.linkDefs[label] = Object(:dest, :title)
		}

	GetLinkDefs(label)
		{
		return .linkDefs.GetDefault(label, false)
		}
	}