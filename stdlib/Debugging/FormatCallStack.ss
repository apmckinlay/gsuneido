// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
// calls is assumed to be in the format of exception.Callstack or Handler
class
	{
	CallClass(calls, levels = 10, indent = false)
		{
		n = 0
		s = ""
		for call in calls
			{
			line = .display(call)
			if line.Has?('WorkSpaceControl.WorkSpaceControl_run')
				break
			line = Unprivatize(line)
			if indent
				line = "   ".Repeat(n) $ line
			s $= line $ "\n"
			if ++n >= levels
				break
			}
		return s[.. -1]
		}
	display(call)
		{
		s = Display(call.fn)
		lib = s.AfterFirst('/* ').BeforeFirst(' ')
		if Libraries().Has?(lib)
			s = lib $ ":" $ s.BeforeFirst(" /*")
		if call.Member?(#srcpos) and
			false isnt src = SourceCode(call.fn)
			s $= ":" $ (1 + src.LineFromPosition(call.srcpos))
		return s
		}
	}
