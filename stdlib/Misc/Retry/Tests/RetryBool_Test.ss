// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Main()
		{
		mock = Mock(RetryBool)
		mock.Eval(RetryBool, 2, 1, { true })
		mock.Verify.Never().retrySleep([anyArgs:])

		mock = Mock(RetryBool)
		mock.Eval(RetryBool, 2, 1, {|count| count is 1 ? false : true })
		mock.Verify.retrySleep([anyArgs:])

		// should always run sleep one fewer time than it runs the block
		mock = Mock(RetryBool)
		mock.Eval(RetryBool, 3, 1, { false })
		mock.Verify.Times(2).retrySleep([anyArgs:])

		mock = Mock(RetryBool)
		mock.Eval(RetryBool, 3, 1, {|count| #(false, false, true)[count-1] })
		mock.Verify.Times(2).retrySleep([anyArgs:])

		mock = Mock(RetryBool)
		Assert(mock.Eval(RetryBool, 1, 1, { 'failed' }) is: 'failed (too many retries)')
		}
	}
