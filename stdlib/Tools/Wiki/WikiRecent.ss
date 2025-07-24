// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function (n = 7)
	{
	cutoff = Date().Plus(days: -n)
	body = '<dl>\n'
	lastdate = false
	changed = Object()
	output_changed =
		{
		for page in changed.Members().Sort!()
			body $= '<dd><a href="Wiki?' $ page $ '">' $ page $ '</a></dd>\n'
		changed = Object()
		}
	QueryApply('wiki where edited > ' $ Display(cutoff) $ ' sort edited')
		{ |x|
		if x.edited.NoTime() isnt lastdate
			{
			output_changed()
			body $= '<dt>' $ (lastdate is false ? '' : '<br>') $
				x.edited.Format('ddd. MMM. d, yyyy') $ '</dt>\n'
			lastdate = x.edited.NoTime()
			}
		changed[x.name] = true
		}
	output_changed()
	return	body $ '</dl>\n'
	}