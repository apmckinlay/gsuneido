// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
Container
	{
	Name: "Form"
	ComponentName: "Form"
	New(@args)
		{
		.ctrls = Object()
		.ComponentArgs = Object()
		for item in args.Values(list:)
			if item in ('nl', 'nL')
				{
				.ctrls.Add('nl')
				.ComponentArgs.Add('nl')
				}
			else if item isnt ''// control
				{
				ctrl = .Construct(item)
				.ctrls.Add(ctrl)
				if Object?(item) and item.Member?("group")
					ctrl.AddExtraSpec(#group, item.group)
				.ComponentArgs.Add(ctrl.GetLayout())
				}
		if args.Member?('left')
			.ComponentArgs.left = args.left
		}

	GetChildren()
		{
		return .ctrls.Filter(Instance?)
		}
	}