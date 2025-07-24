// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
class
	{
	Check(ob, maxLevel)
		{
		stack = Stack()
		.check(ob, 1, stack, maxLevel)
		}

	check(ob, level, stack, maxLevel)
		{
		if level >= maxLevel
			SuneidoLog('FlatObject.Check', params: Object(stack: stack.Stack_stack, :ob))
		else
			{
			for m in ob.Members()
				{
				if not Object?(ob[m])
					continue
				stack.Push(ob[0])
				.check(ob[m], level + 1, stack, maxLevel)
				stack.Pop()
				}
			}
		}
	}