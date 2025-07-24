// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass(message, block, quiet? = false, font = '', size = '',
		weight = '', color = '')
		{
		if quiet? or not Sys.GUI?()
			{
			block()
			return
			}

		message = TranslateLanguage(message)
		result = (.dialog)(Object(this, message, block, font, size, weight, color))
		if Type(result) is #Except
			throw result
		return result
		}

	dialog: Dialog
		{
		CallClass(control)
			{
			style = WS.POPUP | WS.DLGFRAME
			w = new this(control, :style, show: false, skipStartup?:)
			return w.InternalRun(.GetParent(), control, keep_size: false, posRect: false)
			}
		}

	New(.message, .block, .font = '', .size = '', .weight = '', .color = '')
		{
		Defer(.runBlock)
		}

	Controls()
		{
		return Object('Border',
			Object('Static', .message, .font, .size, .weight, color:.color),
			20 /*=border*/, xstretch: 0, ystretch: 0)
		}

	runBlock()
		{
		result = true
		try
			(.block)()
		catch (e)
			result = e
		.Window.Result(result)
		}

	// The methods below prevent forcibly escaping the window via shortcut keys
	On_Cancel() 			// Prevents VK.ESCAPE
		{ return 0 }

	Ok_to_CloseWindow?() 	// Sent from Dialog.CLOSE, prevents Alt+F4
		{ return false }
	}
