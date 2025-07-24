// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// creates multiple static controls
// '*' toggles bold (unless 'alwaysBold' is passed in true), '_' toggles underline
// e.g. "*Hello* there _world_" would bold "Hello" and underline "world"
HorzControl
	{
	New(@args)
		{
		super(@.controls(args))
		}
	controls(args)
		{
		text = args[0]
		args.Add('Static', at: 0)
		statics = Object()
		bold = args.GetDefault('alwaysBold', false)
		markers = bold ? '_' : '_*'
		underline = false
		while text isnt ""
			{
			i = text.Find1of(markers)
			s = text[.. i]
			if s isnt ''
				{
				c = args.Copy()
				c[1] = s
				c.weight = bold ? 'bold' : 'normal'
				c.underline = underline
				statics.Add(c)
				}
			if text[i] is '*'
				bold = not bold
			else if text[i] is '_'
				underline = not underline
			text = text[i + 1 ..]
			}
		return statics
		}
	}