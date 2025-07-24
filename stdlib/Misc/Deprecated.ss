// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
function (name, after = '', use = '')
	{
	SuneidoLog('INFO: ' $
		Join(', ',
			name $ ' is deprecated',
			Join(' ',
				Opt('for BuiltDate() after ', after),
				Opt('use ', use))))
	}