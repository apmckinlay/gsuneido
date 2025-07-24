// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
// TODO: make types plugins
ChooseListControl
	{
	New(@args)
		{
		super(@.make_args(args))
		}
	make_args(args)
		{
		types = CustomFieldTypes(args.GetDefault(#reporter, false),
			filterBy: args.GetDefault(#filterBy, false))
		typenames = types.Map!({ it.name })
		args.Add(typenames, at: 0)
		args.width = 12
		return args
		}
	}