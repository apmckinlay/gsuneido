// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
PassthruController
	{
	ComponentName: 'ShowOne'
	New(@ctrls)
		{
		.i = 0
		.ctrls = Object()
		.data_ctrls = Object().Set_default(#())
		.ComponentArgs = Object()
		for (ind in ctrls.Members(list:))
			{
			c = ctrls[ind]
			.ctrls.Add(ctrl = .Construct(c))
			++.i
			.ComponentArgs.Add(ctrl.GetLayout())
			}
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
		.Act(#UpdateShow, .show)
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
	Destroy()
		{
		.Send("RemoveSetObserver", .handleRecordSet)
		super.Destroy()
		}
	}