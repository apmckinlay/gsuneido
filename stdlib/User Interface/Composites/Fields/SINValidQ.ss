// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(num, checksum? = true)
		{
		if num is ""
			return true

		if num !~ SINControl.Pattern
			return false

		num = num.Tr('^0-9')
		if (not num.Number?() or num.Size() isnt 9) /* = SIN: 9 digits*/
			return false
		if not checksum?
			return true
		.validateLastDigit(num)
		}

	validateLastDigit(num)
		{
		sum = 0
		for (i = 0 ; i < num.Size()-1 ; ++i)
			{
			if (i % 2 is 0)
				sum += Number(num[i])
			else
				{
				s = String(Number(num[i]) * 2)
				for (j=0; j < s.Size(); ++j)
					sum += Number(s[j])
				}
			}
		checkdigit = ((((sum/10).Int() + 1) * 10) - sum) % 10 /*= decimal*/
		lastdigit = Number(num[-1])
		return checkdigit is lastdigit
		}
	}