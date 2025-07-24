// Copyright (C) 2022 Suneido Software Corp. All rights reserved worldwide.
// NOTE: Only run with: CLIENT: gsuneido.exe, SERVER: gsport.exe
function ()
	{
	builtin = BuiltinNames().Difference(ServerEval('BuiltinNames'))
	if builtin is Win32BuiltinNames
		Print('Win32BuiltinNames is up to date')
	else
		QueryApply1('stdlib', name: 'Win32BuiltinNames')
			{
			newText = it.text.BeforeFirst('#(') $ '#(\r\n' $ builtin.Join(',\r\n') $ ')'
			SvcTable('stdlib').Update(it, :newText, t: it.Transaction)
			}
	}