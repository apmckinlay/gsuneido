// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function(token, children = #(), value = false, position = -1, length = 0)
	{
	if Instance?(token)
		return token
	if children is #() and value is false and position is -1 and length is 0 and
		token isnt TDOPTOKEN.LIST
		if _nodes.Member?(token)
			return _nodes[token]
		else
			return _nodes[token] = TdopSymbol(token)

	if not children.Empty?()
		{
		for child in children
			{
			if child.Position is -1
				continue
			if position is -1
				position = child.Position
			length = child.Position - position + child.Length
			}
		}

	node = TdopSymbol(token, :children, :position, :length)
	if value isnt false
		node.Value = value
	return node
	}