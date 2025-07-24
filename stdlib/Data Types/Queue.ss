// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
// Author: Victor Schappert
class
	{
	// Unclear whether linked list or circular array is more efficient. This
	// implementation is a circular linked list with a dummy node at the front,
	// so it looks like: [DUMMY NODE] --> FRONT --> ... --> BACK --> [DUMMY NODE]
	node: class
		{
		New(.Value, .Next)
			{ }
		}
	New()
		{
		.back = new .node(false, false)	   // dummy node
		.back.Next = .back                 // circular
		.size  = 0
		}
	Enqueue(@x)
		{
		// allow either Enqueue(value) or Enqueue(member: value, ...)
		if (x.Size() is 1 and x.Member?(0))
			x = x[0]
		newBack = new .node(x, .back.Next)
		.back.Next  = newBack
		.back = newBack
		++.size
		}
	Front()
		{
		if .size < 1
			throw "Queue empty"
		// .back.Next is the dummy value, so .back.Next.Next is the front
		return .back.Next.Next.Value
		}
	Back()
		{
		if .size < 1
			throw "Queue empty"
		return .back.Value
		}
	Dequeue()
		{
		x = .Front()
		--.size
		// Drop the old dummy value, make the previous head the new dummy value.
		dummy = .back.Next
		head  = dummy.Next
		.back.Next = head
		head.Value = false
		return x
		}
	Count()
		{
		return .size
		}
	List()
		{
		list = Object()
		node = .back.Next.Next
		for (i = 0; i < .size; ++i)
			{
			list.Add(node.Value)
			node = node.Next
			}
		return list
		}
	ToString()
		{
		return 'Queue(' $ .List().Join(', ') $ ')'
		}
	}
