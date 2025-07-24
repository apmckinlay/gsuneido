// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
function (node, block)
	{
	stack = Stack()
	stack.Push(node)

	do
		{
		cur = stack.Pop()
		if block(cur) is false
			continue
		for (i = cur.ChildrenSize() - 1; i >= 0; i--)
			stack.Push(cur.Children[i])
		}
	while stack.Count() isnt 0
	}
