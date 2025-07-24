// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(s, block = false)
		{
		ff = new this(s)
		if block isnt false
			return block(ff)
		return ff
		}
	New(.s)
		{
		.ite = 0
		}
	Read(n = false)
		{
		if .ite is .s.Size()
			return false
		if n is false
			n = .s.Size()
		res = .s[.ite::n]
		.ite += res.Size()
		return res
		}
	Readline()
		{
		line = .s[.ite..].BeforeFirst('\n')
		set = 1
		if line[-1] is '\r'
			{
			line = line[..-1]
			set++
			}
		.Seek(.ite + line.Size() + set)
		if line is "" and .Tell() is .s.Size()
			return false
		return line
		}
	Tell()
		{
		return .ite
		}
	Seek(offset, origin = 'set')
		{
		switch (origin)
			{
		case 'set':
			.ite = offset
		case 'end':
			.ite = .s.Size() + offset
		case 'cur':
			.ite += offset
			}
		if .ite < 0
			.ite = 0
		if .ite > .s.Size()
			.ite = .s.Size()
		}
	Write(str)
		{
		.s $= str
		}
	Writeline(str)
		{
		.s $= str $ '\r\n'
		}
	Close()
		{
		}
	Get()
		{
		return .s
		}
	Reset()
		{
		.s = ''
		.ite = 0
		}
	}
