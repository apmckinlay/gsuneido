// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
StaticControl
	{
	Name: Heading
	// public members so MainField can use them
	Weight: 'bold'
	New(text, size = '+2')
		{
		super(text, size: .defaultSize(size), weight: .Weight, color: CLR.Highlight)
		}

	Size: false // overriden in Heading1, Heading2, Heading3
	defaultSize(size)
		{
		return .Size isnt false ? .Size : size
		}
	}
