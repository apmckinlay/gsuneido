// Copyright (C) 2013 Axon Development Corporation All rights reserved worldwide.
Hwnd
	{
	New(text)
		{
		text = TranslateLanguage(text)
		.CreateWindow( "button", text, WS.VISIBLE | BS.GROUPBOX, WS_EX.TRANSPARENT)
		.SetFont(text: text $ ' ')
		.TextHeight = .Ymin
		}
	SetReadOnly(readOnly/*unused*/)	// override Hwnd: read-only not applicable to groupbox
		{ }
	GetReadOnly() // override Hwnd: read-only not applicable to groupbox
		{ return true }
	}