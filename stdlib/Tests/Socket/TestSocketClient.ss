// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass()
		{
		client = new this
		client.Run()
		}
	Run()
		{
		.small = "helloworld".Repeat(10) /*= small string: 100 */
		.large = "helloworld".Repeat(10000) /*= large string: 100,000 */
		.line = .small
		.longline = "helloworld".Repeat(400) /*= 4000 = max Readline */
		Print('with same connection')
		SocketClient('127.0.0.1', TestSocketServer.Port)
			{|sc|
			.cases(sc)
			sc.Writeline('during 100')
			.cases(sc)
			sc.Writeline('before 100')
			.cases(sc)
			}
		Print('with separate connections')
		.closed({ .readline(it, 'line', .line) })
		.closed({ .readline(it, 'longline', .longline) })
		.closed({ .read1(it, 'small', .small) })
		.closed({ .read4(it, 'small', .small) })
		.closed({ .read1(it, 'large', .large) })
		.closed({ .read4(it, 'large', .large) })
		for chunksize in #(7, 31, 201, 1024)
			{
			.closed({ .readchunked(it, 'small', .small, chunksize) })
			.closed({ .readchunked(it, 'large', .large, chunksize) })
			}
		Print('DONE')
		}
	cases(sc)
		{
		.readline(sc, 'line', .line)
		.readline(sc, 'longline', .longline)
		.reads(sc, 'small', .small)
		.reads(sc, 'large', .large)
		}
	readline(sc, request, expected)
		{
		sc.Writeline(request)
		actual = sc.Readline()
		Assert(actual.Size() is: expected.Size())
		Assert(actual is: expected)
		}
	reads(sc, request, expected)
		{
		.read1(sc, request, expected)
		.read4(sc, request, expected)
		}
	read1(sc, request, expected)
		{
		sc.Writeline(request)
		actual = sc.Read(expected.Size())
		Assert(actual.Size() is: expected.Size())
		Assert(actual is: expected)
		}
	read4(sc, request, expected)
		{
		quarter = expected.Size() / 4 /*= quarters */
		sc.Writeline(request)
		actual = sc.Read(quarter)
		actual $= sc.Read(quarter)
		actual $= sc.Read(quarter)
		actual $= sc.Read(quarter)
		Assert(actual.Size() is: expected.Size())
		Assert(actual is: expected)
		}
	port: 1234
	closed(block)
		{
		SocketClient('127.0.0.1', TestSocketServer.Port)
			{|sc|
			sc.Writeline('close')
			block(sc)
			}
		SocketClient('127.0.0.1', TestSocketServer.Port)
			{|sc|
			sc.Writeline('close')
			Thread.Sleep(10) /*= milliseconds */
			block(sc)
			}
		SocketClient('127.0.0.1', TestSocketServer.Port)
			{|sc|
			sc.Writeline('before 10')
			sc.Writeline('close')
			block(sc)
			}
		SocketClient('127.0.0.1', TestSocketServer.Port)
			{|sc|
			sc.Writeline('during 10')
			sc.Writeline('close')
			block(sc)
			}
		SocketClient('127.0.0.1', TestSocketServer.Port)
			{|sc|
			sc.Writeline('after 10')
			sc.Writeline('close')
			block(sc)
			}
		SocketClient('127.0.0.1', TestSocketServer.Port)
			{|sc|
			sc.Writeline('before 10')
			sc.Writeline('during 10')
			sc.Writeline('after 10')
			sc.Writeline('close')
			block(sc)
			}
		}
	readchunked(sc, request, expected, chunksize)
		{
		sc.Writeline(request)
		actual = ''
		while false isnt chunk = sc.Read(chunksize)
			actual $= chunk
		Assert(actual.Size() is: expected.Size())
		Assert(actual is: expected)
		}
	}
