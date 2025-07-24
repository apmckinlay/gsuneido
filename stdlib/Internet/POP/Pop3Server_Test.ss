// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_findMsg()
		{
		mock = Mock(Pop3Server)
		mock.Pop3Server_messages = Object(0, 1, 2, 3)
		mock.When.findMsg([anyArgs:]).CallThrough()
		Assert(mock.findMsg('1') is: 0)
		Assert(mock.findMsg('2') is: 1)
		Assert(mock.findMsg('3') is: 2)
		Assert(mock.findMsg('4') is: 3)
		Assert(mock.findMsg('5') is: false)
		}

	Test_Transaction()
		{
		mock = Mock(Pop3Server)
		mock.When.Writeline([anyArgs:]).Do({|call/*unused*/|  })
		mock.When.Pop3Server_isMember([anyArgs:]).Do(
			{|call| Pop3Server.Member?(call[1]) })

		mock.When.Transaction([anyArgs:]).CallThrough()
		mock.Transaction('NOOP', '')
		mock.Verify.Tran_NOOP(#(anyArgs:))
		mock.Verify.Writeline(#(anyArgs:))

		mock.Pop3Server_user = 'test'
		mock.Transaction('QUIT', '')
		mock.Verify.Tran_QUIT(#(anyArgs:))
		mock.Verify.Times(2).Writeline(#(anyArgs:))
		mock.Verify.Complete()

		mock.Transaction('RSET', '')
		mock.Verify.Tran_RSET(#(anyArgs:))
		mock.Verify.Times(3).Writeline(#(anyArgs:))

		mock.Transaction('WRONG_CALL', '')
		mock.Verify.Writeline("-ERR expecting STAT, LIST, RETR, DELE, NOOP, or RSET")
		}
	}