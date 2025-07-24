// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function (source)
	{
	ctrl = source.Ctrl
	editor = ctrl.FindControl('Editor')
	if editor is false or not editor.Method?('GetCurrentWord')
		{
		GotoUserManual()
		return
		}
	selection = editor.GetCurrentWord()
	editor.SelectCurrentWord()
	match = QueryLast('suneidoc' $
		' where name is ' $ Display(selection) $ ' and path !~ "^/res\>"' $
		' sort path')
	GotoUserManual(match isnt false ? match.path $ '/' $ match.name : '')
	}