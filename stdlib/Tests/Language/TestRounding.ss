// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		File("rounding.txt", 'w')
			{|f|
			for .. 1000000
				{
				n = .random()
				digits = Random(8)
				mask = '################.' $ '#'.Repeat(digits)
				f.Writeline(n $ '\t' $ digits $ '\t' $ n.Round(digits))
				if Number(n.Format(mask)) isnt n.Round(digits)
					Print(n, digits, 'Format', n.Format(mask), 'Round', n.Round(digits))
				}
			}
		}
	CheckFile()
		{
		nlines = 0
		File("/dev/suneido/rounding.txt", 'r')
			{|f|
			while false isnt line = f.Readline()
				{
				++nlines
				split = line.Split('\t')
				n = Number(split[0])
				digits = Number(split[1])
				expected = Number(split[2])
				if n.Round(digits) isnt expected
					Print(n, digits, 'expected', expected, "got", n.Round(digits))
				mask = '################.' $ '#'.Repeat(digits)
				if Number(n.Format(mask)) isnt n.Round(digits)
					Print(n, digits, 'Format', n.Format(mask), 'Round', n.Round(digits))
				}
			}
		Print(nlines)
		}
	random()
		{
		ndigits = Random(16) + 1
		s = ''
		for i in .. ndigits
			s $= Random(10)
		i = Random(ndigits + 1)
		s = s[.. i] $ '.' $ s[i ..]
		return Number(s)
		}
	}