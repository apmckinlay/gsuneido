// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Window
	{
	ComponentName: 'Dialog'
	CallClass(parentHwnd, control, style = 0, exStyle = 0, border = 5, title = false,
		posRect = false, keep_size = false, closeButton? = false, useDefaultSize = false,
		backdropDismiss? = false, beforeRun = false)
		{
		reserved = SuRenderBackend().ExtractReserved()
		if style is 0
			style = WS.CAPTION | WS.SIZEBOX | (closeButton? ? WS.SYSMENU : 0)

		parentHwnd = 0
		w = super.CallClass(control, :title, :style, :exStyle, :border, :parentHwnd,
			skipStartup?:)
		if backdropDismiss? is true
			w.Act(#SetBackdropDismiss)
		if beforeRun isnt false
			beforeRun(w)
		result = w.InternalRun(control, keep_size, posRect, useDefaultSize)
		SuRenderBackend().MergeReserved(reserved)
		return result
		}

	InternalRun(control, keep_size, posRect, useDefaultSize = false)
		{
		.setPos(control, keep_size, posRect, useDefaultSize)
//		SetWindowPos(.Hwnd, HWND.TOP, 0, 0, 0, 0, SWP.NOSIZE | SWP.NOMOVE)
//		.setDefaultButtonStyle()
		DoStartup(.Ctrl)
		CreateContributionSysMenus('DialogMenus', this)
		return .ActivateDialog()
		}

	setPos(control, keep_size, posRect, useDefaultSize)
		{
		.DialogCenterSize(.GetWindowTitle(), control, keep_size, useDefaultSize)
		if posRect isnt false
			.Act('AlignToField', fieldHwnd: posRect)
		}

	ActivateDialog()
		{
		JsWebSocketServer.MessageLoop()
		result = .resultValue
		.DESTROY()
		return result
		}

	FlushDelays() {}

	resultValue: false
	Result(value)
		{
		.resultValue = value
		throw WebSocketHandler.QUITLOOP
		}

	CLOSE()
		{
		if not .AllowCloseWindow?()
			return 0
		.closeDialog()
		return 0 // meaning we handled it so don't do default DestroyWindow
		}
	Destroy() // called by Controller.On_Close
		{
		if not .AllowCloseWindow?()
			return
		.closeDialog()
		}
	closeDialog()
		{
		.Result(false)
		}
	DoWithWindowsDisabled(block, exclude/*unused*/ = 0)
		{
		block()
		}
	}
