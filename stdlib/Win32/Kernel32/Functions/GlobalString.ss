// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// Retrieves a nul terminated string
function (hm)
	{
	return GlobalData(hm).BeforeFirst('\x00')
	}