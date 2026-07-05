// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		return Object(
			// gives errors when comparing across types like `false < ""` or `34 > #()`
			// even though the ordering is well defined, more often than not we want
			// comparisons on the same type
			strictCrossTypeCompares: 'warn',
			// gives errors when two operands of the concat `$` operators are not provably
			// strings, throws a warning/error even if operands auto-coerces
			// eg: 5 + "123"
			strictStringConcat: 'off',
			// only show diagnostics that match this confidence score
			// confidence scores are in the range 0.00 <= score <= 1.00
			// can be passed as #(<, <=, >, >=)
			confidence: '>=0.70',
		)
		}
	}