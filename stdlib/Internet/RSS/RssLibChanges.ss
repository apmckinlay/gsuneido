// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
function (@unused)
	{
	r = new Rss2
	r.Channel(title: 'Stdlib Changes', link: 'abc.com', description: '')
	n = 0
	QueryApply('stdlib sort reverse lib_modified')
		{ |x|
		r.AddItem(title: x.name, description: Display(x.lib_modified))
		if ++n >= 10
			break
		}
	return r.ToString()
	}