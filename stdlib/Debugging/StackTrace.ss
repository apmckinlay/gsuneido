// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
/* Examples:
	StackTrace(10) 								// print to workspace
	s = StackTrace(10, printFn: false) 			// return string
	StackTrace(10, TracePrint)					// print to trace window
	StackTrace(10, { AddFile('1.txt', it) }) 	// print to a text file
*/
function (n = 99, printFn = function (s) { Print(s) }, skip = 1)
	{
	try
		throw #foo
	catch (e)
		{
		s = FormatCallStack(e.Callstack()[skip..], indent:, levels: n)
		if printFn is false
			return s
		printFn(s)
		}
	}
