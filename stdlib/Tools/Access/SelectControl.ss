// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Xmin: 750
	Ymin: 475 // within a few pixels of previous min. height
	CallClass(ctrl, name = '', okbutton = false, defaultButton = '',
		noUserDefaultSelects? = false, setSelectVals = false, hideCount = false)
		{
		ToolDialog(0,
			Object(this, ctrl, name, :okbutton, :defaultButton, :noUserDefaultSelects?,
				:setSelectVals, :hideCount),
				closeButton?: false, keep_size: 'Select~' $ name)
		}

	Title: "Select"
	subTableSelects: false
	New(access, .name, okbutton = false, defaultButton = '',
		noUserDefaultSelects? = false, setSelectVals = false, .hideCount = false)
		{
		super(.layout(access, okbutton, defaultButton, noUserDefaultSelects?,
			setSelectVals))
		.select2 = .FindControl(.name)
		.subSelectRepeats()
		}
	subSelectRepeats()
		{
		if .access.Method?('UseSubTableFilters?') and
			.access.UseSubTableFilters?() is true
			{
			.subTableSelects = Object()
			for idx in .linkedBrowses.Members().Sort!()
				{
				linked = .linkedBrowses[idx]
				linkedBrowse = linked.browse
				saveName = .subTableFilterName(linkedBrowse)
				.subTableSelects.Add(.FindControl(linked.name) at: saveName)
				}
			}
		}

	layout(access, okbutton, defaultButton, noUserDefaultSelects?, setSelectVals)
		{
		.DefaultButton = defaultButton
		.access = access


		.sf = access.GetSelectFields()
		selects = setSelectVals is false ? access.Select_vals : setSelectVals
		.remove_invalid_selects(selects)

		layout = Object('Vert',
			Object('SelectRepeat', .sf, selects, .name, option: access.Option,
				title: access.Title, :noUserDefaultSelects?))

		// .access could also be a KeyListView
		if .access.Method?('UseSubTableFilters?') and
			.access.UseSubTableFilters?() is true
			{
			.linkedBrowses = .access.GetLinkedBrowseTabs()
			.layoutSubtables(layout)
			}

		buttons = okbutton is true ? .buttons("Select") : .buttons(@.AccessButtons)
		return layout.Add(buttons)
		}

	layoutSubtables(layout)
		{
		AccessSubtables.Layout(.access)
			{ |control, unused, unused|
			layout.Add(control)
			}
		}

	// this name is the saved name in userselects
	subTableFilterName(linkedBrowse)
		{
		return linkedBrowse.GetColumnsSaveName() $ ' Filter'
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
		if not .hideCount
			list.Add('Count')
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
		.access.ModifyWhere(where)
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
		return .access.SetWhere(where)
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
		headerWhere = .sf.Joins(where.joinflds) $ where.where

		lineItemWhere = ""
		if .access.Method?('UseSubTableFilters?') and
			.access.UseSubTableFilters?() is true
			lineItemWhere $= .access.SubTables_Where(.subTableSelects)
// Need to handle if this returns something NOT in the current re-name
		return headerWhere $ '\r\n' $ lineItemWhere
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