// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_editorTextLimit()
		{
		m = EditorControl.EditorControl_editorTextLimit

		Assert(m(false) is: EditorTextLimit)
		Assert(m(600) is: 600)
		Assert(m(EditorTextLimit + 1) is: EditorTextLimit)
		}
	}
