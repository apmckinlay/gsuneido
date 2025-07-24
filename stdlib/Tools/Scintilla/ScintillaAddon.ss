// Copyright (C) 2009 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	For: 'Editor'
	Getter_IDE()
		{
		return .Parent.IDE
		}
	Getter_Hwnd()
		{ // allow .Hwnd
		return .Parent.Hwnd
		}
	Getter_Window()
		{ // allow .Window
		return .Parent.Window
		}
	}
