// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Window
	{
	ComponentName: 'ModalWindow'
	CallClass(control, title  = false, border = 5, closeButton? = true,
		onDestroy = false, keep_size  = true, useDefaultSize = false)
		{
		w = super.CallClass(control, :title, :border,
			style: .style(closeButton?),
			:onDestroy, :keep_size , :useDefaultSize)
		w.DialogCenterSize(w.GetWindowTitle(), control, keep_size, :useDefaultSize)
		CreateContributionSysMenus('DialogMenus', w)
		return w
		}

	style(closeButton?)
		{
		return (closeButton? ? WS.SYSMENU : WS.CAPTION) | WS.SIZEBOX
		}

	Result(result)
		{
		if not .AllowCloseWindow?()
			return

		_windowResult = result
		.Window.DESTROY()
		}

	DESTROY()
		{
		.UnregisterWindow()
		super.DESTROY()
		}
	}