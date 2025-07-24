// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
function (el, parent, at)
	{
	if parent is false or false is parent = QuerySelector(parent)
		return

	if Object?(at)
		at = at.parent.El is parent ? at.at : false

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
