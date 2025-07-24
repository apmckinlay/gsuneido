// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (book, option)
	{
	// WARNING: assumes text does not start or end with whitespace
	return book $ '
		where path !~ `^/res\>`
		where (text is ' $ Display(option) $ ' or text is ' $ Display(option $ '()') $ ')
		project path, name'
	}
