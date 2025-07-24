// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_showWarnings?()
		{
		mock = Mock(Addon_libview_todo)
		mock.When.showWarnings?([anyArgs:]).CallThrough()
		mock.When.todoHeight([anyArgs:]).Return(0)
		mock.Addon_libview_todo_prevWarningOb = #()

		// '' will not show warnings
		Assert(mock.showWarnings?('') is: false)

		// Warnings have a new error
		Assert(mock.showWarnings?('Warning: line 48 is ill formed'))

		// Error has moved a line, not a new warning
		Assert(mock.showWarnings?(text = 'Warning: line 48 is ill formed') is: false)

		// Error message changed for said line, new warning
		Assert(mock.showWarnings?(text = 'Warning: line 48 is broken'))

		// Additional error has been added, new warning
		Assert(mock.showWarnings?('ERROR: another line is broken\n' $ text))

		// One of the errors was fixed, no new warning
		Assert(mock.showWarnings?(text) is: false)

		// Same error with empty lines and duplicate message, no new warning
		Assert(mock.showWarnings?(text $ '\n \n\n\n\t' $ text) is: false)

		// Clear out the warnings, text is seen as a new warning
		// todoHeight is greater than 0 though, so we do not need to expand the todo
		mock.Addon_libview_todo_prevWarningOb = #()
		mock.When.todoHeight([anyArgs:]).Return(10)
		Assert(mock.showWarnings?(text $ '\n \n\n\n\t' $ text) is: false)

		// Repeat the process to verify the prior statement
		// Now that the warnings are cleared out again and the todoHeight is 0,
		// the method will return true (indicating we need to expand the panel)
		mock.Addon_libview_todo_prevWarningOb = #()
		mock.When.todoHeight([anyArgs:]).Return(0)
		Assert(mock.showWarnings?(text $ '\n \n\n\n\t' $ text))
		}
	}
