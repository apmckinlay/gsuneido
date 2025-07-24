// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
ChooseListControl
	{
	Name: 'Prov'
	New(mandatory = false)
		{
		super(Provinces, :mandatory, status: 'A two letter province abbreviation')
		}
	}