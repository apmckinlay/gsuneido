// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(@args)
		{
		.stack = Object()
		args.Each(.Push)
		}
	Push(@x)
		{
		// allow either Push(value) or Push(member: value, ...)
		if x.Size() is 1 and x.Member?(0)
			x = x[0]
		.stack.Add(x)
		}
	Top(i = 0)
		{
		if i >= .stack.Size()
			throw "Stack underflow"
		return .stack[.stack.Size() - 1 - i]
		}
	Pop()
		{
		x = .Top()
		.stack.Delete(.stack.Size() - 1)
		return x
		}
	Count()
		{
		return .stack.Size()
		}
	List()
		{
		return .stack
		}
	ToString()
		{
		return 'Stack(' $ .stack.Join(', ') $ ')'
		}
	}