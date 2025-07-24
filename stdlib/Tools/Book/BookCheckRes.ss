// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
function (book)
	{
	Working('Checking res usage')
		{
		QueryApply(book $ ' where path =~ "^/res\>" sort name')
			{ |x|
			if QueryEmpty?(book $ ' where text =~ ' $ Display(x.name))
				Print('Image not used ' $ Display(x.name))
			}
		}
	}
