// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
class
	{
	windows: #()
	activeWindow: false
	New()
		{
		.windows = Object()
		.modalWindows = Object()
		.taskbar = Object()
		}

	IsWindow(hwnd)
		{
		return .windows.HasIf?({ it.Hwnd is hwnd })
		}

	GetActive()
		{
		return .activeWindow
		}

	ActivateWindow(window)
		{
		if Same?(.activeWindow, window)
			return window.Hwnd

		prev = .activeWindow
		if window isnt false
			{
			.windows.RemoveIf({ Same?(it, window) })
			.windows.Add(window)
			}

		.activate(prev, false)
		.activeWindow = window
		.activate(window, true)
		.syncWindowOrder()
		.UpdateTaskbar()
		return prev is false or not prev.Member?(#Hwnd)
			? NULL
			: prev.Hwnd
		}

	GetWindows()
		{
		return .windows
		}

	GetWindow(hwnd)
		{
		return .windows.FindOne({ it.Hwnd is hwnd })
		}

	IsWindowEnabled(hwnd)
		{
		for (i = .windows.Size() - 1; i >= 0; i--)
			{
			if .windows[i].Hwnd is hwnd
				return true
			if .windows[i].ComponentName in ('Dialog', 'ModalWindow')
				return false
			}
		Assert(false, msg: 'IsWindowEnabled: not find window ' $ Display(hwnd))
		}

	UnregisterWindow(window)
		{
		.windows.RemoveIf({ Same?(it, window) })
		.modalWindows.RemoveIf({ Same?(it, window) })
		.RemoveTaskbarWindow(window)
		if Same?(.activeWindow, window)
			{
			.activate(window, false)
			.activeWindow = .windows.GetDefault(.windows.Size() - 1, false)
			.activate(.activeWindow, true)
			}
		.syncWindowOrder()
		.UpdateTaskbar()
		}

	activate(window, active)
		{
		if window is false
			return

		window.ACTIVATE(active)
		}

	syncWindowOrder()
		{
		for i in .windows.Members()
			.windows[i].Act(#UpdateOrder, i, active?: Same?(.windows[i], .activeWindow))
		}

	AddModalWindow(window)
		{
		.modalWindows.Add(window)
		}

	ShowingModalWindow?(excludeModalWindow = false)
		{
		return .modalWindows.Filter({ not Same?(it, excludeModalWindow) }).NotEmpty?()
		}

	AddTaskbarWindow(window)
		{
		.taskbar.Add(window)
		.UpdateTaskbar()
		}

	RemoveTaskbarWindow(window)
		{
		if false is pos = .taskbar.FindIf({ Same?(it, window) })
			return
		.taskbar.Delete(pos)
		.UpdateTaskbar()
		}

	UpdateTaskbar()
		{
		res = Object()
		for window in .taskbar
			{
			if '' isnt title = window.GetWindowTitle()
				res.Add(Object(:title, id: window.UniqueId,
					active?: Same?(window, .GetActive())))
			}

		SuRenderBackend().CancelAction(false, 'SuTaskbarUpdate')
		SuRenderBackend().RecordAction(false, 'SuTaskbarUpdate', [titles: res])
		}

	ActivateNextWindow(current)
		{
		Assert(Same?(.activeWindow, current))
		next = false
		for (i = .windows.Size() - 1; i >= 0; i--)
			{
			if not Same?(current, .windows[i]) and .windows[i].IsMinimized?() isnt true
				{
				next = .windows[i]
				break
				}
			}
		.ActivateWindow(next)
		}
	}