// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
KeyListViewBase
	{
	New(query, columns, .saveInfoName = "", prefix = "", .access = false,
		prefixColumn = false, keys = false, field = "", value = "", .excludeSelect = #(),
		customizeQueryCols = false, optionalRestrictions = #(), startLast = false)
		{
		super(query, columns, .saveInfoName, prefix, prefixColumn, keys,
			field, value, :customizeQueryCols, :optionalRestrictions, :startLast)

		if saveInfoName isnt ''
			{
			.selectMgr = AccessSelectMgr(name: saveInfoName)
			.selectMgr.LoadSelects(this)
			}
		.Redir('On_Select', this)
		}

	Layout(query, columns, saveInfoName, prefix, prefixColumn, keys)
		{
		if (.access isnt false and ((Suneido.Member?('AccessGoToCount') and
			Suneido.AccessGoToCount >= 2) or
			AccessGoTo.CheckPermission(.access).permission is false))
			.access = false

		layoutOb = .BaseLayout(query, columns, saveInfoName, prefix,
			prefixColumn, keys, false, false)

		layoutOb.Add(#(Skip 5),
			Object('Horz',
				'Fill',
				(.access is false ? 'Skip' : #(Button 'Access')),
				'Skip',
				#(Button 'Select...'),
				'Skip',
				#(Button Cancel),
				#(Skip 2)),
			#(Skip 2))
		return layoutOb
		}

	VirtualList_Return()
		{
		.VirtualList_LeftClick(.GetList().GetSelectedRecord())
		}
	VirtualList_LeftClick(rec)
		{
		if rec is false
			return
		return .Window.Result(Object(rec, .GetPrefixBy().Get()))
		}
	VirtualList_Escape()
		{
		.On_Cancel()
		}
	VirtualList_DoubleClick(@unused)
		{
		return false
		}
	Commands: #(("Select",	"Alt+S"))
	VirtualList_On_Select()
		{
		.On_Select()
		return true
		}

	// search/select interface
	OverrideSelectManager?()
		{
		return true
		}
	On_Select()
		{
		SelectControl(this, okbutton:,  name: .saveInfoName, hideCount:)
		SetFocus(.GetField().Hwnd)
		.FieldChange()
		}
	Getter_Option() // used by Select presets
		{ return '' }
	GetFields()
		{
		return QueryColumns(.GetBaseQuery())
		}
	sf: false
	GetSelectFields()
		{
		if .sf is false
			.sf = SelectFields(.GetFields(), .GetExcludeSelectFields())
		return .sf
		}
	GetExcludeSelectFields()
		{
		return .excludeSelect
		}
	SetSelectVals(select_vals)
		{
		if .saveInfoName isnt ''
			.selectMgr.SetSelectVals(select_vals, .GetSelectFields())
		}
	Getter_Select_vals()
		{
		if  .saveInfoName isnt ''
			return .selectMgr.Select_vals()
		return Record()
		}

	On_Access()
		{
		field = .GetField()
		AccessGoTo(.access, .GetPrefixColumn(), field.Get(), .Window.Hwnd,
			onDestroy: {
				if not .Empty?()
					{
					list = .GetList()
					// Force the virtual list to reload if it has read all the data
					if list.GetModel().AllRead? is true
						list.Refresh()
					SetFocus(field.Hwnd)
					.FieldChange()
					}
				 })
		}

	Destroy()
		{
		if .saveInfoName isnt ''
			.selectMgr.SaveSelects()
		super.Destroy()
		}
	}
