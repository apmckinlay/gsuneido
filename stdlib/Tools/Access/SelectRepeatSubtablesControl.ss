// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	// Our problem is AccessControl needs to drive all this, as it is the one that
	// knows what the subtables are. List/MultiView DOES NOT
	// accessControl has the ference to AccessSubtables
	// Second issue, we have two instances of the SelectRepeatControls, only one is
	// consturted at a time.

	New(.view, .colModel, .selectRepeatName)
		{
		.Name = .selectRepeatName
		hdr = .FindControl('Header')
		hdr.Redir('Record_NewValue', this)
		.selectControls = Object(Header: hdr)
		.gatherControls()
		}

	gatherControls()
		{
		if .extraFilterNames is false
			return
		// we need the at:  to be the save name
		// we need the .FindControl to use the linked.name
		for saveName in .extraFilterNames.Members()
			{
			x = .FindControl(.extraFilterNames[saveName])
			x.Redir('Record_NewValue', this)
			.selectControls.Add(x at: saveName)
			}
		}

	Recv(@args)
		{
		source = args.source
		if source.Parent.Name is 'buttons' and args[0].Prefix?('On_')
			.Controller.Send(@args)
		return 0
		}

	extraFilterNames: false
	Controls()
		{
		layout = Object('Vert',
			Object('SelectRepeat',
				.view.GetSelectFields(), .view.Select_vals, 'Header',
				option: .view.Option, title: .view.GetTitle(),
				selChanged: .view.GetDefault('SelectChanged?', false),
				noUserDefaultSelects?: not .colModel.UserDefaultSelectEnabled?()),
			)

		extraFilters = .Send('Select_ExtraFilters')
		.extraFilterNames = Object()
		for extraFilter in extraFilters
			{
			layout.Add(extraFilter.control)
			.extraFilterNames[extraFilter.saveName] = extraFilter.linkedName
			}

		layout.Add(.buttons())
		return layout
		}

	buttons()
		{
		if 0 is extraLayout = .Send('Select_ExtraLayout')
			extraLayout = #()
		buttons = Object('HorzEqual',
			Object('EnhancedButton', 'Select', command: 'Select', buttonStyle:,
				mouseEffect:, weight: 'bold', pad: 40,
				name: 'loadButton', xstretch: 0), name: 'buttons')
		buttons.Add(#(Skip, small:))
		buttons.Add(Object('Button', 'Count', tip: SelectControl.CountBtnTip))
		buttons.Add('Fill')
		VirtualListTopLayoutControl.BuildLayout(extraLayout, buttons)
		buttons.name = 'buttons'
		return buttons
		}

	Record_NewValue(field, value/*unused*/)
		{
		if field is 'conditions'
			.SetSelectApplied(false)
		}

	Valid?(quiet = false)
		{
		for ctrl in .selectControls
			if ctrl.Valid?(quiet) is false
				return false
		return true
		}

	On_Count()
		{
		if not .Valid?()
			return
		.Send('On_Count')
		}

	On_Select()
		{
		if not .Valid?()
			return
		.Send('Select_Apply')
		.SetSelectApplied(true)
		}

	SetSelectApplied(applied = false)
		{
		loadButton = .FindControl('loadButton')
		if loadButton is false
			return
		.Send('SelectControl_SetSelectApplied', .Valid?)
		loadButton.SetTextColor(applied ? CLR.BLACK : CLR.RED)
		}

	SelectChanged?()
		{
		loadButton = .FindControl('loadButton')
		return loadButton.GetTextColor() is CLR.RED
		}

	Get(item = false)
		{
		ob = Object()
		for name in .selectControls.Members()
			ob.Add(.selectControls[name].Get() at: name)
		if item isnt false
			return ob[item]
		return ob
		}

	// THIS is currently called from  Addon_VirtualListTopFilters.Select_Apply
	Where(selectFields)
		{
		where = .FindControl('Header').Where()
		x = selectFields.Joins(where.joinflds) $ where.where
		x $= .Send('Select_ExtraWhere', .selectControls)
		return x
		}

	Set(vals, extra_select_vals = false)
		{
		.FindControl('Header').Set(vals)
		if extra_select_vals is false
			return

		for saveName in extra_select_vals.Members()
			{
			selval = extra_select_vals[saveName]
			.selectControls[saveName].Set(selval)
			}
		}

	HighlightLastRow(row)
		{
		.FindControl('Header').HighlightLastRow(row)
		}

	}