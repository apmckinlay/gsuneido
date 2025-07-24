// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(hwnd, parent)
		{
		OkCancel(Object(this, parent), 'Import Images', hwnd)
		}

	New(parent)
		{
		super(.layout(parent))
		}

	layout(parent)
		{
		return Object('Vert'
			Object('Static' 'Import into "' $ parent $ '"?', textStyle: 'main')
			#(Static '(clicking "Cancel" will import into "/res")', textStyle: 'main'),
			xstretch: 0)
		}
	}
