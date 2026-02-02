// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
// purpose: to assist users in entering comma separated lists
ChooseField
	{
	Name: CommaList
	New(mandatory = false, .max_items = false, .numeric? = false, noFieldEdit? = false,
		.protectLineFunc = false)
		{
		super(Object('Field', readonly: noFieldEdit?), :mandatory)
		}

	Getter_DialogControl()
		{
		return Object('CommaListDlg', .Field.Get(), .max_items, .protectLineFunc)
		}

	Get()
		{
		return .Field.Get().Split(',').Map!(#Trim).Join(',')
		}

	Valid?()
		{
		list = .Field.Get().Split(',')
		if .numeric? is true and not .all_numeric?(list)
			return false
		if .max_items isnt false and list.Size() > .max_items
			return false
		return .Field.Valid?()
		}

	FieldKillFocus()
		{
		.SetValid(.Valid?())
		}

	all_numeric?(list)
		{
		return list.Every?({ it.Trim().Number?() })
		}

	On_DropDown()
		{
		.Send('CommaList_SetNextNums')
		super.On_DropDown()
		}
	}