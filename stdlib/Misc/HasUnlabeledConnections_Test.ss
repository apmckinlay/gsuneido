// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_hasInvalidConnection?()
		{
		Assert(HasUnlabeledConnections?(#('(abc)')) is: false)
		Assert(HasUnlabeledConnections?(#('(abc)', '127.0.0.1')))
		Assert(HasUnlabeledConnections?(#('(abc)', 'tt@192.168.1.153')) is: false)
		Assert(HasUnlabeledConnections?(#('(abc)', '192.168.1.11')))
		Assert(HasUnlabeledConnections?(#('(abc)', '127.0.0.1:SocketServer-thread-371'))
			is: false)
		Assert(HasUnlabeledConnections?(#('(abc)', '127.0.0.1:suneido-thread-0'))
			is: false)
		Assert(HasUnlabeledConnections?(#('(abc)', 'tt@192.168.1.153', 'sat')))
		Assert(HasUnlabeledConnections?(#('(abc)',
			"127.0.0.1:SocketServer-0-connection-7 HTTP Server"))
			is: false)
		Assert(HasUnlabeledConnections?(#('(abc)',
			"127.0.0.1:SocketServer-0-connection-445 Rack Server"))
			is: false)
		Assert(HasUnlabeledConnections?(#('(abc)',
			"127.0.0.1:SocketServer-0-connection-445 Rack Server",
			"127.0.0.1:SocketServer-0-connection-445 HTTP Server"))
			is: false)
		Assert(HasUnlabeledConnections?(#('(http-server-monitor)',
			'127.0.0.1:SocketServer-0-connection-291 Rack Server', '(rackserver)',
			'(rackserver-server-monitor)', '(http)', '(Sys.Connections)'))
			is: false)
		Assert(HasUnlabeledConnections?(#('(abc)',
			"127.0.0.1:SocketServer-0-connection-7 HTTP Server",
			"127.0.0.1:SocketServer-0-connection-335 HTTP Server"
			"127.0.0.1:SocketServer-100-connection-335 HTTP Server"))
			is: false)
		Assert(HasUnlabeledConnections?(#('(abc)', '127.0.0.1',
			'127.0.0.1:SocketServer-thread-371',
			'127.0.0.1:SocketServer-0-connection-291 Rack Server',
			'tt@192.168.1.153')))

		Assert(HasUnlabeledConnections?(#('user@wts5', '134.234.2.123:test-thread'))
			is: false)
		Assert(HasUnlabeledConnections?(#('user@wts5', '134.234.2.123:Thread-9')))

		connections = #(
			"192.168.1.101:ThreadTotal_Pr_Transactions.CalcTotalOnServer_#20210722.1030",
			"joe@192.168.1.101")
		Assert(HasUnlabeledConnections?(connections) is: false)
		}
	}