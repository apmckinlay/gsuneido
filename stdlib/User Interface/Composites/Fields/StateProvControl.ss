// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
ChooseListControl
	{
	Name: 'StateProv'
	New(mandatory = false, hidden = false, readonly = false)
		{
		super(.GetList(), :mandatory, width: 3,
			status: 'A two letter state or province abbreviation', :hidden, :readonly)
		}
	GetList() // needed by AxonMobileBuildMessage
		{
		return StatesProvsMex
		}
	ValidationList()
		{
		return .GetList()
		}
	}