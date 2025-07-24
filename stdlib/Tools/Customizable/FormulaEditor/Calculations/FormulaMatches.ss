// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaContains
	{
	DisplayName: 'MATCHES'
	Calc(text, substring)
		{
		return  text.value =~ substring.value
		}
	}