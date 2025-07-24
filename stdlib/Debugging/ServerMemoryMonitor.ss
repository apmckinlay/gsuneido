// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
function ()
	{
	ctrl = Controller
		{
		New()
			{
			.arena = .Border.Vert.arena.arena
			.dbsize = .Border.Vert.dbsize.dbsize
			.transactions = .Border.Vert.transactions.transactions
			.cursors = .Border.Vert.cursors.cursors
			.timer = SetTimer(.WndProc.Hwnd, id: 100, ms: 1000, f: .onTimer)
			}
		Controls: (Border (Vert
			(Pair (Static Arena) (Static "99,999 k" name: 'arena') name: 'arena')
			(Pair (Static 'DB Size') (Static "9999 mb" name: 'dbsize') name: 'dbsize')
			(Pair (Static Transactions) (Static "9999" name: 'transactions')
				name: 'transactions')
			(Pair (Static Cursors) (Static "9999" name: 'cursors') name: 'cursors')
			) xstretch: 0, ystretch: 0)
		onTimer(@unused)
			{
			kb = 1024
			.arena.Set((ServerEval("MemoryArena") / kb).Format("##,### k"))
			.dbsize.Set((Database.CurrentSize() / (kb * kb)).Format("#### mb"))
			.transactions.Set(Database.Transactions().Size())
			.cursors.Set(Database.Cursors())
			return 0
			}
		Destroy()
			{
			KillTimer(.WndProc.Hwnd, .timer)
			ClearCallback(.onTimer)
			super.Destroy()
			}
		}
	Window(ctrl)
	}