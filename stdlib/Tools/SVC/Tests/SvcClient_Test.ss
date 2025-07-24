// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	socketClient: class
		{
		Default(@unused){ } // Catches the various writes / other unespecified calls

		readReturn: ''
		SetReadReturn(.readReturn) { }

		Read(n)
			{
			res = .readReturn[.. n]
			.readReturn = .readReturn[n ..]
			return res
			}

		Readline()
			{
			res = .readReturn.BeforeFirst('\r\n')
			.readReturn = .readReturn.AfterFirst('\r\n')
			return res
			}
		}

	svccl: class
		{
		readReturn: #()
		SetReadReturn(.readReturn) { }

		Read(@unused)
			{
			return .readReturn
			}
		}

	svcClient(name, path, id, lib_committed, comment)
		{
		client = new SvcClient
		svccl = new .svccl
		svccl.SetReadReturn([:name, :path, :id, :lib_committed, :comment])
		client.SvcClient_svccl = svccl
		return client
		}

	Test_readRec()
		{
		date = #20210101.1010
		text = 'Here is the method\r\nIt would have a lot of text\r\nand newlines'
		sc = new .socketClient
		sc.SetReadReturn(text.Size() $ '\r\n' $	text)
		comment = 'this is a test\x03\x04with newlines'
		svcClient = .svcClient(#test, `test\folder`, #default, date, comment)

		// Standard read
		x = svcClient.SvcClient_readRec(sc, 'test')
		Assert(x.name is: 'test')
		Assert(x.lib_committed is: date)
		Assert(x.path is: `test\folder`)
		Assert(x.comment is: 'this is a test\r\nwith newlines')
		Assert(x.text is: text)

		// An error message is returned, SvcSocketClient will handle any logging
		svcClient.SvcClient_svccl.SetReadReturn('ERR Message')
		Assert(svcClient.SvcClient_readRec(sc, 'test') is: false)

		// Returned record name doesn't match requested name
		svcClient = .svcClient(#diffName, `test\folder`, #default, date, comment)
		Assert(svcClient.SvcClient_readRec(sc, 'test') is: false)
		}
	}
