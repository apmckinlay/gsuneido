// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(node, preVisit, postVisit = false, reverse = false)
		{
		stack = Stack()
		stack.Push(Object(type: 'PRE', :node))

		do
			{
			cur = stack.Pop()
			switch (cur.type)
				{
			case 'PRE':
				stack.Push(Object(type: 'POST', node: cur.node))
				if preVisit(cur.node) is false
					continue
				.PushChildren(reverse, cur, stack)

			case 'POST':
				if postVisit isnt false
					postVisit(cur.node)
				}
			}
		while stack.Count() isnt 0
		}
	PushChildren(reverse, cur, stack)
		{
		if reverse is false
			for (i = cur.node.ChildrenSize() - 1; i >= 0; i--)
				stack.Push(Object(type: 'PRE', node: cur.node.Children[i]))
		else
			for (i = 0; i < cur.node.ChildrenSize(); i++)
				stack.Push(Object(type: 'PRE', node: cur.node.Children[i]))
		}
	}