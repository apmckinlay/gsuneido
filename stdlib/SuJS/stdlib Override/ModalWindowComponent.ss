// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
WindowComponent
	{
	AnchorElement: 'dialog'
	DisableMinimize?: true
	New(@args)
		{
		super(@args)
		.AnchorEl.AddEventListener('cancel', .On_Cancel)
		.AnchorEl.showModal();
		SuRender().CancelDelayedTask(id: 'any')
		}

	SetBackdropDismiss()
		{
		.AnchorEl.AddEventListener('click', .backdropDismiss)
		}

	On_Close()
		{
		.On_Cancel()
		}

	On_Cancel()
		{
		.Event('EscapeCancel')
		}

	backdropDismiss(event)
		{
		if event.target is .AnchorEl
			.On_Cancel()
		}

	PlaceActive()
		{
		}

	Destroy()
		{
		.AnchorEl.Close()
		super.Destroy()
		}
	}
