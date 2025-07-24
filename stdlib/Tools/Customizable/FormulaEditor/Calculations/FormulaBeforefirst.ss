// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaBeforeAfterBase
	{
	DisplayName: 'BEFOREFIRST'
	Calc(s, delimiter)
		{
		return s.value.BeforeFirst(delimiter.value)
		}
	}
