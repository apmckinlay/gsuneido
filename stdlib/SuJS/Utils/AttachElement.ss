// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
function (el, parent, at)
	{
	if parent is false or false is parent = QuerySelector(parent)
		return

	if Object?(at)
		{
		if at.parent.El is parent
			{
			if at.Member?(#parentEl)
				parent = at.parentEl
			at = at.at
			}
		else
			at = false
		}

	if at is false
		parent.AppendChild(el)
	else
		{
		childNodes = parent.childNodes
		if at >= childNodes.length
			parent.AppendChild(el)
		else
			parent.InsertBefore(el, childNodes.item(at))
		}
	}
