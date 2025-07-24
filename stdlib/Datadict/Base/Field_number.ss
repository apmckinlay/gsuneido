// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Control: ('Number')
	Format: ('Number')
	Encode(val)
		{
		try
			{
			fmtVal = val
			if String?(fmtVal) and fmtVal.Trim().Prefix?('$')
				fmtVal = fmtVal.Trim().Replace('^\$', '')

			if fmtVal isnt ""
				{
				num = Number(fmtVal)
				if not IsInf?(num)
					return num
				}
			}
		return val
		}
	}