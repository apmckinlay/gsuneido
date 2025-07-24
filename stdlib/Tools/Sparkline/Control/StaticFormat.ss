// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
TextFormat
	{
	New(@args)
		{
		super(@.convert(args))
		}
	convert(args)
		{
		if args.Member?('font')
			{
			font = args.font
			args.font = Object?(font) ? font : Object(name: font)
			if args.Member?('size')
				{
				args.font = args.font.Copy()
				if Number?(args.size)
					args.font.size = args.size
				}
			}
		if args.Member?('justify')
			args.justify = args.justify.Lower()
		return args
		}
	}