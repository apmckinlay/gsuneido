// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// base for controls that are window procedures
// map should be an object mapping notification message numbers to method names
// TODO: switch to using SetWindowSubclass
// http://blogs.msdn.com/b/oldnewthing/archive/2003/11/11/55653.aspx
Hwnd
	{
	prevproc: false
	SubClass()
		{
		if .prevproc isnt false
			return
		.prevproc = SetWindowProc(.Hwnd, GWL.WNDPROC, this)
		Assert(.prevproc isnt 0)
		.thisproc = GetWindowLongPtr(.Hwnd, GWL.WNDPROC)
		}
	NCDESTROY(wParam, lParam)
		{
		if .prevproc isnt false
			{
			if GetWindowLongPtr(.Hwnd, GWL.WNDPROC) is .thisproc
				SetWindowProc(.Hwnd, GWL.WNDPROC, .prevproc)
			ClearCallback(this)
			prevproc = .prevproc
			.prevproc = false
			return CallWindowProc(prevproc, .Hwnd, WM.NCDESTROY, wParam, lParam)
			}
		// have to return 'callsuper' (rather than just 0)
		// or else Windows leaks something
		// and after about 1500 windows you get weird errors
		return 'callsuper'
		}
	Call(hwnd, msg, wParam, lParam)
		{
		_hwnd = .WindowHwnd()
		if 0x200 <= msg and msg <= 0x209 /*= mouse messages*/
			(.relay)(.Hwnd, msg, wParam, lParam)
		if WMmap.Member?(msg) and .Method?(method = WMmap[msg])
			{
			result = this[method](:wParam, :lParam)
			if result isnt 'callsuper'
				return result
			}
		if .prevproc isnt false
			return CallWindowProc(.prevproc, hwnd, msg, wParam, lParam)
		else
			return DefWindowProc(hwnd, msg, wParam, lParam)
		}
	Callsuper(hwnd, msg, wParam, lParam)
		{
		if .prevproc isnt false
			return CallWindowProc(.prevproc, hwnd, msg, wParam, lParam)
		else
			return DefWindowProc(hwnd, msg, wParam, lParam)
		}
	relay(@unused) { } // dummy
	ToolTip(tip)
		{
		.SubClass() // needed to relay
		.SetRelay(.Window.Tips().RelayEvent)
		super.ToolTip(tip)
		}
	SetRelay(func)
		{ // NOTE: only supports one relay
		.relay = func
		}
	COMMAND(wParam, lParam) /*internal*/
		{
		if .Window.HwndMap.Member?(lParam)
			// old style notification - reflect to sender
			return .Window.HwndMap[lParam].Notify(HIWORD(wParam), 0)
		else
			return 'callsuper'
		}
	NOTIFY(lParam) /*internal*/
		{
		// new style notification - reflect to sender
		nmhdr = NMHDR(lParam)
		if .Window.HwndMap.Member?(n = nmhdr.hwndFrom)
			return .Window.HwndMap[n].Notify(nmhdr.code, lParam)
		else
			return 'callsuper'
		}
	CTLCOLOREDIT(wParam, lParam)
		{
		if lParam isnt .Hwnd and .Window.HwndMap.Member?(lParam)
			{
			ctrl = .Window.HwndMap[lParam]
			if ctrl.Method?('CTLCOLOREDIT') and ctrl isnt this
				return ctrl.CTLCOLOREDIT(:wParam, :lParam)
			}
		return 'callsuper'
		}
	CTLCOLORSTATIC(wParam, lParam)
		{
		if lParam isnt .Hwnd and .Window.HwndMap.Member?(lParam)
			{
			ctrl = .Window.HwndMap[lParam]
			if ctrl.Method?('CTLCOLORSTATIC') and ctrl isnt this
				return ctrl.CTLCOLORSTATIC(:wParam, :lParam)
			}
		return 'callsuper'
		}
	CTLCOLORLISTBOX(wParam, lParam)
		{
		if lParam isnt .Hwnd and .Window.HwndMap.Member?(lParam)
			{
			ctrl = .Window.HwndMap[lParam]
			if ctrl.Method?('CTLCOLORLISTBOX') and ctrl isnt this
				return ctrl.CTLCOLORLISTBOX(:wParam, :lParam)
			}
		return 'callsuper'
		}
	CTLCOLORBTN(wParam, lParam)
		{
		if lParam isnt .Hwnd and .Window.HwndMap.Member?(lParam)
			{
			ctrl = .Window.HwndMap[lParam]
			if ctrl.Method?('CTLCOLORBTN') and ctrl isnt this
				return ctrl.CTLCOLORBTN(:wParam, :lParam)
			}
		return 'callsuper'
		}
	DRAWITEM(wParam, lParam)
		{
		dis = DRAWITEMSTRUCT(lParam)
		hwnd = dis.hwndItem
		if hwnd isnt .Hwnd and .Window.HwndMap.Member?(hwnd)
			{
			ctrl = .Window.HwndMap[hwnd]
			if ctrl.Method?('DRAWITEM') and ctrl isnt this
				return ctrl.DRAWITEM(:wParam, :dis)
			}
		return 'callsuper'
		}
	CONTEXTMENU(wParam, lParam)
		{
		// The low-order word of lParam gives the mouse's x-coordinate
		// and the high-order word gives the y-coordinate (screen coordinates)
		if .Window.HwndMap.Member?(wParam)
			{
			x = LOSWORD(lParam)
			if x is -1
				x = 0
			y = HISWORD(lParam)
			if y is -1
				y = 0
			return .Window.HwndMap[wParam].ContextMenu(x, y)
			}
		else
			return 'callsuper'
		}
	VSCROLL(wParam, lParam)
		{
		if lParam isnt .Hwnd and .Window.HwndMap.Member?(lParam)
			{
			ctrl = .Window.HwndMap[lParam]
			if ctrl.Method?('VSCROLL') and ctrl isnt this
				return ctrl.VSCROLL(:wParam, :lParam)
			}
		return 'callsuper'
		}
	HSCROLL(wParam, lParam)
		{
		if lParam isnt .Hwnd and .Window.HwndMap.Member?(lParam)
			{
			ctrl = .Window.HwndMap[lParam]
			if ctrl.Method?('HSCROLL') and ctrl isnt this
				return ctrl.HSCROLL(:wParam, :lParam)
			}
		return 'callsuper'
		}
	DoContextMenu(menu, x, y)
		{
		i = ContextMenu(menu).Show(.Hwnd, x, y) - 1
		if i is -1
			return 0
		.ContextMenuCall(menu[i])
		return 0
		}
	ContextMenuCall(menu)
		{
		this['On_' $ ToIdentifier(menu.BeforeFirst('\t'))]()
		}
	DevMenu: ('', 'Inspect Control', 'Copy Field Name', 'Go To Field Definition')
	On_Inspect_Control()
		{
		Inspect(this)
		}
	On_Go_To_Field_Definition()
		{
		GotoLibView('Field_' $ .getFieldName())
		}
	On_Copy_Field_Name()
		{
		fieldName = .getFieldName()
		.Copy_Field_Name(fieldName)
		}
	Copy_Field_Name(fieldName)
		{
		ClipboardWriteString(fieldName)
		InfoWindowControl(fieldName $ ' is copied to clipboard',
			titleSize: 0, marginSize: 7, autoClose: 1)
		}
	getFieldName()
		{
		name = .Send('GetFieldName')
		return String?(name)
			? name
			: .getFieldNameRec(this)
		}
	getFieldNameRec(inst)
		{
		if ((inst.Name in ('Value', 'Horz', 'twoList', 'Uom') or
			Gotofind('Field_' $ inst.Name, LibraryTables(), true).Empty?()) and
			inst.Member?('Parent'))
			return .getFieldNameRec(inst.Parent)
		return inst.Name
		}
	DoDevContextMenu(x, y)
		{
		if Suneido.User is 'default'
			return .DoContextMenu(.DevMenu, x, y)
		return 'callsuper'
		}
	}