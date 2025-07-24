// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
/* AddItem and InsertItem expect as 'item' parameter an object with this structure:

	#(	name: 'string',
		id: ###,
		type: [Windows menu type flags],
		state: [Windows menu state flags],
		submenu: [Windows menu handle],
		data: [],
		bitmap: [hBitmap],
		cmd: 'string' 	// if the cmd you want to call is different from the display,
						// only works in ShowCallCascade
	)

All members with the exception of ARE optional. If the bitmap member
is present and not NULL, it is presumed to be a valid handle to a bitmap
and will be displayed either alone (if name is an empty string or
not present) or along side it, if name is present.  Data is an application
defined 32-bit value. */
class
	{
	currentid: 1
	New(@items)
		{
		// Create Windows menu object and get handle to it
		.Hmenu = .createPopupMenu()
		.items = Object()
		.menuMap = Object()
		while items.Size() is 1 and items.Member?(0) and Object?(items[0])
			items = items[0]
		if items.HasNamed?()
			items = Object(items)
		.AddItems(items, .Hmenu)
		}
	// Maintenance (add/remove/insert) functions
	InsertItem(item = false, pos = 0, newhandle = false, prefixes = #())
		{
		if not Object?(item)
			return
		newMenu = false
		if item.Member?('name')
			{
			if item.name is ""
				item = Object(type: MFT.SEPARATOR)
			else
				{
				name = TranslateLanguage(item.name)
				.items[name] = item.name
				newMenu = Object(:name, :prefixes, cmd: item.GetDefault(#cmd, item.name))
				item.name = name
				// Specifying 'type' if it is just a string is optional...
				item.type = item.GetDefault('type', 0) | MFT.STRING
				}
			}
		if not item.Member?('id')
			item.id = .currentid++
		else
			.currentid += (item.id < .currentid) ? 1 : (1 + item.id - .currentid)
		.insertMenuItem(item, newhandle, pos, newMenu)
		}
	insertMenuItem(item, newhandle, pos, newMenu)
		{
		if 0 is mask = .GetMask(item)
			return
		info = Object(
			cbSize:		MENUITEMINFO.Size(),
			fMask:		mask,
			fType: 		item.GetDefault(#type, 0),
			fState:		item.GetDefault(#state, 0),
			wID:		item.GetDefault(#id, 0),
			dwItemData:	item.GetDefault(#data, 0),
			dwTypeData:	item.GetDefault(#name, 0),
			cch:		item.Member?('name') ? item.name.Size(): 0,
			//hbmpItem:	item.Member?('bitmap') ? item.bitmap : 0
			)
		// Must have at least one property...
		hMenu = newhandle is false ? .Hmenu : newhandle
		InsertMenuItem(hMenu, pos, true, info)
		if item.Member?('def')
			SetMenuDefaultItem(.Hmenu, item.id, false)
		if newMenu isnt false and -1 isnt id = GetMenuItemID(hMenu, pos)
			.menuMap[id] = newMenu
		}

	AddItems(items, hCurrentSubMenu, prefixes = #()) // Recursive
		{
		for (i = 0; i < items.Size(); i++)
			{
			item = items[i]
			x = .getMenuItemCount(hCurrentSubMenu)
			if Object?(item)
				.InsertItem(item.Copy(), x, hCurrentSubMenu, :prefixes)
			else if String?(item)
				.InsertItem(Object(name: item), x, hCurrentSubMenu, :prefixes)
			else
				continue

			i += .handleSubMenu(items, hCurrentSubMenu, prefixes, i, x)
			}
		}

	// extracted for testing
	getMenuItemCount(hCurrentSubMenu)
		{
		return GetMenuItemCount(hCurrentSubMenu)
		}

	handleSubMenu(items, hCurrentSubMenu, prefixes, i, pos)
		{
		if i + 1 >= items.Size() or not Object?(items[i + 1]) or items[i + 1].HasNamed?()
			return 0

		if 0 isnt hNewMenu = .createPopupMenu()
			{
			.setMenuItemInfo(hCurrentSubMenu, hNewMenu, pos)
			.AddItems(items[i + 1], hNewMenu, prefixes.Copy().Add(.getItemCmd(items[i])))
			}
		return 1
		}

	// extracted for testing
	createPopupMenu()
		{
		return CreatePopupMenu()
		}
	setMenuItemInfo(hCurMenu, hNewMenu, pos)
		{
		SetMenuItemInfo(hCurMenu, pos, true,
			[cbSize: MENUITEMINFO.Size(), fMask: MIIM.SUBMENU, hSubMenu: hNewMenu])
		}

	getItemCmd(item)
		{
		if String?(item)
			return .toIdentifier(item)
		if Object?(item) and item.Member?(#name) and item.name isnt ""
			return .toIdentifier(item.name)
		throw "Invalid menu item"
		}

	// Utility functions
	GetMask(obj)
		{
		// Return a valid MENUITEMINFO fMask member
		return (obj.Member?('type') ? MIIM.TYPE : 0) |
				(obj.Member?('state') ? MIIM.STATE : 0) |
				(obj.Member?('data') ? MIIM.DATA : 0) |
				(obj.Member?('bitmap') ? MIIM.BITMAP : 0) |
				//(obj.Member?('name') ? MIIM.STRING : 0) |
				(obj.Member?('id') ? MIIM.ID : 0) |
				(obj.Member?('submenu') ? MIIM.SUBMENU : 0)
		}
	// Interface functions
	Show(hwnd, x = 0, y = 0, left = false, rcExclude = 0, buttonRect = #())
		// pre: rcExclude is like #(left: 0, right: 10, top: 0, bottom: 10)
		// 		*** in screen (not window) coordinates
		{
		align = .adjustAlign(hwnd, x, buttonRect)
		x = .adjustX(align, x, buttonRect)
		if Object?(rcExclude)
			rcExclude = Object(cbSize: TPMPARAMS.Size(), :rcExclude)
		button = left is false ? TPM.RIGHTBUTTON : TPM.LEFTBUTTON
		// Display the context menu and return the index of the item that was selected
		i = TrackPopupMenuEx(.Hmenu, align | button | TPM.RETURNCMD,
			x, y, hwnd, rcExclude)
		.Destroy()
		return i
		}
	adjustAlign(hwnd, x, buttonRect = #())
		{
		align = TPM.LEFTALIGN  //default
		if buttonRect.Empty?()
			return align
		menuWidth = .menuWidth(hwnd)
		wa = GetWorkArea(buttonRect)
		//handle multiple-monitors checking
		return wa.left < 0
			? x.Abs() < menuWidth
				? TPM.CENTERALIGN : align
			: wa.right < x + menuWidth
				? TPM.CENTERALIGN : align
		}
	adjustX(align, x, buttonRect = #())
		{
		if buttonRect.Empty?()
			return x
		wa = GetWorkArea(buttonRect)
		//popup should "stick to" edge of screen as much as possible
		x = align is TPM.CENTERALIGN and not buttonRect.Empty?()
			? buttonRect.right //right edge adjustment
			: x
		return wa.left > x ? wa.left : x  //left edge adjustment
		}

	menuWidth(hwnd)
		{
		WithDC(hwnd)
			{|dc|
			DoWithHdcObjects(dc, [Suneido.hfont])
				{
				longest = ''
				for item in .items
					if .textExtent(hwnd, item) > .textExtent(hwnd, longest)
						longest = item
				// Boarder in pixels for each side of item text
				boarder = 50 /*= side boarder*/ * 2
				return .textExtent(hwnd, longest) + boarder + 1
				}
			}
		}

	textExtent(dc, s)
		{
		GetTextExtentPoint32(dc, s, s.Size(), ex = Object())
		return ex.x
		}

	ShowCall(ctrl, x = 0, y = 0, hwnd = false, rcExclude = 0)
		{
		result = .showcall(ctrl, x, y, hwnd, rcExclude)
		.Destroy()
		return result
		}

	popupMenu(ctrl, x, y, hwnd, rcExclude)
		{
		// Exit if ctrl isnt valid control
		if not ((Instance?(ctrl) and (ctrl.Base?(Control))) or Number?(hwnd))
			return false
		tpm = Object?(rcExclude)
			? Object(cbSize: TPMPARAMS.Size(), :rcExclude)
			: 0
		// when openning context menu and destroying control (alt+F4) at same time
		if ctrl.Empty?()
			return false
		// Display context menu
		if hwnd is false
			hwnd = ctrl.Window.Hwnd
		return TrackPopupMenuEx(.Hmenu, TPM.LEFTALIGN | TPM.RIGHTBUTTON | TPM.RETURNCMD,
			x, y, hwnd, tpm)
		}

	showcall(ctrl, x, y, hwnd, rcExclude)
		{
		i = .popupMenu(ctrl, x, y, hwnd, rcExclude)
		if i <= 0
			return i
		if false is item = GetMenuItemInfoText(.Hmenu, i)
			return false
		// get the untranslated text from .items
		if .items.Member?(item)
			item = .items[item]
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

	ShowCallCascade(ctrl, x = 0, y = 0, hwnd = false, rcExclude = 0)
		{
		i = .popupMenu(ctrl, x, y, hwnd, rcExclude)
		result = .callCascade(i, ctrl)
		.Destroy()
		return result
		}

	callCascade(i, ctrl)
		{
		if i <= 0
			return i

		prefixes = #()
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

	MakeItemIntoCommand(item)
		{
		return "On_Context_" $ .toIdentifier(item)
		}

	toIdentifier(item)
		{
		cmd = item.BeforeFirst('\t') // strip keyboard accelerator
		return ToIdentifier(cmd)
		}

	Destroy()
		{
		// Free the menu's resources:
		DestroyMenu(.Hmenu)
		.Hmenu = 0
		}
	}
