// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBeforeAfterBase
	{
	DisplayName: 'AFTERLAST'
	Calc(s, delimiter)
		{
		return s.value.AfterLast(delimiter.value)
		}
	}
