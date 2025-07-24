// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getDiffs()
		{
		list1 = 'line1\nlineA\nline3'
		list2 = 'line1\nlineB\nline3'
		base = 'line1\nline2\nline3'
		fn = DiffBase.GetDiffs
		result = fn(list1, list2, base)
		Assert(result hasMember: 'diffs')
		Assert(result hasMember: 'model')
		Assert(result.diffs is: #(#("line1", "", "line1"),
			#("line2", "-", ""),
			#("", "+1", "lineB"),
			#("", "+2", "lineA"),
			#("line3", "", "line3")))
		}
	Test_normalizeLineEnds()
		{
		fn = DiffBase.DiffBase_normalizeLineEnds
		Assert(fn("") is: "")
		Assert(fn("\r") is: "\r\n")
		Assert(fn("\r\n") is: "\r\n")
		Assert(fn("\n") is: "\n")
		Assert(fn("abc") is: "abc")
		Assert(fn("\r\n\ta\r\nb\r\n") is: "\r\n\ta\r\nb\r\n")
		Assert(fn("\r\n\ta\rb\r\n") is: "\r\n\ta\r\nb\r\n")
		Assert(fn("\r\n\ta\rb\r") is: "\r\n\ta\r\nb\r\n")
		Assert(fn("\r\ta\rb\r") is: "\r\n\ta\r\nb\r\n")
		Assert(fn("\r\ta\r\rb\r") is: "\r\n\ta\r\n\r\nb\r\n")
		Assert(fn("\r\ta\r\rb\r\r") is: "\r\n\ta\r\n\r\nb\r\n\r\n")
		}
	}