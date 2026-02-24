// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
ModalWindowComponent
	{
	New(@args)
		{
		super(@args)

		// This is to trigger the server side to flush the delay tasks
		// JsWebSocketServer.MessageLoop doesn't flush delay tasks
		// These delays should be called after the Dialog is contructed and displayed
		.Event('FlushDelays')
		}

	AlignToField(fieldHwnd)
		{
		method = args = false
		// fieldHwnd can be an object (stdlib:GotoLibView.multiMatch)
		if Object?(fieldHwnd)
			{
			method = fieldHwnd.method
			args = fieldHwnd.args
			fieldHwnd = fieldHwnd.target
			}
		if false isnt field = SuRender().GetRegisteredComponent(fieldHwnd)
			{
			r = method is false
				? SuRender.GetClientRect(field.El)
				: (field[method])(@args)
			rcExclude = Object(left: -9999, right: 9999, top: r.top, bottom: r.bottom)
			// context menu place itself instead of the window container
			container = .Ctrl.Base?(ContextMenuListComponent)
				? .Ctrl.El
				: .GetContainerEl()
			fieldWindowRect = SuRender.GetClientRect(field.Window.GetContainerEl())
			PlaceElement(container, Max(0, Max(r.left, fieldWindowRect.left)), r.bottom,
				rcExclude, r)
			}
		}
	}
