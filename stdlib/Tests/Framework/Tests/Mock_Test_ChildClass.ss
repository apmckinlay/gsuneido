// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
// used for Mock_Test
Mock_Test_Class
	{
	m: 100
	M: 4
	Bar()
		{
		return .Foo() + .m
		}
	bar()
		{
		return .Foo() * .M
		}

	M2()
		{
		.m3()
		}
	m3()
		{
		throw 'should not be called here'
		}
	}