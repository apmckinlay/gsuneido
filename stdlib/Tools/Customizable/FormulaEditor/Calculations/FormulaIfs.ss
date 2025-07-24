// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
FormulaIf
	{
// suggestion 24047: FormulaIfs was renamed to FormulaIf. This is to handle old formulas
// TEMPORARY: check client status for the message, fix customers, then delete this record
	CallClass(@args)
		{
SuneidoLog('INFO: FormulaIfs in use')
		return super.CallClass(@args)
		}

	Validate(@args)
		{
		return super.Validate(@args)
		}
	}