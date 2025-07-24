// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// a Window that disables its parent while it's open
// similar to a Dialog, but does not run a nested message loop
// automatically centers itself on its parent
// automatically saves and restores size (if stretchable)
Window
	{
	New(control, title = false, border = 5, closeButton? = true,
		onDestroy = false, keep_size = true, useDefaultSize = false)
		{
		super(control, :title, :border,
			style: .style(closeButton?), exStyle: WS_EX.TOOLWINDOW,
			:onDestroy, parentHwnd: .parentHwnd = .GetParent(), show: false)
		.DialogCenterSize(.parentHwnd, .GetWindowTitle(), control, keep_size,
			:useDefaultSize)
		.prevDisabled = EnableWindow(.parentHwnd, false)
		CreateContributionSysMenus('DialogMenus', this)
		.Show(SW.SHOWNORMAL)
		}
	style(closeButton?)
		{
		return (closeButton? ? WS.SYSMENU : WS.CAPTION) | WS.SIZEBOX
		}

	Result(result)
		{
		if not .AllowCloseWindow?()
			return
		.reenableParent()
		_windowResult = result
		.Window.Destroy()
		}
	CLOSE()
		{
		if 0 isnt result = super.CLOSE()
			.reenableParent()
		return result
		}
	// this should be called before destroying, so Windows does not switch to other app
	restoreParent?: false
	reenableParent()
		{
		if .restoreParent?
			return
		if not .prevDisabled
			EnableWindow(.parentHwnd, true)
		.ModalClose(.parentHwnd)
		.restoreParent? = true
		}
	Destroy()
		{
		// handle when this method is called directly
		.reenableParent()
		super.Destroy()
		}
	}
