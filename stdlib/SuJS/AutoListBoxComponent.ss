// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
ListBoxComponent
	{
	New(@args)
		{
		// delete xmin: 1000, ymin: 1000 in AutoChooseList
		super(@args.Delete(#xmin, #ymin))
		}

	InsertItem(s, i)
		{
		el = super.InsertItem(s, i)
		el.AddEventListener('mouseleave', { .Select() })
		el.AddEventListener('mouseenter', { .Select(el) })
		el.AddEventListener('mousedown',
			{ |event|
			.SELCHANGE(el)
			event.PreventDefault() })

		if .Xmin < w = SuRender().GetTextMetrics(el, s).width + 40/*=extra*/
			{
			.Xmin = w
			.SetMinSize()
			}
		}
	}
