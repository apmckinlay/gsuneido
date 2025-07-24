// Copyright (C) 2011 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_FromInt()
		{
		.fromInt(0, "0")
		.fromInt(123, "123")
		.fromInt(999, "999")
		.fromInt(1023, "1023")
		.fromInt(1024, "1 kb")
		.fromInt(12345, "12 kb")
		.fromInt(123456, "120 kb")
		.fromInt(1234567, "1.17 mb")
		.fromInt(1234567890, "1.14 gb")
		.fromInt(88504877056, "82.4 gb")
		.fromInt(86430544.Kb(), "82.4 gb")
		.fromInt(36836763641460, '33.5 tb')
		.fromInt(3683676364146000, '3.27 pb')
		.fromInt(3683676364146000000, '3271.76 pb')

		.fromInt(100000, '97.6 kb')
		}
	fromInt(n, s)
		{
		Assert(ReadableSize.FromInt(n) is: s)
		}
	TestToInt()
		{
		.toInt('0', 0)
		.toInt('123', 123)
		.toInt('999', 999)
		.toInt('1023', 1023)
		.toInt('123 Kb', 125952)
		.toInt('123kb', 125952)
		.toInt('456 mb', 478150656)
		.toInt('456MB', 478150656)
		.toInt('789 gb', 847182299136)
		.toInt('789gb', 847182299136)
		.toInt('4.5 mb', 4718592)
		.toInt('4.5mb', 4718592)
		.toInt('4.5tb', 4947802324992)
		.toInt('4.5 tb', 4947802324992)

		.toInt('97.5 kb', 99840)
		}
	toInt(s, n)
		{
		Assert(ReadableSize.ToInt(s) is: n)
		Assert(ReadableSize.FromInt(n).Tr(' ').Lower() is: s.Tr(' ').Lower())
		}
	}