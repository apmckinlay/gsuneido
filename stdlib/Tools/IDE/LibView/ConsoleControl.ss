// Copyright (C) 2001 Suneido Software All rights reserved worldwide.

// 	 begun 20010724 [vcs]
// Purpose:	generic output window

Controller
	{
	/* purpose: allows an output window to be displayed */
	// data:
	Name:	"Console"
	Title:	"Console"
	Xmin:	200
	Ymin:	100
	Commands:
		(
		(Clear,	"",	"", Delete)
		(Break, "", "", Flag)
		)
	Controls:
		(Vert
			(Toolbar Clear Break)
			(Scintilla readonly:)
		)
	// interface:
	New()
		// post:	returns a new instance of ConsoleControl
		{
		.messages = .Vert.Editor
		.messages.SetFont(StdFonts.Mono())
		if not Suneido.Member?("Console") or Suneido.Console is false
			Suneido.Console = this
		}
	Destroy()
		// post:	Suneido.Console does not refer to this
		{
		if Suneido.Member?("Console") and Suneido.Console is this
			Suneido.Console = false
		super.Destroy()
		}
	Get()
		// post:	returns current messages text
		{ return .messages.Get() }
	Set(text)
		// pre:	text is a string
		// post:	sets messages text to text
		{
		.messages.SetReadOnly(false)
		.messages.Set(text)
		.messages.SetReadOnly(true)
		}
	Append(text)
		// pre:	text is a string
		// post:	text is appended to end of this' text
		{
		.messages.AppendText(text)
		.messages.Update()
		}
	Clear()
		// post:	this contains no text
		{ .Set("") }
	On_Clear()
		{ .Clear() }
	On_Break()
		{ .Append("-------------------------\n") }
	}