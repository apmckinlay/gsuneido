// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	For: 'LibView'
	Getter_Initialized()
		{ return true /* Assumed initialized unless specified elsewhere */ }

	AddonReady?()
		{
		if not .Initialized
			.Init()
		return .Initialized
		}

	Getter_Editor()
		{ return .Parent.Editor }

	Getter_Explorer()
		{ return .Parent.Explorer }

	Getter_View()
		{ return .Parent.View }

	Getter_Window()
		{ return .Parent.Window }

	Libs()
		{ return .Parent.Libs }
	}
