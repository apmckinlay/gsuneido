// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
PairControl
	{
	New(@args)
		{
		super(@.make_args(args))
		}
	make_args(args)
		{
		size = args.GetDefault(#size, '')
		if Datadict(args[0]).Control[0] is 'CheckBox'
			return Object(
				Object('CheckBox', Prompt(args[0]), name: args[0], :size,
					weight: HeadingControl.Weight),
				#(Static ''))
		heading = Object('Heading', Prompt(args[0]), size)
		weight = args.GetDefault(#weight, 'semibold')
		field = args.Add('NoPrompt', at: 0).Merge([:size, :weight])
		return Object(heading, field)
		}
	}
