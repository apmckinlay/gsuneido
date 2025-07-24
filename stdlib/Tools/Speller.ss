// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
// Spell checking and correction suggestions
// Test with: Speller("recieve")
// Uses hunspell
// Keeps the hunspell program active after first use to avoid startup cost.
class
	{
	CallClass(word)
		{
		try
			return .checkSpeller(word)
		catch (err)
			return .errorHandling(err)
		}

	errorHandling(err)
		{
		.Close()
		if err.Prefix?("RunPiped:")
			SuneidoLog("INFO: (CAUGHT) Speller: " $ err, caughtMsg: 'no spell checker')
		else if err.Prefix?('FileExists?:')
			AlertDelayed("Could not access spell checker from " $ ApplicationDir() $
				', this is possibly caused by network issues\r\n\r\n' $
				'Please contact your system administrator')
		else
			throw err
		return #()
		}

	checkSpeller(word)
		{
		suggestions = #()
		if not .Open() or Suneido.Speller is false or word is ""
			return suggestions

		Suneido.Speller.Writeline('^' $ word)
		Suneido.Speller.Flush()
		do
			{
			if false is reply = Suneido.Speller.Readline()
				return suggestions
			}
			while reply[0] not in ('*', '+', '-', '#', '&')

		reply = reply.Trim()
		if reply isnt "*" and reply isnt ""
			suggestions = reply.Trim().AfterFirst(': ').Split(', ')

		return suggestions
		}
	Open()
		{
		if Suneido.Member?(#Speller)
			return true // already open
		if "" is result = .open()
			return true
		if Suneido.User is 'default'
			Print("WARNING: " $ result)
		else
			SuneidoLog("ERRATIC: " $ result)
		return false
		}
	open()
		{
		Suneido.Speller = false
		if Sys.Linux?()
			prog = 'hunspell'
		else
			if false is prog = ExternalApp("hunspell")
				return "can't find hunspell"

		dict = Suneido.Language.GetDefault('dict', false)
		if dict is false
			return "can't find Suneido.Language.dict"
		if not Sys.Linux?()
			{
			dict = Paths.ParentOf(prog) $ "/" $ dict
			if not FileExists?(file = dict $ ".aff")
				return "can't find " $ file
			if not FileExists?(file = dict $ ".dic")
				return "can't find " $ file
			}

		cmdline = '"' $ prog $ '" -d "' $ dict $ '"'
		try
			Suneido.Speller = RunPiped(cmdline)
		catch (e)
			return "from RunPiped(" $ cmdline $ ") => " $ e
		Suneido.Speller.Readline()
		return ""
		}
	Close() // not called automatically
		{
		if false is Suneido.Speller
			return
		try
			Suneido.Speller.Close()
		Suneido.Delete(#Speller)
		}
	}
