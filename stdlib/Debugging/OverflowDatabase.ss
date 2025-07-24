// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	Database("ensure big (k,s) key(k)")
	s = "1234567890".Repeat(1000) /*= data size*/
	for (i = 0; i < 2000; ++i) /*= number of records */
		QueryOutput('big', Record(k: Timestamp(), :s))
	}