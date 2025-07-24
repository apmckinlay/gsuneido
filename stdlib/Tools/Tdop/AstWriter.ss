// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.source, .astTree, .comments, .blanks, .startNode = false)
		{
		.writerId = Timestamp()
		.init()
		}

	init()
		{
		.offsets = Object()
		.length = .startNode is false ? .source.Size() : .startNode.Length
		.content = .startNode is false
			? .source
			: .source[.startNode.Position-1::.startNode.Length]
		.initOffeset = .startNode is false
			? 0
			: .startNode.Position - 1
		}

	Length()
		{
		return .length
		}

	updateLength(change)
		{
		.length += change
		}

	Add(node, index, value)
		{
		if not node.Match(TDOPTOKEN.LIST)
			throw 'ADD: target must be LIST'

		.addEvent(node, #ADD, index, value)
		}

	Replace(node, value)
		{
		.addEvent(node, #REPLACE, false, value)
		}

	Remove(node)
		{
		.addEvent(node, #REMOVE, false, false)
		}

	addEvent(node, op, index, value)
		{
		event = Object(:node, :op, :index, :value)
		eventStore = node.GetInit(#events, Object()).GetInit(.writerId, Object())
		switch (op)
			{
		case #REMOVE:
			eventStore.Delete(all:).Add(event)
		case #REPLACE:
			if not eventStore.Empty?() and eventStore[0].op is #REMOVE
				throw 'REPLACE: target node removed'
			eventStore.Delete(all:).Add(event)
		case #ADD:
			if not eventStore.Empty?() and eventStore[0].op is #REMOVE
				throw 'ADD: target node removed'
			eventStore.Add(event)
			}
		}

	ToString()
		{
		.init()
		node = .startNode is false
			? .GetRoot()
			: .startNode
		TdopTraverse2(node, .handleEvent)
		return .content
		}

	handleEvent(node)
		{
		eventStore = node.GetDefault(#events, Object()).GetDefault(.writerId, Object())
		continueTraverse = true
		for event in eventStore
			{
			switch (event.op)
				{
			case 'ADD':
				.handleAdd(event)
			case 'REPLACE':
				continueTraverse = false
				.handleReplace(event)
			case 'REMOVE':
				continueTraverse = false
				.handleRemove(event)
				}
			}
		return continueTraverse
		}

	handleAdd(event)
		{
		node = event.node
		index = event.index
		value = event.value

		oriStart = index isnt node.ChildrenSize()
			? node.Children[index].Position - 1
			: node.Position - 1 + node.Length
		start = .GetCurPos(oriStart)
		.content = .content[..start] $ value $ .content[start..]
		change = String?(value) ? value.Size() : value.Length()
		.updateOffsets(oriStart, oriStart, change)
		.updateLength(change)
		}

	handleReplace(event)
		{
		node = event.node
		value = event.value

		oriStart = node.Position - 1
		oriEnd = oriStart + node.Length - 1
		start = .GetCurPos(oriStart)
		end = .GetCurPos(oriEnd)

		.content = .content[..start] $ value $ .content[end + 1..]
		change = (String?(value) ? value.Size() : value.Length()) - (end - start + 1)
		.updateOffsets(oriStart, oriEnd + 1, change)
		.updateLength(change)
		}

	handleRemove(event)
		{
		node = event.node

		oriStart =node.Position - 1
		oriEnd = oriStart + node.Length - 1
		start = .GetCurPos(oriStart)
		end = .GetCurPos(oriEnd)

		.content = .content[..start] $ .content[end + 1..]
		change = - (end - start + 1)
		.updateOffsets(oriStart, oriEnd + 1, change)
		.updateLength(change)
		}

	updateOffsets(start, end, length)
		{
		.offsets[end] = .offsets.GetDefault(end, 0) + length
		if start < end
			for offsetPos in .offsets.Copy().Members()
				if offsetPos > start and offsetPos < end
					.offsets.Delete(offsetPos)
		}

	GetCurPos(pos)
		{
		curPos = pos - .initOffeset
		for offsetPos in .offsets.Members()
			if offsetPos <= pos
				curPos += .offsets[offsetPos]
		return curPos
		}

	GetRoot()
		{
		return .astTree
		}

	GetSource()
		{
		return .source
		}

	ReWrite(node = false, offset = 0, skipTail = false)
		{
		s = ''
		if node is false
			node = .GetRoot()
		iString = s.Size() + offset + 1
		TdopTraverse(node)
			{ |node|
			if node.Position isnt -1 and node.ChildrenSize() is 0
				{
				while iString < node.Position
					if .comments.Member?(iString)
						{
						s $= .comments[iString].text
						iString = s.Size() + offset + 1
						}
					else if .blanks.Member?(iString)
						{
						s $= .blanks[iString].text
						iString = s.Size() + offset + 1
						}
					else
						throw "mismatch"

				if iString isnt node.Position
					throw "mismatch"

				s $= node.ToWrite()
				iString = s.Size() + offset + 1
				}
			node.Position isnt -1
			}

		stop = false
		while stop is false and skipTail is false
			{
			stop = true
			if .comments.Member?(iString)
				{
				s $= .comments[iString].text
				iString = s.Size() + 1
				stop = false
				}
			else if .blanks.Member?(iString)
				{
				s $= .blanks[iString].text
				iString = s.Size() + 1
				stop = false
				}
			}
		return s
		}
	}