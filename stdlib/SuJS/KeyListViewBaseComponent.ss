// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	New(@args)
		{
		super(@args)
		.fieldCtrl = .FindControl('Field')
		.list = .FindControl('VirtualListGrid')
		.go = .FindControl('')
		.fieldCtrl.El.AddEventListener('keydown', .keydown)
		.fieldCtrl.El.AddEventListener('wheel', .wheel)
		}

	keydown(event)
		{
		if event.key in ('ArrowUp', 'ArrowDown')
			{
			.list.Event('SetFocus')
			.list.Event(#KEYDOWN, event.key is 'ArrowUp' ? VK.UP : VK.DOWN, 0,
				ctrl: event.ctrlKey, shift: event.shiftKey)
			}
		else if event.key is #Enter
			.Event('FieldEnter')
		}

	wheel(event)
		{
		.list.El.scrollTop = .list.El.scrollTop + event.deltaY
		}
	}
