// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
EditorControl
	{
	New(list, width, style = 0, mandatory = false, .listField = false,
		font = "", readonly = false, .delimiter = ',', height = 2)
		{
		super(:width, :style, :mandatory, :font, :readonly, :height)
		.listarg = list.Copy().Map!(String)
		}
	getter_list()
		{
		if .listField is false
			return .listarg
		list = .Send("GetField", .listField)
		if String?(list)
			list = list.Split(.delimiter)
		for i in list.Members()
			list[i] = list[i].Trim()
		return list
		}
	KillFocus()
		{
		dirty? = .Dirty?()
		.Set(.Get().Split(.delimiter).Map!(#Trim).Join(.delimiter))
		// resume the dirty? flag so EditorControl.EN_KILLFOCUS will send NewValue if needed
		.Dirty?(dirty?)
		}

	Valid?()
		{
		if false is .validCheck?(.Get().Split(.delimiter).Map!(#Trim), .list)
			return false

		return super.Valid?()
		}

	validCheck?(selectedList, validList)
		{
		return selectedList.UniqueValues() is selectedList and // no duplicates
			selectedList.Subset?(validList)
		}

	SetList(list)
		{
		.list = list.Copy().Map!(String) // disables getter_list
		}
	}