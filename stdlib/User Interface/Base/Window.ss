// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// top level window
WindowBase
	{
	sizeMenu:	0
	// override Control's CallClass
	CallClass(@args)
		{
		new this(@args)
		}
	New(control, title = false,
		x = false, y = false, w = false, h = false,
		exStyle = 0, style = 0, wndclass = "SuBtnfaceArrow",
		.show = true, .newset = false, .exitOnClose = false,
		.keep_placement = false, border = 0, .onDestroy = false, parentHwnd = 0,
		useDefaultSize = false, skipStartup? = false)
		{
		super(:border)
		.menutips = Object().Set_default("")
		.dynmenu = Object()
		.dynid = Object()
		.commands = Object()
		if (style is 0)
			style = WS.OVERLAPPEDWINDOW
		.createWindow(x, w, y, h, wndclass, style, exStyle, parentHwnd)

		.SubClass()
		.SetupCommands(control)
		.create(control)
		.title = .FindTitle(title, control)
		SetWindowText(.Hwnd, TranslateLanguage(.title))
		.createmenu()
		.createSizeMenu()
		.RunPendingRefresh()

		if .keep_placement isnt false and .restoreWindowPlacement()
			UpdateWindow(.Hwnd)
		else
			.showWindow(style, exStyle, useDefaultSize, show, border)

		.SetAccels()

		if not skipStartup?
			DoStartup(.Ctrl)
		}
	createWindow(x, w, y, h, wndclass, style = 0, exStyle = 0, parentHwnd = 0)
		{
		.Hwnd = CreateWindowEx(exStyle, wndclass, "", style,
			.adjust(x), .adjust(y), .adjust(w), .adjust(h),
			parentHwnd, 0, Instance(), 0)
		if .Hwnd is 0
			throw "CreateWindow failed"
		}
	adjust(n)
		{
		return n is false
			? CW_USEDEFAULT
			// if width/height <= 1 then assumed to be fractions of screen size
			: n isnt false and 0 < n and n < 1
				? n * GetSystemMetrics(SM.CXMAXIMIZED)
				: n
		}
	create(control)
		{
		if (.newset isnt false)
			{
			try
				.Ctrl = .Construct(control)
			catch (x)
				{
				Alert(x, title: "Error Creating PersistentWindow",
					flags: MB.ICONERROR)
				Exit()
				}
			}
		else
			.Ctrl = .Construct(control)
		if .show isnt false and .show isnt SW.SHOWNA and .show isnt SW.SHOWNOACTIVATE
			SetFocus(GetNextDlgTabItem(.Hwnd, NULL, false))
		}
	createmenu()
		{
		if (not .Ctrl.Member?("Menu"))
			return
		menus = .Ctrl.Val_or_func(#Menu).Copy()
		.addPluginMenus(menus)
		menubar = CreateMenu()
		for (menu in menus)
			{
			pulldown = CreateMenu()
			.makesubmenu(menu, pulldown)
			AppendMenu(menubar, MF.POPUP, pulldown, TranslateLanguage(menu[0]))
			}
		SetMenu(.Hwnd, menubar)
		}
	addPluginMenus(menus)
		{
		attached = Object().Set_default(false)
		Plugins().ForeachContribution("UI", 'attach')
			{|c|
			if c.to is .Ctrl.Title
				attached[c.menu] = true
			}
		map = Object()
		Plugins().ForeachContribution("UI", 'action')
			{|c|
			if attached[c.menu] is false
				continue
			if not map.Member?(c.menu)
				menus.Add(map[c.menu] = Object(c.menu))
			map[c.menu].Add(c.name)
			if c.Member?(#target)
				.menu_redir(c.name, c.target)
			if c.Member?(#accel)
				.AddAccel(.menuify(c.name), c.accel)
			}
		.SetAccels()
		order = Object()
		Plugins().ForeachContribution("UI", 'menu')
			{|c|
			order[c.menu] = c.order
			}
		defaultOrder = 5
		menus.Sort!({|x, y|
			order.GetDefault(x[0], defaultOrder) < order.GetDefault(y[0], defaultOrder) })
		}
	menu_redir(name, target)
		{
		if String?(target)
			target = Global(target)
		if Object?(name)
			.Ctrl.Redir("Menu_" $ .menuify(name = name[0]), target)
		.Ctrl.Redir("On_" $ .menuify(name), target)
		}
	menuify(name)
		{
		return ToIdentifier(name)
		}
	makesubmenu(menu, pulldown)
		{
		for (i = 1; i < menu.Size(); ++i)
			{
			m = menu[i]
			if String?(m)
				.makemenuentry(m, pulldown)
			else if Object?(m)
				{
				p = CreateMenu()
				if (m.Size() is 1)
					.dynmenu[p] = .menuify(m[0])
				.makesubmenu(m, p)
				AppendMenu(pulldown, MF.POPUP, p, TranslateLanguage(m[0]))
				}
			}
		}
	makemenuentry(m, pulldown)
		{
		if (m is '')
			{
			AppendMenu(pulldown, MF.SEPARATOR, 0, 0)
			return
			}

		cmd = .menuify(m)
		c = .Commands().GetDefault(cmd,
			Object(id: .Mapcmd(cmd), accel: "", help: ""))
		m = TranslateLanguage(m)
		accel = TranslateLanguage(c.accel)
		if m is '&Copy' and accel is ''
			accel = 'Ctrl+C'
		AppendMenu(pulldown, MF.STRING, c.id, m $ '\t' $ accel)
		.menutips['On_' $ cmd] = TranslateLanguage(c.help)
		}

	showWindow(style, exStyle, useDefaultSize, show, border)
		{
		GetClientRect(.Hwnd, r = Object())
		r.right = .Ctrl.Xstretch <= 0
			? .Ctrl.Xmin + 2 * border : Max(.Ctrl.Xmin, r.right)
		r.bottom = .Ctrl.Ystretch <= 0
			? .Ctrl.Ymin + 2 * border
			: Min(Max(.Ctrl.Ymin, r.bottom), .Ctrl.MaxHeight)
		AdjustWindowRectEx(r, style, .Ctrl.Member?("Menu"), exStyle)
		if useDefaultSize is true
			.SetWinSize(@.DefaultSize())
		else
			.SetWinSize(r.right - r.left, r.bottom - r.top)
		if show isnt false
			.Show(show is true ? SW.SHOWNORMAL : show)
		}
	INITMENUPOPUP(wParam)
		// pre: wParam is the menu handle
		// post: menu is updated
		{
		if (not .dynmenu.Member?(wParam) or
			0 is menu = .Send('Menu_' $ .dynmenu[wParam]))
			return 0
		for (i = GetMenuItemCount(wParam) - 1; i >= 0; --i)
			DeleteMenu(wParam, i, MF.BYPOSITION)
		for (m in menu)
			{
			if (m is "")
				AppendMenu(wParam, MF.SEPARATOR, 0, 0)
			else
				{
				id = .Mapcmd(m)
				.dynid[id] = Object('On_' $ .dynmenu[wParam], m)
				AppendMenu(wParam, MF.STRING, id, TranslateLanguage(m))
				}
			}
		return 0
		}
	COMMAND(wParam, lParam) /*internal*/
		{
		id = LOWORD(wParam)
		if (.dynid.Member?(id))
			{
			.Send(@.dynid[id])
			return 0
			}
		else
			return super.COMMAND(wParam, lParam)
		}

	focus: 0
	ACTIVATE(wParam)
		{
		if (LOWORD(wParam) is WA.INACTIVE)
			{
			.focus = GetFocus()
			.Send("Inactivate")
			PubSub.Publish(#WindowInactivated)
			}
		else
			{
			if .focus isnt 0
				SetFocus(.focus)
			.Send("Activate")
			PubSub.Publish(#WindowActivated)
			}
		return 0
		}
	MENUSELECT(wParam)
		{
		cmd = .Cmdmap(LOWORD(wParam))
		tip = .menutips[cmd]
		.Send("MenuSelect", tip)
		return 0
		}
	Show(cmd)
		{
		// Call ShowWindow, passing cmd
		ShowWindow(.Hwnd, cmd)
		// If the window is not being minimized, repaint it...
		if ((cmd isnt SW.MINIMIZE) and (cmd isnt SW.SHOWMINIMIZED))
			UpdateWindow(.Hwnd) // sends WM_PAINT message
		}
	sizeMenuOptions: (
		(width: 640, height: 480),
		(width: 800, height: 600),
		(width: 1024, height: 768),
		(width: 1280, height: 1024),
		(width: 1600, height: 900),
		(width: 1920, height: 1080),
		(width: 2560, height: 1440)
		)
	createSizeMenu()
		{
		.windowMenus.Add(menu = CreateMenu())
		for v in .sizeMenuOptions
			.appendWindowMenu(menu, v.width $ ' x ' $ v.height,
				.buildSizeCallable(v.width, v.height))

		m = GetSystemMenu(.Hwnd, false)
		item = .initMenuItem('S&et Size', menu, MFT.SEPARATOR)
		InsertMenuItem(m, 5, true, item) /*= between Maximize and Close*/
		item.fType = MFT.STRING
		InsertMenuItem(m, 6, true, item) /*= between Maximize and Close*/
		}

	getter_windowMenus()
		{ return .windowMenus = Object() }

	// Have to isolate the building of the callable to avoid issues with blocks / local
	// variables. (Otherwise the last callable interposes over all the previous ones)
	buildSizeCallable(w, h)
		{
		return { .setSize(ScaleWithDpiFactor(w), ScaleWithDpiFactor(h)) }
		}

	appendWindowMenu(menu, prompt, callable)
		{
		cmd = .windowMenuCommands.Empty?()
			? 0xEEF1 /*= first command value */
			: .windowMenuCommands.Members().Max() + 1
		AppendMenu(menu, MFT.STRING, cmd, prompt)
		.windowMenuCommands[cmd] = callable
		}

	getter_windowMenuCommands()
		{ return .windowMenuCommands = Object() }

	initMenuItem(prompt, menu, fType)
		{
		item = Object()
		item.cbSize = MENUITEMINFO.Size()
		item.fMask = MIIM.TYPE | MIIM.ID | MIIM.SUBMENU
		item.fType = fType
		item.wID = .windowWID++
		item.dwTypeData = prompt
		item.cch = item.dwTypeData.Size()
		item.hSubMenu = menu
		return item
		}

	getter_windowWID()
		{ return .windowWID = 0xEEFF }

	AddWindowMenuOptions(menus, idx = 7) // idx: 7 puts the added options AFTER "Set Size"
		{
		// Init menu and insert a separator
		item = .initMenuItem('', NULL, MFT.SEPARATOR)
		InsertMenuItem(m = GetSystemMenu(.Hwnd, false), idx, true, item)
		menus.Each()
			{
			if it.GetDefault('options', false) is false
				{
				InsertMenuItem(m, ++idx, true,  .initMenuItem(it.root, NULL, MFT.STRING))
				.windowMenuCommands[GetMenuItemID(m, idx)] = it.cmd
				}
			else
				{
				.windowMenus.Add(menu = CreateMenu())
				it.options.Members().Sort!().Each()
					{ |m| .appendWindowMenu(menu, m, it.options[m]) }
				InsertMenuItem(m, ++idx, true, .initMenuItem(it.root, menu, MFT.STRING))
				}
			}
		}

	SYSCOMMAND(wParam)
		{
		if false is .windowMenuCommands.Member?(wParam)
			return 'callsuper'
		(.windowMenuCommands[wParam])()
		return 0
		}

	setSize(w, h)
		{
		h -= 24	/*= allow for "normal" taskbar to better "preview"
					what window would look like at that resolution */

		//TODO fix this to handle multi-monitor
		rc = SPI_GetWorkArea()
		wa_width = rc.right - rc.left
		wa_height = rc.bottom - rc.top
		w = Min(w, wa_width)
		h = Min(h, wa_height)
		if (w is wa_width and h is wa_height)
			ShowWindow(.Hwnd, SW.MAXIMIZE)
		else
			{
			ShowWindow(.Hwnd, SW.NORMAL)
			SetWindowPos(.Hwnd, NULL,
				rc.left + ((wa_width - w) / 2).Int(),
				rc.top + ((wa_height - h) / 2).Int(),
				w,
				h,
				SWP.NOZORDER)
			}
		}
	APP_SETFOCUS(lParam)
		{
		SetFocus(lParam)
		}

	// ignore e.g. when called by On_Cancel
	Result(unused)
		{ }

	EnableClose(enable = true)
		{
		EnableMenuItem(GetSystemMenu(.Hwnd, false), SC_CLOSE, MF.BYCOMMAND |
			(enable is true ? 0 : MF.GRAYED))
		}

	getter_windowPlacementSaveName()
		{
		saveName = 'wp: '
		if String?(.keep_placement) and .keep_placement isnt ''
			saveName $= .keep_placement
		else
			saveName $= .title
		return saveName
		}
	orig_info: false
	restoreWindowPlacement()
		{
		if false is info = KeyListViewInfo.Get(.windowPlacementSaveName)
			return false
		.orig_info = info.window_info.Copy()
		if info.window_info.showCmd is SW.SHOWMINIMIZED
			info.window_info.showCmd = SW.SHOWNORMAL
		req = .RequiredWindowSize()
		r = info.window_info.rcNormalPosition
		.adjustToScreen(r)
		if r.right - r.left < req.w or .Ctrl.Xstretch <= 0
			r.right = r.left + req.w
		if r.bottom - r.top < req.h or .Ctrl.Ystretch <= 0
			r.bottom = r.top + req.h
		SetWindowPlacement(.Hwnd, info.window_info)
		return true
		}
	adjustToScreen(r)
		{
		wa = GetWorkArea(r)
		// make sure the window size is not larger than the screen size
		r.right = Min(r.right, r.left + wa.right - wa.left)
		r.bottom = Min(r.bottom, r.top + wa.bottom - wa.top)
		// if part of the window is out of the screen, move it in
		if r.left < wa.left
			{
			r.right += wa.left - r.left
			r.left = wa.left
			}
		if r.right > wa.right
			{
			r.left -= r.right - wa.right
			r.right = wa.right
			}
		if r.top < wa.top
			{
			r.bottom += wa.top - r.top
			r.top = wa.top
			}
		if r.bottom > wa.bottom
			{
			r.top -= r.bottom - wa.bottom
			r.bottom = wa.bottom
			}
		}
	saveWindowPlacement()
		{
		// program errors while creating print previews causes title to be uninitialized,
		// could also be other cases
		if not .Member?("Window_title")
			return
		GetWindowPlacement(.Window.Hwnd,
			place = Object(length: WINDOWPLACEMENT.Size()))
		if place isnt .orig_info
			KeyListViewInfo.Save(.windowPlacementSaveName, place)
		}

	GetWindowTitle()
		{
		return .title
		}

	ResetStyle()
		{
		.Ctrl.Send('ResetAddons')
		.Ctrl.TopDown('ResetAddons')
		.Ctrl.Send('ResetTheme')
		.Ctrl.TopDown('ResetTheme')
		}

	DESTROY(_windowResult = false)
		{
		if .keep_placement isnt false
			.saveWindowPlacement()
		.windowMenus.Each(DestroyMenu)
		.ResetAccels()
		if .onDestroy isnt false
			(.onDestroy)(:windowResult)
		super.DESTROY()
		if .exitOnClose
			Exit()
		return	0
		}
	}
