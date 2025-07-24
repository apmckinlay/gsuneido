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
		if false isnt field = SuRender().GetRegisteredComponent(fieldHwnd)
			{
			r = SuRender.GetClientRect(field.El)
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
