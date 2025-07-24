// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
WindowComponent
	{
	New(@args)
		{
		super(@args)
		.edit = SuRender().GetRegisteredComponent(args.editHwnd)
		}

	PlaceElement()
		{
		if .edit isnt false
			{
			r = .edit.GetListPos()
			rcExclude = Object(left: -9999, right: 9999, top: r.top, bottom: r.bottom)
			PlaceElement(.GetContainerEl(), r.left, r.bottom, rcExclude, r)
			}
		}

	PlaceActive()
		{
		}
	}
