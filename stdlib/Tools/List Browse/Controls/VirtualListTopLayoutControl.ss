// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New(.extraLayout)
		{
		}

	Controls()
		{
		return .BuildLayout(.extraLayout)
		}

	BuildLayout(extraLayout, layout = false)
		{
		if layout is false
			layout = Object('HorzEqual', 'Fill')

		if extraLayout.checkBoxAmountField isnt false
			layout.Add(#(Pair
				(Static 'Total Selected')
				(Number mask: "-###,###,###.##", readonly:, name: 'totalSelected')))
		extraLayout.buttons.Each()
			{ layout.Add(#('Skip', small:), Object?(it) ? it : Object(#Button, it)) }
		return layout
		}

	Recv(@args)
		{
		if args.source.Base?(ButtonControl) and args[0].Prefix?('On_')
			.Controller.Send(@args)
		return 0
		}
	}
