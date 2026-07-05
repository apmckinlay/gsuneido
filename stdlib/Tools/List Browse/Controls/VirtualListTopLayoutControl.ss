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

		.addTotalSelected(extraLayout, layout)
		extraLayout.buttons.Each()
			{ layout.Add(#('Skip', small:), Object?(it) ? it : Object(#Button, it)) }
		return layout
		}

	addTotalSelected(extraLayout, layout)
		{
		if extraLayout.checkBoxAmountField is false
			return

		amtFormat = Object?(extraLayout.checkBoxAmountField)
			? #(Field width: 15, readonly:, name: 'totalSelected')
			: #(Number mask: "-###,###,###.##", width: 15,
				readonly:, name: 'totalSelected')

		layout.Add(Object('Pair'
			#(Static 'Total Selected')
			amtFormat))
		}

	Recv(@args)
		{
		if args.source.Base?(ButtonControl) and args[0].Prefix?('On_')
			.Controller.Send(@args)
		return 0
		}
	}
