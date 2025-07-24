// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
function (value)
	{
	return Type(value) in ('Class','Instance') ? value.Base() : false
	}