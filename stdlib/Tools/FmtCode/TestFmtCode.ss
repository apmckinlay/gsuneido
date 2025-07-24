// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
function (lib = #stdlib)
	{
	n = 0
	QueryApply(lib $ " where group is -1 sort name")
		{|x|
		++n
		try
			{
			v = Suneido.Compile(x.text, oldwarn = Object())
			if Type(v) not in (#Function, #Class)
				continue
			}
		catch
			continue
		try
			{
			newtext = FmtCode(x.text)
			Suneido.Compile(newtext, newwarn = Object())
			Assert(oldwarn.Size(), is: newwarn.Size())
			}
		catch (e)
			{
			Print(n, x.name, e)
			break
			}
		}
	}