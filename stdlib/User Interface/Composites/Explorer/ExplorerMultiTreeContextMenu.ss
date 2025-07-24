// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.explorer, .readonly) { }

	getter_controller()
		{ return .controller = .explorer.Controller }

	getter_tree()
		{ return .tree = .explorer.Tree }

	// Purposely not setting .model = .explorer.Model as
	// .Model can be reconstructed multiple times in ExplorerMultiControl.
	// As a result, we need to ALWAYS point to the ExplorerMultiControl.Model directrly
	getter_model()
		{ return .explorer.Model }

	editOptions: 	(New, Delete, Rename, Cut, Copy, 'Copy To...', Paste, 'Move To...')
	disableMember: 	DisableContextMenuOptions
	getter_disabledOptions()
		{
		disabledOptions = .controller.GetDefault(.disableMember, []).Copy()
		if .readonly
			disabledOptions.MergeUnion(.editOptions)
		return .disabledOptions = disabledOptions
		}

	Default(@args)
		{ return .controller[args[0]](@+1args) }

	ShowTreeContextMenu(item, x, y)
		{
		if item is 0 // Right-Click occurs in empty space
			return 0
		selectionSize = .tree.Selection.Size()
		static? = .tree.Static?(item)
		container? = .tree.Container?(item)
		expanded? = .tree.Expanded?(item)
		children? = .tree.Children?(item)
		menu = .buildContextMenu(selectionSize, static?, container?, expanded?, children?)
		addonMenus = .controller.Method?(#Collect) ? .Collect(#ContextMenu, item) : []
		.orderMenu(menu, addonMenus)
		.cleanUpMenu(menu)
		return menu.Any?({ it.name isnt '' })
			? ContextMenu(menu).ShowCall(this, x, y, hwnd: .tree.Hwnd)
			: 0
		}

	segments: 0
	buildContextMenu(selectionSize, static?, container?, expanded?, children?)
		{
		.segments = 0
		menu = Object()
		.buildOpenCloseMenu(menu, selectionSize, container?, expanded?, children?)
		.buildItemMenu(menu, selectionSize, static?)
		return menu
		}

	buildOpenCloseMenu(menu, selectionSize, container?, expanded?, children?)
		{
		if selectionSize is 1 and container?
			{
			if not expanded? and children?
				.addToMenu(menu, options: [[name: '&Open', def:], ''])
			else if expanded?
				.addToMenu(menu, options: [[name: 'C&lose', def:], ''])
			}
		}

	paddingFactor: 10
	addToMenu(menu, options)
		{
		order = ++.segments * 10 /*= padding to allow addons to fit more easily*/
		if Object?(options)
			for option in options
				{
				if String?(option)
					option = [name: option]
				option.order = order++
				menu.Add(option)
				}
		else
			menu.Add([name: options, :order])
		}

	buildItemMenu(menu, selectionSize, static?)
		{
		.addCopyPaste(menu, selectionSize, static?)
		if not rootSelected? = .explorer.RootSelected?()
			.addToMenu(menu, #('&Delete', 'Rena&me', ''))
		.addToMenu(menu, '&New')
		if rootSelected?
			.addRootMenu(menu)
		}

	addRootMenu(menu)
		{
		rootMenu = ['', 'Dump']
		if .checkContoller(#AllowRootDelete?)
			rootMenu.Add('&Delete')
		.addToMenu(menu, rootMenu)
		}

	checkContoller(method)
		{ return .controller.Method?(method) and .controller[method]() }

	isLibViewControl?()
		{
		return .controller.Base?(LibViewControl)
		}

	addCopyPaste(menu, selectionSize, static?)
		{
		if not static?
			.addToMenu(menu, #('C&ut', '&Copy'))
		if allowPaste? = .checkContoller(#CanPaste?) and selectionSize is 1
			.addToMenu(menu, '&Paste')
		if not static? and .isLibViewControl?()
			.addToMenu(menu, #('Copy To...', '&Move To...'))
		if allowPaste? or not static?
			.addToMenu(menu, '')
		}

	orderMenu(baseMenu, addonMenus)
		{
		prevOrder = 0
		for addonMenu in addonMenus
			for option in addonMenu
				{
				if not option.Member?(#name) // Submenu of previous element
					option.order = prevOrder
				else if false isnt duplicate = .duplicateOption(baseMenu, option.name)
					{
					duplicate = baseMenu.Extract(duplicate)
					option.order = duplicate.order
					}
				baseMenu.Add(option)
				prevOrder = option.order
				}
		if not .disabledOptions.Empty?()
			baseMenu.RemoveIf({ .disabledOptions.Has?(it.name.Tr('&')) })
		baseMenu.Sort!(By(#order))
		baseMenu.Map!({ it.Delete(#order); it })
		}

	duplicateOption(baseMenu, name)
		{
		return baseMenu.
			FindIf({ it.Member?(#name) and it.name isnt '' and it.name is name })
		}

	// Removes any duplicate dividers as well as dividers at the start / end of the list
	cleanUpMenu(menu)
		{
		prevOption = ''
		menu.RemoveIf()
			{
			remove? = false
			if it.Member?(#name)
				{
				remove? = prevOption is '' and it.name is ''
				prevOption = it.name
				}
			remove?
			}
		if not menu.Empty?() and menu.First().GetDefault(#name, false) is ''
			menu.PopFirst()
		if not menu.Empty?() and menu.Last().GetDefault(#name, false) is ''
			menu.PopLast()
		}

	On_Context_Close()
		{ .tree.ExpandItem(.explorer.GetSelected(), true) }

	On_Context_Copy()
		{
		not_static? = not .explorer.Static?(.explorer.CurItem)
		if not_static?
			.clipboardCopy()
		ClipboardWriteString(.tree.GetName(.explorer.CurItem), add?: not_static?)
		}

	clipboardCopy(selected = false)
		{
		.Save()
		if selected is false
			selected = .tree.Selection

		data = Object()
		for hItem in selected
			data.Add(.clipboardCopyData(hItem))
		ClipboardWriteData(Pack(data), .controller.Clipformat)
		}

	clipboardCopyData(hItem)
		{
		item = .model.Get(.tree.GetParam(hItem), origText?:).Copy()
		return .copyData(item, .path(hItem), .tree.RootName(hItem))
		}

	copyData(item, path, table)
		{
		item.lib_before_path = item.GetDefault(#lib_before_path, path)
		item.table = table
		item.children = Object()
		path = Opt(path, '/') $ item.name
		for child in .model.Children(item.num)
			item.children.Add(.copyData(child, path, table))
		return item
		}

	path(item)
		{
		path = .tree.Path(item)
		return path.RemovePrefix(.controller.CurrentTable() $ '/').BeforeLast('/')
		}

	On_Context_Cut()
		{ .clipboardCut() }

	clipboardCut()
		{
		.clipboardCopy()
		.explorer.On_Delete_Item(false)	// false means do not give confirmation message
		}

	On_Context_Delete()
		{ .explorer.On_Delete_Item() }

	On_Context_Import_Records()
		{ .On_Import_Records(.getDisplayName()) }

	getDisplayName()
		{ return .tree.GetName(.tree.Selection[0]) }

	On_Context_Item()
		{ .explorer.On_New_Item() }

	On_Context_Folder()
		{ .explorer.On_New_Folder() }

	On_Context_Open()
		{ .tree.ExpandItem(.explorer.GetSelected()) }

	On_Context_Paste()
		{
		if false is curitem = .explorer.GetSelected()
			return
		else if not .tree.Container?(curitem)
			curitem = .tree.GetParent(curitem)
		.paste(curitem, Unpack(ClipboardReadData(.controller.Clipformat)))
		}

	paste(curitem, data)
		{
		table = .tree.RootName(curitem)
		if Object?(data)
			for pasteobj in data
				.tree.Insert(curitem, pasteobj, fromtree:, :table)
		}

	On_Context_Move_To()
		{
		result = .getSelectedAndNewPath('Move')
		if result.newpath is false
			return

		.explorer.GotoPath(result.newpath)
		newParent = .explorer.GetSelected()
		lastRec = result.selected.Last()
		result.selected.Each({ .tree.Move(it, newParent, last?: it is lastRec) })
		}
	On_Context_Copy_To()
		{
		result = .getSelectedAndNewPath('Copy')
		if result.newpath is false
			return

		.clipboardCopy(result.selected)
		.explorer.GotoPath(result.newpath)
		newParent = .explorer.GetSelected()
		.paste(newParent, Unpack(ClipboardReadData(.controller.Clipformat)))
		}
	getSelectedAndNewPath(title)
		{
		selected = .tree.Selection.Copy()
		if selected.Empty?()
			return Object(newpath: false, :selected)

		names = Object()
		selected.Each({ names.Add(.tree.GetName(it)) })
		pathVal = .tree.Path(selected[0])
		newpath = LibViewCopyMovePathControl(names.Sort!(), pathVal, title)
		return Object(:newpath, :selected)
		}

	On_Context_Rename()
		{
		item = .tree.GetSelectedItem()
		if .tree.Selection.Size() > 1
			.tree.UnselectAll(item)
		.tree.EditLabel(item)
		}

	On_Context_Dump()
		{ Database.Dump(.controller.CurrentTable()) }
	}
