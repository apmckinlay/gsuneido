// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	CallClass()
		{
		Window(this, exStyle: WS_EX.TOPMOST)
		}
	New()
		{
		.arena = .FindControl('arena')
		.dbsize = .FindControl('dbsize')
		.threads = .FindControl('threads')
		.connections = .FindControl('connections')
		.readtrans = .FindControl('readtrans')
		.updatetrans = .FindControl('updatetrans')
		.cursors = .FindControl('cursors')
		.timer = SetTimer(hwnd: .WndProc.Hwnd, id: 100, ms: 200, f: .OnTimer)
		}
	Controls: (Border (Vert
		(Heading 'Server Monitor')
		Skip
		(Pair (Static Arena) (Static "99,999 k" name: 'arena'))
		(Pair (Static 'DB Size') (Static "9999 mb" name: 'dbsize'))
		(Pair (Static Threads) (Static "999" name: 'threads'))
		(Pair (Static Connections) (Static "999" name: 'connections'))
		(Pair (Static 'Read Transactions') (Static "9999" name: 'readtrans'))
		(Pair (Static 'Update Transactions') (Static "9999" name: 'updatetrans'))
		(Pair (Static Cursors) (Static "9999" name: 'cursors'))
		) xstretch: 0, ystretch: 0)
	OnTimer(@unused)
		{
		.arena.Set(ReadableSize(ServerEval("MemoryArena")))
		.dbsize.Set(ReadableSize(Database.CurrentSize()))
		.threads.Set(Thread.List().Join(', '))
		.connections.Set(Sys.Connections().Join(', '))
		.transactions()
		.cursors.Set(Database.Cursors())
		return 0
		}
	transactions()
		{
		trans = Database.Transactions()
		.readtrans.Set(.format(trans, #Even?))
		.updatetrans.Set(.format(trans, #Odd?))
		}
	format(trans, filter)
		{
		t = trans.Filter(filter)
		return '(' $ t.Size() $ ') ' $ t.Join(', ')
		}
	Destroy()
		{
		KillTimer(.WndProc.Hwnd, .timer)
		ClearCallback(.OnTimer)
		super.Destroy()
		}
	}
