// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	return QueryMin(LibViewImportRestoreControl.Table $
		' where lvimport_filename is ' $ Display(.lvimport_filename), 'lvimport_num')
	}