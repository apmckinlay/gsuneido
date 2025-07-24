// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// centers subsequent dialog on parent
// e.g. for messagebox or common dialogs
// Used by BrowseFolderName, OpenFileName, SaveFileName
// NOTE: Not needed for Dialog - it centers itself
class
	{
	CallClass(hwnd, block)
		{
		hook = new .hookproc
		hook.Instance = GetWindowLongPtr(hwnd, GWL.HINSTANCE)
		hook.Hook = SetWindowsHookEx(WH.CALLWNDPROCRET,
			hook, hook.Instance, GetCurrentThreadId())
		hook.Owner = hwnd
		return Finally(block, {
			if not hook.Unhooked?
				{
				if not UnhookWindowsHookEx(hook.Hook)
					throw "can't UnhookWindowsHookEx(" $ hook.Hook $ ")"
				}
			ClearCallback(hook) })
		}

	hookproc: class
		{
		Unhooked?: false
		Call(nCode, wParam, lParam)
			{
			if nCode >= 0
				{
				msg = CWPRETSTRUCT(lParam)
				if msg.message is WM.INITDIALOG
					{
					.centerWindow(.Owner, msg.hwnd)
					.Unhooked? = UnhookWindowsHookEx(.Hook)
					}
				}
			return CallNextHookEx(.Hook, nCode, wParam, lParam )
			}
		centerWindow(hwndParent, hwndChild)
			{
			if false is IsWindow(hwndParent) or false is IsWindow(hwndChild)
				return
			if Object?(rcParent = GetWindowRect(hwndParent)) and
				Object?(rcChild = GetWindowRect(hwndChild))
				{
				x = rcParent.left + ((rcParent.right - rcParent.left) -
					(rcChild.right - rcChild.left)) / 2
				y = rcParent.top + ((rcParent.bottom - rcParent.top) -
					(rcChild.bottom - rcChild.top)) / 2
				SetWindowPos(hwndChild, NULL, x, y, 0, 0,
					SWP.NOSIZE | SWP.NOZORDER | SWP.NOACTIVATE)
				}
			}
		}
	}