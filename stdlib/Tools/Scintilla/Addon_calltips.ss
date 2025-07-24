// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
ScintillaAddon
	{
	ContextMenu()
		{
		return #('Show Parameters\tCtrl+F11')
		}
	On_Show_Parameters()
		{
		.WordRight()
		.CharRight()
		.start_calltip()
		}
	CharAdded(c)
		{
		if .CallTipActive() is 1
			{
			if c is '('
				++.parens
			else if c is ')' and --.parens <= 0
				.end_calltip()
			}
		else if c is '('
			{
			.parens = 1
			.start_calltip()
			}
		}
	start_calltip()
		{
		if false is word = .GetCurrentReference()
			return
		word = word.RightTrim('.')
		if false is params = .params(word)
			return
		SendMessageTextIn(.Hwnd, SCI.CALLTIPSHOW,
			.GetCurrentPos() - word.Size() - 1,	word $ params)
		}
	params(word)
		{
		if word.Prefix?('Suneido.')
			return false
		try
			value = Global(word)
		catch
			return false
		try
			return .checkForUsage(value)
		try
			{
			if false isnt c = value.MethodClass('CallClass')
				return .checkForUsage(c.CallClass)
			if false isnt c = value.MethodClass('New')
				return .checkForUsage(c.New)
			}
		return false
		}
	// Standard Usage:
	//		EX: CallClass(@args) /*usage: [ 0: val, 1: field, ... ]*/
	// 		Things to note:
	//			- Works with: functions, CallClass, New, public methods of classes
	//			- args doesn't need to be named args (args is standard though)
	//			- the usage doesn't need to be on the same line as @args
	//			-- However, it must come before the first { though
	//			- the usage can contain line breaks in code
	//			-- However, it will be displayed as one line in the tool tip
	checkForUsage(call)
		{
		params = call.Params()
		if not params.Has?('@')
			return params

		c = SourceCode(call)
		rxCall = Type(call) is 'Method'
			? Display(call).AfterFirst('.').BeforeFirst(' ')
			: 'function'
		rxUsage = rxCall $ `(?i)[[:blank:]]?(?q)` $ params $ `(?-q)[[:space:]]*/\*usage`
		if c.Size() is i = c.FindRx(rxUsage)
			return params
		usage = c[i ..].AfterFirst(`/*usage:`).BeforeFirst(`*/`)
		return params $ ' - usage: ' $ usage.Tr('\r\n', ' ').Tr('\t').Trim()
		}
	end_calltip()
		{
		if .CallTipActive() is 1
			.CallTipCancel()
		}
	}