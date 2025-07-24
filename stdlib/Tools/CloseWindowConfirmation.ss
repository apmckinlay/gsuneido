// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: "Close Window Confirmation"
	CallClass(hwnd = 0)
		{
		return ToolDialog(hwnd, this, closeButton?: false)
		}
	Controls:
		(Vert
			(Static	"The information on the current screen is invalid.")
			(Static "Choose 'Correct Changes' to remain on this screen to fix it.")
			(Static "If you choose 'Discard Changes', the window will be closed")
			(Static "and the changes will be lost.")
			Skip
			(Horz Fill (Button 'Correct Changes') Skip (Button 'Discard Changes') Fill)
			xstretch: 0)
	On_Discard_Changes()
		{
		if Suneido.Member?('AccessRecordDestroyed') and
			Suneido.AccessRecordDestroyed isnt ''
			{
			SuneidoLog('Changes discarded by ' $ Suneido.User $
				' (' $ Display(Suneido.AccessRecordDestroyed) $ ')')
			Suneido.AccessRecordDestroyed = ''
			}

		.Window.Result(true)
		}
	On_Correct_Changes()
		{
		Suneido.AccessRecordDestroyed = ''
		.Window.Result(false)
		}
	On_Cancel()
		{
		// disable ESC key
		}
	}