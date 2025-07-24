// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
ChooseListControl
	{
	Name: 'FieldHistory'
	New(width = 20, mandatory = false, selectFirst = false,
		font = "", trim = true, size = "", field = #(), cue = false)
		{
		super(#(), :width, :mandatory, allowOther:,	:font, :size, :trim, :field, :cue)
		.SetList(.list = .getlist())
		if (selectFirst isnt false and not .list.Empty?())
			{
			.Set(.list[0])
			.Send('NewValue', .list[0])
			}
		}
	getlist()
		{
		if .Name is 'FieldHistory'
			return Object()
		if not Suneido.Member?('FieldHistory')
			Suneido.FieldHistory = Object()
		if not Suneido.FieldHistory.Member?(.Name)
			Suneido.FieldHistory[.Name] = Object()
		return Suneido.FieldHistory[.Name]
		}

	NewValue(value)
		{
		.updateList(value)
		super.NewValue(value)
		}
	updateList(value)
		{
		if not .list.Has?(value)
			.list.Delete(9) /*= max size */
		.list.Remove(value)
		.list.Add(value, at: 0)
		}
	Getter_Hwnd()
		{
		return .Field.Hwnd
		}

	FieldReturn()
		{
		.Send('FieldReturn')
		}
	}