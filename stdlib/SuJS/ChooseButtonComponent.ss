// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
MenuButtonComponent
	{
	New(.text, .list, width = false)
		{
		super(text, left:, :width)
		}

	Keydown(event)
		{
		if event.key =~ "^\w$"
			.char(event.key)
		else
			super.Keydown(event)
		}

	char(c)
		{
		list = .list.Assocs().Filter({ it[1].Prefix?(c) })
		if list.Empty?()
			return 0
		i = list.FindIf({ it[1] is .El.innerText })
		i = i is false ? 0 : (i + 1) % list.Size()
		return .Event(#On_ChooseButton, list[i][1], index: list[i][0])
		}

	SetList(.list) { }
	}
