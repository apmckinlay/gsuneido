// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
EditorControl
	{
	LBUTTONUP()
		{
		.Send("FormulaEditor_Click")
		return 'callsuper'
		}

	Get()
		{
		return super.Get().Trim()
		}
	}