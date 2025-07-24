// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Addon
	{
	For: 'OpenImage|OpenImageWithLables'

	State(disabled? = false)
		{
		return disabled? ? MFS.DISABLED : MFS.ENABLED
		}

	Getter_Window()
		{ // allow .Window
		return .Parent.Window
		}

	FileEmpty?()
		{
		return .Parent.File.Blank?()
		}
	}
