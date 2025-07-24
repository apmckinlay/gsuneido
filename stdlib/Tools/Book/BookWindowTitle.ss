// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
function (book, option = '')
	{
	titleFunc = OptContribution('BookWindowTitle', function(@unused){ return false })
	return titleFunc(book, option)
	}