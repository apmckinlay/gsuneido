// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	return TableExists?('users') and not QueryEmpty?('users')
	}