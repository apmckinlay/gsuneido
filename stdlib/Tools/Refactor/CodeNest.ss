// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
// Finds the closest enclosing parenthesis/braces/brackets/quotes to a given position
// Written for Addon_unwrap
// Previously Addon_unwrap would work backwards from the position
// but that meant it couldn't use Scanner
// which meant potential issues
// and it didn't handle if the first delimiter it found working backwards
// wasn't the same as the first delimiter working forwards.
// This version does more work because it uses Scanner from the beginning of the code
// but it handles more scenarios.
// And Scanner is built-in so it is fast.
class
	{
	other: #('(': ')', '[': ']', '{': '}')
	CallClass(code, pos)
		{
		scan = Scanner(code)
		stack = .findPos(scan, pos)
		return .findWrap(scan, pos, stack)
		}
	findPos(scan, pos)
		{
		stack = [[""]] // dummy value to avoid checking for empty
		while scan.Position() < pos and scan isnt type = scan.Next2()
			{
			if type isnt ""
				continue
			tok = scan.Text() // only get text when required (less memory allocation)
			if tok is stack.Last()[0]
				stack.PopLast()
			else if tok in ('(', '[', '{')
				stack.Add([.other[tok], scan.Position()])
			}
		return stack
		}
	findWrap(scan, pos, stack)
		{
		if scan.Type() is #STRING
			{
			tok = scan.Text()
			endPos = scan.Position()
			if endPos > pos
				return [endPos - tok.Size(), endPos-1]
			}
		if stack.Size() > 1
			return .findClosing(scan, stack.Last())
		return false
		}
	findClosing(scan, last)
		{
		startPos = last[1] - 1
		stack = [last[0]] // don't need positions
		while stack.NotEmpty?() and scan isnt type = scan.Next2()
			{
			if type isnt ""
				continue
			tok = scan.Text() // only get text when required (less memory allocation)
			if tok is stack.Last()
				stack.PopLast()
			else if tok in ('(', '[', '{')
				stack.Add(.other[tok])
			}
		endPos = scan.Position() - 1
		return [startPos, endPos]
		}
	}