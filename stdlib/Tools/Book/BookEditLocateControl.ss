// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Name: BookEditLocate
	table: false
	New(.table = '')
		{
		super(['AutoChoose', .RecordNames, width: 20, autoSelect:,
			cue: 'Locate by name', allowOther:])
		}

	RecordNames(prefix)
		{
		records = .records(.table)
		if prefix is ""
			return records
		match = '\<(?i)(?q)' $ prefix
		list = records.Filter({ it.name =~ match or it.trimmed =~ match}) // exact matches
		return .buildlist(list)
		}

	SetTable(.table)
		{}

	records(table)
		{
		if table is false
			return #()
		records = Object()
		QueryApply(table $ ' where text isnt "" sort name')
			{|x|
			records.Add([path: x.path, name: x.name, trimmed: x.name.Tr(" ")])
			}
		return records
		}

	buildlist(list)
		{
		namesWithPath = Object()
		for item in list
			{
			path = item.path.AfterFirst(`/`).Replace(`/`,' > ')
			name = path.Blank?()
				? item.name
				: path $ ' > ' $ item.name
			namesWithPath.Add(name)
			}
		return namesWithPath.Sort!()
		}

	NewValue(value)
		{
		.value = value
		path = value.Replace(' > ', '/')
		path = path[0] isnt "/" ? "/" $ path : path
		.Send('Locate', .table $ path)
		}

	value: ''
	Get()
		{
		return .value
		}
	}