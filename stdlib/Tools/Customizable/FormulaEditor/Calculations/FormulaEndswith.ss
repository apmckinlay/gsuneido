// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaContains
	{
	DisplayName: 'ENDSWITH'
	Calc(text, substring)
		{
		return  text.value =~ "(?i)(?q)" $ substring.value $ "(?-q)$"
		}
	}