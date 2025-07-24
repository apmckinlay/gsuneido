// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function (address = false)
	{
	ctrl = GotoPersistentWindow(TranslateLanguage("User's Manual"),
		#(Book, suneidoc, "User's Manual", help_book:), title?: true)

	if address isnt false
		ctrl.Goto(address)
	try SetActiveWindow(ctrl.Window.Hwnd)
	}