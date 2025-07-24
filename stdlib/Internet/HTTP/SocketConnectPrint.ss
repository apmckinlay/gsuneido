// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
// wraps a socket and prints the data written and read
class
	{
	New(.sc)
		{
		}
	Write(line)
		{
		Print(Write: line)
		.sc.Write(line)
		}
	Writeline(line)
		{
		Print(Writeline: line)
		.sc.Writeline(line)
		}
	Read(n)
		{
		s = .sc.Read(n)
		Print(Read: s)
		return s
		}
	Readline()
		{
		s = .sc.Readline()
		Print(Readline: s)
		return s
		}
	}