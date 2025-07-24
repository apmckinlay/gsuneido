// Copyright (C) 2018 Axon Development Corporation All rights reserved worldwide.
EditComponent
	{
	Name: "Editor"
	Xstretch: 1
	Ystretch: 1
	DefaultHeight: 4
	Hasfocus?:	false
	New(@args)
		{
		super(@args)
		.El.AddEventListener('keydown', .onKeyDown)
		}

	onKeyDown(event)
		{
		pressed = Object(
			control: event.ctrlKey, shift: event.shiftKey, alt: event.altKey)
		EditorKeyDownComponentHandler(this, event, pressed)
		}
	}
