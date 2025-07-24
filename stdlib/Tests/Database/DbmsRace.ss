// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
//
// The purpose of this is to randomly and concurrently run dbmsserver.go code
// to use with the Go race detector
class
	{
	CallClass(nthreads = 2)
		{
		Assert(Sys.Client?())
		Assert(ExePath().Has?("race"))
		for .. nthreads
			Thread(.run)
		}
	run()
		{
		fns = .Members().Filter({ it.Prefix?("F_") })
		forever
			try
				this[fns.RandVal()]()
			catch(unused, "endthread")
				return
		}
	F_EndThread() // close/restart thread/session
		{
		Thread(.run)
		throw "endthread"
		}
	F_Cursor()
		{
		Cursor("stdlib")
			{|unused|
			}
		}
	F_Query()
		{
		WithQuery("stdlib")
			{|unused|
			}
		}
	F_ReadTran()
		{
		Transaction(read:)
			{|unused|
			}
		}
	F_UpdateTran()
		{
		Transaction(update:)
			{|unused|
			}
		}
	F_Rollback()
		{
		t = Transaction(update:)
		t.Rollback()
		}
	F_Transactions()
		{
		Database.Transactions()
		}
	F_Cursors()
		{
		Database.Cursors()
		}
	}