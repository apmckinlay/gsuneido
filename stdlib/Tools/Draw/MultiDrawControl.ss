// Copyright (C) 2015 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	Name: 'MultiDraw'
	New(items, .protectedTabs = #())
		{
		super(.layout(items))
		.tabs = .FindControl('Tabs')
		.tabs.SetImageList(AccessMarkTab.GetInitImageList())
		}

	layout(items)
		{
		tabs = Object('Tabs', buttonTip: 'Add New', newTab:)
		items.Each({ tabs.Add(.drawControl(it)) })
		return Object('Vert', tabs)
		}

	drawControl(item)
		{
		readOnly = .protectedTabs.Has?(item.name)
		return Object('Draw', valueOb: item.value, :readOnly, Tab: item.name)
		}

	Set(items)
		{
		itemCount = items.Size()
		tabsAvailable = .tabs.GetAllTabCount()
		for (i = tabsAvailable-1; i >= itemCount; i--) // Remove extra tabs
			.tabs.Remove(i)
		for i in ..itemCount
			if i < tabsAvailable
				.updateTabControls(i, items[i])
			else
				.insertNewTab(items[i].name, items[i].value)

		if not .tabs.Constructed?(i = .tabs.GetSelected())
			.tabs.SelectTab(i)
		}

	updateTabControls(i, item)
		{
		// Tab control is already constructed
		if false isnt control = .tabs.GetControl(i)
			control.Set(item.value)
		else
			{
			// Tab control is not yet constructed
			// Simply update the data to be used during construction
			data = .tabs.TabConstructData(i)
			data.valueOb = item.value
			}
		.SetTabData(i, [], item.name)
		}

	SetTabData(i, data, name = false)
		{
		.tabs.SetTabData(i, data, name)
		}

	Get()
		{
		return Seq(.TabsCount()).Map(.get1)
		}

	TabsCount()
		{
		return .tabs.GetAllTabCount()
		}

	get1(i, skipFormat = false)
		{
		value = .tabs.Constructed?(i)
			? .tabs.GetControl(i).Get()
			: .tabs.TabConstructData(i).valueOb
		tabData = .tabs.GetTabData(i)
		tabData.idx = i
		data = [name: .tabs.TabName(i), :value, :tabData]
		if not skipFormat
			.Send('FormatData', data)
		return data
		}

	On_Add_New()
		{
		if false isnt .Send('EditMode?')
			.addNewDraw()
		}

	existingNames()
		{
		return .tabs.GetAllTabNames()
		}

	addNewDraw()
		{
		result = (.addDrawDialog)(.Window.Hwnd, .existingNames())

		if result isnt false and result isnt ''
			{
			copyTabIndex = .Get().FindIf({ it.name is result.copy_from_draw_item })
			tabValue = copyTabIndex isnt false
				? .Get()[copyTabIndex].value
				: Object(items: #(), resources: #())
			.insertNewTab(result.new_draw_item_name, tabValue, newDraw:)
			}
		}

	ImportNewTab(name, value, data = false, title = 'Multi-Draw')
		{
		if .existingNames().Has?(name) and false is name = .renameImport(name, title)
			return 'Import cancelled'
		.insertNewTab(name, value, data, newDraw:)
		return true
		}

	renameImport(name, title)
		{
		return Ask(name $ ' already exists.\r\nPlease specify a new, ' $
			'unused name', title, hwnd: .Window.Hwnd, valid: .ValidateName)
		}

	// If newDraw is false, then the tab is not constructed until required
	insertNewTab(name, value, data = false, newDraw = false)
		{
		.tabs.Insert(name, .drawControl([:name, :value]), :data, noSelect: not newDraw)
		if newDraw
			.Send('NewValue')
		}

	ValidateName(name, existingNames = false)
		{
		if existingNames is false
			existingNames = .existingNames()
		if name.Blank?() or name =~ "[^a-zA-Z0-9 ]"
			return 'Name can not be blank or contain symbols'
		if existingNames.Has?(name)
			return name $ ' already exists.\r\nPlease enter another name.'
		return ''
		}

	addDrawDialog: Controller
		{
		Title: 'New Item'
		CallClass(hwnd, names)
			{
			return OkCancel(Object(this, names), .Title, hwnd)
			}
		New(existingItemNames)
			{
			super(Object('Record',
				Object('Vert'
					'new_draw_item_name'
					'copy_from_draw_item')))
			.existingItemNames = existingItemNames
			.data = .FindControl('Data')
			.data.Set(Record(draw_item_existing_names: .existingItemNames))
			}

		valid?()
			{
			if .data.Valid() isnt true
				return false
			name = .data.Get().new_draw_item_name
			if '' isnt msg = MultiDrawControl.ValidateName(name, .existingItemNames)
				.AlertInfo(.Title, msg)
			return '' is msg
			}
		OK()
			{
			if not .valid?()
				return false

			return .data.Get()
			}
		}

	Highlight(i, remove = false)
		{
		for j in .. .TabsCount()
			.tabs.SetImage(j, -1)
		.tabs.SetImage(i, remove ? -1 : 0)
		}

	SelectedTabData(format = false)
		{
		return .get1(.tabs.GetSelected(), skipFormat: not format)
		}

	CurTabIdx()
		{
		return .tabs.GetSelected()
		}

	DeleteDraw(i)
		{
		if .protectedTabs.Has?(.tabs.TabName(i))
			return
		.tabs.Select(Max(0, i - 1))
		.tabs.Remove(i)
		.Send('NewValue')
		}

	PreviewCurrent()
		{
		.tabs.GetControl().On_Print()
		}
	}
