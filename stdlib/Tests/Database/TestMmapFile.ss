// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(secs = 10)
		{
		try Database('drop tmp')
		Database('create tmp (a,b) key(a)')
		(new this)(secs)
		}
	Call(secs)
		{
		.t = Date().Plus(seconds: secs)
		Thread({ .Writer('a') })
		Thread({ .Writer('b') })
		Thread(.Reader)
		Thread(.Reader)
		Thread(.Reader)
		Thread(.Reader)
		Thread(.Dumper)
//		Thread(.Checker)
		}
	Writer(prefix)
		{
		i = 0
		while Date() < .t
			try QueryOutput('tmp', [a: prefix $ i++, b: 'helloworld'.Repeat(200)])
			catch (e) Print(e)
		}
	Reader()
		{
		while Date() < .t
			QueryLast('tmp sort a')
		}
	Dumper() // will also do check
		{
		while Date() < .t
			try Database.Dump()
			catch (e) Print(e)
		}
	Checker()
		{
		while Date() < .t
			try Print(Check: Database.Check())
			catch (e) Print(e)
		}
	}