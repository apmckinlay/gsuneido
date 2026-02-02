// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
class
	{
	currentid: 1
	New(@items)
		{
		.menu = Object()
		.menuMap = Object()
		while items.Size() is 1 and items.Member?(0) and Object?(items[0])
			items = items[0]
		if items.HasNamed?()
			items = Object(items)
		.AddItems(items, .menu)
		}

	AddItems(items, menu, prefixes = #()) // Recursive
		{
		for (i = 0; i < items.Size(); i++)
			{
			item = items[i]
			if Object?(item)
				.insertItem(item.Copy(), menu, :prefixes)
			else if String?(item)
				.insertItem(Object(name: item), menu, :prefixes)
			else
				continue

			i += .handleSubMenu(items, prefixes, i, menu)
			}
		}

	insertItem(item, menu, prefixes = #())
		{
		if not Object?(item) or not item.Member?('name')
			return

		if item.name is ""
			{
			menu.Add("")
			.currentid++
			}
		else
			{
			if not item.Member?('id')
				item.id = .currentid++
			else
				.currentid = Max(item.id, .currentid) + 1
			name = TranslateLanguage(item.name)
			.menuMap[item.id] = Object(:name, :prefixes,
				cmd: item.GetDefault(#cmd, item.name))
			item.name = name

			menu.Add(item)
			}
		}

	handleSubMenu(items, prefixes, i, menu)
		{
		if i + 1 >= items.Size() or not Object?(items[i + 1]) or items[i + 1].HasNamed?()
			return 0

		newMenu = Object()
		.AddItems(items[i + 1], newMenu, prefixes.Copy().Add(.getItemCmd(items[i])))
		menu.Last().submenu = newMenu
		return 1
		}

	getItemCmd(item)
		{
		if String?(item)
			return .toIdentifier(item)
		if Object?(item) and item.Member?(#name) and item.name isnt ""
			return .toIdentifier(item.name)
		throw "Invalid menu item"
		}

	toIdentifier(item)
		{
		cmd = item.BeforeFirst('\t') // strip keyboard accelerator
		return ToIdentifier(cmd)
		}

	Show(hwnd, x = 0, y = 0, left/*unused*/ = false, rcExclude = 0,
		buttonRect = #(), parent = false)
		{
		if .menu.Empty?()
			return 0

		if false is id = Dialog(0,
			Object('ContextMenuList', .menu.Map(.formatMenuItem), x, y,
				rcExclude, buttonRect, :parent),
			border: 0, style: WS.POPUP, backdropDismiss?:,
			posRect: x is false and y is false ? hwnd : false)
			return 0 // dismissed
		return id
		}

	ShowCall(ctrl, x = 0, y = 0, hwnd = false, rcExclude = 0)
		{
		return .showcall(ctrl, x, y, hwnd, rcExclude)
		}

	showcall(ctrl, x, y, hwnd, rcExclude)
		{
		i = .popupMenu(ctrl, x, y, hwnd, rcExclude)
		if i <= 0
			return i

		if not .menuMap.Member?(i)
			return false

		// get the untranslated text from .items
		item = .menuMap[i].cmd
		callmethod = .MakeItemIntoCommand(item)
		// Call method in control corresponding to item's text string, allowing
		// ctrl to try and handle non-existent methods
		try
			ctrl[callmethod](:item)
		catch (err, "method not found:")
			{
			if err.Has?(callmethod)
				SuneidoLog('ERROR: ' $ err)
			else
				throw err
			}
		return i
		}

	popupMenu(ctrl, x, y, hwnd/*unused*/, rcExclude)
		{
		// Exit if ctrl isnt valid control
		if not (Instance?(ctrl) and ctrl.Base?(Control))
			return false
		// when openning context menu and destroying control (alt+F4) at same time
		if ctrl.Empty?()
			return false
		if .menu.Empty?()
			return false
		// Display context menu
		if false is id = Dialog(0,
			Object('ContextMenuList', .menu.Map(.formatMenuItem), x, y, rcExclude),
			border: 0, style: WS.POPUP, backdropDismiss?:)
			return 0 // dismissed
		return id
		}

	formatMenuItem(item)
		{
		if not Object?(item)
			return item
		ob = item.Project(#id, #name)
		if item.Member?(#type)
			ob.type = item['type']
		if item.Member?(#dynamic)
			ob.dynamic = true
		if item.Member?(#submenu)
			ob.submenu = item.submenu.Map(.formatMenuItem)
		if ((item.GetDefault(#state, 0) & MFS.DISABLED) isnt 0)
			ob.disable = true
		return ob
		}

	MakeItemIntoCommand(item)
		{
		return "On_Context_" $ .toIdentifier(item)
		}

	ShowCallCascade(ctrl, x = 0, y = 0, hwnd = false, rcExclude = 0)
		{
		i = .popupMenu(ctrl, x, y, hwnd, rcExclude)
		return .callCascade(i, ctrl)
		}

	callCascade(i, ctrl)
		{
		if i <= 0
			return i

		if not .menuMap.Member?(i)
			return false

		prefixes = .menuMap[i].prefixes
		cmd = .menuMap[i].cmd

		calls = prefixes.Copy().Add(cmd)
		sep = calls.Size()
		while sep isnt 0
			{
			callmethod = .MakeItemIntoCommand(calls[..sep].Join('_'))
			if ctrl.Method?(callmethod) or ctrl.Method?('Default')
				ctrl[callmethod](item: calls[sep..])
			sep--
			}
		return i
		}
	}