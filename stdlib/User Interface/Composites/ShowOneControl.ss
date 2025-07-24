// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
PassthruController
	{
	New(@ctrls)
		{
		.i = 0
		.ctrls = Object()
		.data_ctrls = Object().Set_default(#())
		for (ind in ctrls.Members(list:))
			{
			c = ctrls[ind]
			.ctrls.Add(.Construct(c))
			++.i
			}
		.Recalc()
		.Send("AddSetObserver", .handleRecordSet)
		.handleRecordSet()
		}
	handleRecordSet()
		{
		.Defer(.update, uniqueID: 'recordSet')
		}
	Data(source)
		{
		.data_ctrls[.i].Add(source)
		.Send('Data', :source)
		}
	Recalc()
		{
		// set the field dimensions to that of the currently visible field
		i = .find_non_empty()
		.Xmin = .ctrls[i].Xmin
		.Ymin = .ctrls[i].Ymin
		.Top = .ctrls[i].Top
		.Xstretch = .ctrls[i].Xstretch
		.Ystretch = .ctrls[i].Ystretch
		}

	show: false
	SetVisible(visible)
		{
		for (i = 0; i < .ctrls.Size(); ++i)
			.ctrls[i].SetVisible(visible and i is .show)
		}
	update()
		{
		i = .find_non_empty()
		if i is .show // no change
			return
		.show = i
		.SetVisible(true)
		}
	find_non_empty()
		{
		for (i = .data_ctrls.Size() - 1; i > 0; --i)
			for c in .data_ctrls[i]
				if c.Get() isnt ''
					return i
		return 0
		}
	GetChildren()
		{
		return .ctrls
		}
	Resize(x, y, w, h)
		{
		for c in .ctrls
			c.Resize(x, y + .Top - c.Top, w, h)
		}
	Destroy()
		{
		.Send("RemoveSetObserver", .handleRecordSet)
		super.Destroy()
		}
	}