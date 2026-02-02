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
		.El.SetAttribute(#maxlength, EditorTextLimit)
		}

	onKeyDown(event)
		{
		pressed = Object(
			control: event.GetDefault(#ctrlKey, false),
			shift: event.GetDefault(#shiftKey, false),
			alt: event.GetDefault(#altKey, false))
		EditorKeyDownComponentHandler(this, event, pressed)
		}
	}
