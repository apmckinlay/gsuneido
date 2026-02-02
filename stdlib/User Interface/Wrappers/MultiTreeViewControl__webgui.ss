// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
TreeViewControl
	{
	New(.multi? = true, readonly = false, style = 0)
		{ super(readonly, style) }

	Getter_Selection()
		{
		selection = Object()
		if 0 isnt item = .GetSelectedItem()
			selection.Add(item)
		return selection
		}
	}
