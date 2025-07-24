// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
ChooseListControl
	{
	Name: 'State'
	New(mandatory = false)
		{
		super(States, :mandatory, status: 'A two letter state abbreviation')
		}
	}