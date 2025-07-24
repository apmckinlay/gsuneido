// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	tsField: false
	RequirementsMet?()
		{
		.tsField = QueryColumns(.GetBaseQuery()).FindOne({ it.Suffix?('_TS') })
		return .tsField isnt false
		}

	Init()
		{
		.lastModified = Display(.GetData()[.tsField])
		}

	RunAddon()
		{
		if QueryEmpty?(.GetKeyQuery() $ ' where ' $ .tsField $ ' is ' $ .lastModified)
			.Reload()
		}
	}
