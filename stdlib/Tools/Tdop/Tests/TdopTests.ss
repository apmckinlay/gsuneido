// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	CheckTdop(src, expected, expectedPositions = false, type = 'statements')
		{
		res = Tdop(src, :type)
		Assert(Display(res) like: expected)
		if expectedPositions isnt false
			{
			positions = .getProperty(res, #Position)
			Assert(positions is: expectedPositions)
			}
		.checkLength(res)
		}

	getProperty(token, property)
		{
		if token.ChildrenSize() is 0
			return token[property]

		res = Object(token[property])
		for child in token.Children.Members().Sort!()
			res.Add(.getProperty(token.Children[child], property))
		return res
		}

	checkLength(token)
		{
		if token.ChildrenSize() is 0
			{
			if token.Position is -1
				Assert(token.Length is: 0)
			else
				Assert(token.Length isnt: 0)
			return
			}
		position = -1
		length = 0
		for child in token.Children
			{
			.checkLength(child)
			if child.Position is -1
				continue
			if position is -1
				position = child.Position
			length = child.Position - position + child.Length
			}
		Assert(token.Position is: position)
		Assert(token.Length is: length)
		}

	CheckTdopCatch(src, error = false, type = 'statements')
		{
		if error is false
			Assert({ Tdop(src, :type) } throws:)
		else
			Assert({ Tdop(src, :type) } throws: error)
		}
	}