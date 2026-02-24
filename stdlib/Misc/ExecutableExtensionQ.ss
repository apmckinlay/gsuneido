// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
class
	{
	InvalidTypeMsg: "This file type is not allowed"
	CallClass(fileName)
		{
		return fileName.AfterLast('.') in ('exe', 'com', 'bat', 'ps', 'ps1', 'vb', 'vbs',
			'pif', 'reg', 'scr', 'sct', 'vbe', 'wsf', 'wsh', 'zip', 'cmd', 'lnk', 'rgs',
			'msi', 'msp', 'msc', 'jse')
		}
	}