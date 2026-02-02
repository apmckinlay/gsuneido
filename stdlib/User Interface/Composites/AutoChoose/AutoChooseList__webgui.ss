// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_AutoChooseList
	{
	ComponentName: 'AutoChooseListComponent'
	CallClass(parent, list)
		{
		reservation = SuRenderBackend().ReserveAction()
		window = new this(parent, list)
		layout = window.GetLayout()
		args = [layout, style: WS.POPUP, uniqueId: window.UniqueId,
			parentHwnd: parent.Controller.UniqueId, editHwnd: parent.UniqueId]
		SuRenderBackend().RecordAction(false, .ComponentName, args, at: reservation.at)
		window.Act(#PlaceElement)
		return window
		}

	move()
		{
		}
	}
