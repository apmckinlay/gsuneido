// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
// This addon handles Toolbar controls / commands which require an .Editor
LibViewAddon
	{
	Commands(cmds)
		{
		cmds.Add(
			#(Comment_Lines,
				'Ctrl+/', 'Comment/Uncomment selected lines', 'CommentLine', seq: 0),
			#(Comment_Selection,
				'Shift+Ctrl+/', 'Comment/Uncomment selected text', 'CommentSpan', seq: 0.1
				),
			#(Flag, "Ctrl+F2", "Mark the current line", seq: 2.0)
			#(Next_Flag, "Ctrl+Shift+F2", "Go to the next marked line", seq: 2.1)
			#(Previous_Flag, "Shift+F2", "Go to the previous marked line", seq: 2.2)
			#(Run, 'F9', 'Run current record or selected text', '!'),
			#(Create_Test_for_Method, 'Ctrl+F12'),
			#(Debug_One_Test, 'Shift+Ctrl+K', '', ''),
			#(Go_To_Associated_Test, 'Ctrl+Shift+T', ''),
			)
		}

	Init()
		{
		if not .Initialized = .Editor isnt false
			return
		.Redir('On_Comment_Lines')
		.Redir('On_Comment_Selection')
		.Redir('On_Flag', .Editor)
		.Redir('On_Next_Flag', .Editor)
		.Redir('On_Previous_Flag', .Editor)
		}

	On_Create_Test_for_Method()
		{
		text = .Editor.Get()
		pos = .Editor.GetSelectionStart()
		if false is method = ClassHelp.MethodName(text, pos)
			return
		name = .CurrentName()
		.On_Go_To_Associated_Test()
		if name is .CurrentName()
			return // test not created
		text = .Editor.Get()
		method = 'Test_' $ method
		if text is GotoLibView.Test_text
			text = 'Test\r\n\t{\r\n\t}'
		else if ClassHelp.Methods(text).Has?(method)
			{
			pos = text.Find(method)
			.Editor.SetSelect(pos)
			return
			}
		test_method = method $ '()\r\n\t\t{\r\n\t\t\r\n\t\t}'
		text = ClassHelp.AddMethodAtEnd(text, test_method)
		.Editor.PasteOverAll(text)
		pos = text.Find(test_method) + test_method.Size() - 5 /*= for cursor position */
		.Editor.SetSelect(pos)
		}

	On_Go_To_Associated_Test()
		{
		name = .CurrentName()
		path = .Explorer.Getpath(.Explorer.GetSelected()).BeforeLast('/')
		.TabMenu_GoToAssociatedTest([:name, :path])
		}

	On_Debug_One_Test()
		{
		.Save()
		name = .CurrentName()
		lib = .CurrentTable()
		x = Global(name)
		if not Class?(x) or not x.Base?(Test)
			return

		ranges = ClassHelp.MethodRanges(.Editor.Get())
		cur = .Editor.GetCurrentPos()
		.Editor.SendToAddons('On_BeforeAllTests')
		if false isnt test_method = .test_method(ranges, cur)
			if LibViewRunTest(.Editor, lib, name, { it.DebugOne(test_method) })
				{
				.Editor.SendToAddons('On_AfterAllTests')
				.AlertInfo('Run Test', name $ '.' $ test_method $ ' SUCCEEDED')
				}
		}

	test_method(ranges, cur)
		{
		test_method = false
		for r in ranges
			if cur > r.from and cur < r.to
				test_method = r.name
		if test_method isnt false and not test_method.Prefix?('Test_')
			test_method = false
		return test_method
		}

	On_Run()
		{
		.Save()

		sel = .Editor.GetSelText().Trim()
		if sel is ""
			.Try_run("all")
		else
			.Try_run("selection", { sel.Eval2() })
		}

	On_ViewRestore_Item_as_of()
		{
		data = .Explorer.Get()
		if data.Member?('text')
			ViewItemControl(.CurrentTable(), .CurrentName(), data.text, .Editor)
		}
	}
