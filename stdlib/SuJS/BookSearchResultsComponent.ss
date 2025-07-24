// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
HtmlDivComponent
	{
	New(@args)
		{
		super(@args)
		fieldHwnd = args.fieldHwnd
		if false isnt field = SuRender().GetRegisteredComponent(fieldHwnd)
			{
			r = SuRender.GetClientRect(field.El)
			rcExclude = Object(left: -9999, right: 9999, top: r.top, bottom: r.bottom)
			PlaceElement(.Window.GetContainerEl(), r.left, r.bottom, rcExclude, r)
			}
		}
	}
