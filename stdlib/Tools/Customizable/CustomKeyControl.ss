// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
KeyControl
	{
	CustomizableOptions: #(mandatory, hidden, readonly, tabover)
	New(@args) // expecting customField
		{
		super(@.Setup_CustomKeyControl(args))
		}

	Setup_CustomKeyControl(args)
		{
		args = args.Copy()
		field = args.customField
		args.field = 'name'
		args.query = field $ '_table'
		args.columns = #("name", "desc")
		args.access = 'Access_' $ field
		args.customizeQueryCols = true
		return args
		}

	ValidData?(@args)
		{
		field = args.customField
		args = args.Copy()
		args.query = field $ '_table'
		args.field = 'name'
		return super.ValidData?(@args)
		}
	}
