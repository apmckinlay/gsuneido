// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_KEYDOWN()
		{
		mock = Mock(DrawCanvasControl)
		mock.When.KEYDOWN([anyArgs:]).CallThrough()
		mock.When.Send([anyArgs:]).Return(0)
		mock.When.DeleteSelected([anyArgs:]).Return(0)
		mock.When.MoveSelected([anyArgs:]).Return(0)

		.MakeLibraryRecord([name: #KeyPressed?,
			text: `function(unused) { return _ctrlPressed?}`])

		// Left arrow
		mock.KEYDOWN(37)
		mock.Verify.MoveSelected(-1, 0)

		// Up arrow
		mock.KEYDOWN(38)
		mock.Verify.MoveSelected(0, -1)

		// Right arrow
		mock.KEYDOWN(39)
		mock.Verify.MoveSelected(1, 0)

		// Left arrow
		mock.KEYDOWN(40)
		mock.Verify.MoveSelected(0, 1)

		// Delete key
		mock.KEYDOWN(46)
		mock.Verify.DeleteSelected()

		_ctrlPressed? = false

		// A key
		mock.KEYDOWN(65)
		mock.Verify.Never().ctrlKeys([anyArgs:])

		_ctrlPressed? = true

		// Ctrl + A key
		mock.KEYDOWN(65)
		mock.Verify.Send(#On_Select_All)

		// Ctrl + X key
		mock.KEYDOWN(88)
		mock.Verify.Send(#On_Cut)

		// Ctrl + Y key
		mock.KEYDOWN(86)
		mock.Verify.Send(#On_Paste)

		// Ctrl + uncaught key (default case)
		mock.KEYDOWN(95)
		mock.Verify.ctrlKeys(95)
		mock.Verify.Times(3).Send([anyArgs:])
		}
	}