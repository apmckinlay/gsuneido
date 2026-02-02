// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Container
	{
	New(controls)
		{
		.ctrls = Object()
		.ComponentArgs = Object()
		for (c in controls.Values(list:))
			{
			if c is ''
				continue
			ctrl = .Construct(c)
			.ctrls.Add(ctrl)
			// VirtualList stub controls return false
			if false isnt layout = ctrl.GetLayout()
				.ComponentArgs.Add(layout)
			}
		}

	Tally()
		{ return .ctrls.Size() }
	ctrls: () // default if destroyed
	GetChildren()
		{ return .ctrls }
	Get()
		{
		ob = Object()
		for ctrl in .ctrls
			if ctrl.Method?(#Get)
				ob[ctrl.Name] = ctrl.Get()
		return ob
		}

	Append(control)
		{
		.Insert(.Tally(), control)
		}

	AppendAll(controls)
		{
		for ctrl in controls
			.Append(ctrl)
		}

	Insert(i, control)
		{
		.ActWith()
			{
			.ctrls.Add(ctrl = .Construct(control), at: i)
			DoStartup(ctrl)
			Object('Insert', i, ctrl.GetLayout())
			}
		return ctrl
		}

	Remove(i)
		{
		if not .ctrls.Member?(i)
			return
		.ctrls[i].Destroy()
		.ctrls.Delete(i)
		.Act('Remove', i)
		}

	RemoveAll()
		{
		.ctrls.Each(#Destroy)
		.ctrls = Object()
		.Act('RemoveAll')
		}
	}