// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
// Retrieves a nul terminated string
// Pairs with GlobalAllocString
function (hm)
	{
	return GlobalData(hm).BeforeFirst('\x00')
	}