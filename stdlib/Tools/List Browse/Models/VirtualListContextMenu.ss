// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
class
	{
	menu: #()
	headerMenu: #()
	recMenu: false

	New(menu = false, headerMenu = false, .recMenu = false, .addCurrentMenu? = false,
		.addGlobalMenu? = false)
		{
		.defaultHeaderMenus = Object("Reset Columns", "", "Print...", "Reporter...")

		if Suneido.User is 'default'
			.defaultHeaderMenus.Add('', 'Go To QueryView',
				'Inspect Control', 'Copy Field Name', 'Go To Field Definition')
		if menu isnt false
			.menu = menu
		if headerMenu isnt false
			.headerMenu = headerMenu
		}

	BuildMenu(model, view, preventCustomExpand?, excludeCustomize?)
		{
		if model.ExpandModel isnt false
			{
			layoutOb = view.Send('VirtualList_Expand', [])
			if 0 isnt layoutOb
				{
				if model.ExpandModel.CustomizableExpand?(layoutOb)
					.AddCustomizeExpand()
				}
			else if not preventCustomExpand?
				.AddCustomizeExpand()
			}
		// move contextMenu and following code into model
		if not excludeCustomize? and model.EditModel.Editable?() and
			model.ColModel.GetCustomKey() isnt false
			.AddCustomize()
		.AddCustomizeColumn(model.ColModel.GetColumnsSaveName())
		}

	AddCustomize()
		{
		if not .defaultHeaderMenus.Has?('Customize...')
			.defaultHeaderMenus.Add('Customize...', at: 0)
		}

	AddCustomizeColumn(columnsSaveName)
		{
		if columnsSaveName isnt false and
			not .defaultHeaderMenus.Has?('Customize Columns...')
			.defaultHeaderMenus.Add('Customize Columns...', at: 0)
		}

	AddCustomizeExpand()
		{
		if not .defaultHeaderMenus.Has?('Customize Expand...')
			.defaultHeaderMenus.Add('Customize Expand...', at: 0)
		}

	ContextRec: false
	ContextCol: false
	ContextRowNum: false
	contextColumns: #()

	SetContext(rec, col, columns, row_num = false)
		{
		.ContextRec = rec
		.ContextCol = col
		.ContextRowNum = row_num
		.contextColumns = columns
		}

	SetHeaderMenu(newHeaderMenu)
		{
		Assert(Object?(newHeaderMenu))
		.headerMenu = newHeaderMenu
		}

	SetMenu(newMenu)
		{
		.menu = newMenu
		}

	ShowMenu(view, rec, col, row_num, point)
		{
		if 0 isnt menu = view.Send("VirtualList_BuildContextMenu", :rec)
			{
			if rec is false
				.SetHeaderMenu(menu)
			else
				.SetMenu(menu)
			}
		model = view.GetModel()
		.SetContext(rec, col, model.ColModel.GetColumns(), row_num)
		.Show(view, point.x, point.y)
		}

	Show(control, x, y, extraMenu = #())
		{
		m = .buildMenu(control, extraMenu)
		ContextMenu(m).ShowCallCascade(control, x, y)
		}

	buildMenu(control, extraMenu = #())
		{
		m = .ContextRec is false
			? .headerMenu.Copy().Add(@.defaultHeaderMenus)
			: .menu.Copy()
		.addToMenu(m, extraMenu, pos: m.Find(""))
		if control.Editable?()
			.addToMenu(m, Object("New"))
		if .ContextRec isnt false and control.Editable?()
			{
			selectedCount = control.GetSelectedRecords().Size()
			.buildRecordMenu(selectedCount, m)
			}
		if .addGlobalMenu? is true
			{
			if m.NotEmpty?()
				.addToMenu(m, Object(""))
			.addToMenu(m, Object(.recMenu.Global, "Global"))
			}
		if control.GetHeaderSelectPrompt() isnt 'no_prompts'
			.addFormatMenu(m)
		return m
		}

	buildRecordMenu(selectedCount, m)
		{
		if selectedCount is 1
			{
			.addToMenu(m, Object("Edit Field"), pos: 0)
			if .addCurrentMenu? isnt false
				for cMenu in .recMenu.Current
					.addToMenu(m, Object(cMenu))
			}
		}

	addToMenu(new_menus, ob, pos = false)
		{
		if ob.Empty?()
			return
		for menu in ob.Reverse!()
			if pos is false
				new_menus.Add(menu)
			else
				new_menus.Add(menu, at: pos)
		}

	addFormatMenu(menu)
		{
		if .ContextCol is false
			return

		fmt = Datadict(.ContextCol).Format[0]
		if not fmt.Suffix?('Format')
			fmt $= 'Format'
		fmt = Global(fmt)
		if fmt.Method?('List_ExtraContext')
			if false isnt contextExtra = fmt.List_ExtraContext()
				menu.Append(Object("", contextExtra))
		}

	Inspect()
		{
		for col in .contextColumns
			.ContextRec[col]
		Inspect.Window(.ContextRec)
		}

	HandlePluginOption(option, ctrl)
		{
		option = option.AfterFirst('On_Context_').Replace('_', ' ')
		if .recMenu isnt false and .recMenu.Current.Has?(option)
			{
			.recMenu.On_Current(option, .ContextRec, ctrl)
			return true
			}
		return false
		}

	UpdateHistory(rec)
		{
		if .recMenu isnt false
			.recMenu.UpdateHistory(rec, rec.New?())
		}

	RedirectContextMenu(view, args)
		{
		event = args[0]
		if event.Prefix?('On_Context_')
			{
			if view.Addons.Send(event)
				return
			if .HandlePluginOption(event, view)
				return
			view.Send(event, rec: .ContextRec, col: .ContextCol, contextMenu: this,
				item: args.item)
			}
		else
			throw 'method not handled: ' $ Display(args)
		}
	}
