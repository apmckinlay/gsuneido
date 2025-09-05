// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin: 750
	Ymin: 475 // within a few pixels of previous min. height
	CallClass(ctrl, name = '', okbutton = false, defaultButton = '',
		noUserDefaultSelects? = false, setSelectVals = false)
		{
		ToolDialog(0,
			Object(this, ctrl, name, :okbutton, :defaultButton, :noUserDefaultSelects?,
				:setSelectVals),
				closeButton?: false, keep_size: 'Select~' $ name)
		}

	Title: "Select"
	New(access, .name, okbutton = false, defaultButton = '',
		noUserDefaultSelects? = false, setSelectVals = false)
		{
		super(.layout(access, okbutton, defaultButton, noUserDefaultSelects?,
			setSelectVals))
		.select2 = .Vert.SelectRepeat
		}
	layout(access, okbutton, defaultButton, noUserDefaultSelects?, setSelectVals)
		{
		.DefaultButton = defaultButton
		.access = access

		.sf = access.GetSelectFields()
		selects = setSelectVals is false ? access.Select_vals : setSelectVals
		.remove_invalid_selects(selects)

		return Object('Vert',
			Object('SelectRepeat', .sf, selects, .name, option: access.Option,
				title: access.Title, :noUserDefaultSelects?),
			okbutton is true ? .buttons("Select") : .buttons(@.AccessButtons))
		}
	AccessButtons: ("First", "Last", "Current") // public to allow overriding
	remove_invalid_selects(selects)
		{
		selects.RemoveIf()
			{
			not it.Member?('condition_field') or not .sf.Fields.Has?(it.condition_field)
			}
		}

	CountBtnTip: 'Show the number of records matching the current Select'
	buttons(@list)
		{
		ob = Object('HorzEven', 'Skip')

		if not list.Has?('Select') // KeyListView
			list.Add("Count")
		tip = ''
		for button in list.Add("Uncheck All", "Cancel")
			{
			if button is 'Count'
				tip = .CountBtnTip
			ob.Add(Object('Button', button, :tip), 'Skip')
			}
		return ob
		}

	On_First()
		{
		if (false is .set_query())
			return
		.access.On_First()
		.closeDialog()
		}

	closeDialog()
		{
		.access.Send('SelectControl_Changed')
		.Window.Result("")
		}

	On_Last()
		{
		if (false is .set_query())
			return
		.access.On_Last()
		.closeDialog()
		}

	On_Current()
		{
		if false is where = .where()
			return false
		.access.ModifyWhere(where, hwnd: .Window.Hwnd)
		.closeDialog()
		}

	On_Select()
		{
		if (false is .set_query())
			return
		.closeDialog()
		}

	set_query()
		{
		if false is where = .where()
			return false
		return .access.SetWhere(where, hwnd: .Window.Hwnd)
		}
	where()
		{
		// Close the select control as the access has been destroyed
		if .access.Destroyed?()
			{
			.Window.Result("")
			return false
			}

		if false is where = .select2.Where()
			return false
		.access.SetSelectVals(.select2.Get().conditions)
		return .sf.Joins(where.joinflds) $ where.where
		}
	On_Uncheck_All()
		{
		newConditions = .select2.Get().conditions.Map({ it.check = false; it })
		.select2.Set([conditions: newConditions])
		}
	On_Cancel()
		{
		.Window.Result("")
		}

	On_Count()
		{
		if (false is .set_query())
			return
		.access.On_Count()
		}
	}
