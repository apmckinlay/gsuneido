// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		tbl = .MakeTable('(name, text, path, lib_committed, type)
			key(name, lib_committed)')
		Assert(SvcSyncServer(tbl) is: #())

		r = [name: 'One', text: '1', path: 'test/path',	lib_committed: #20210901]
		QueryOutput(tbl, r)
		Assert(SvcSyncServer(tbl) isSize: 0)

		SvcSyncServer.ResetCache()
		x = SvcSyncServer(tbl)
		Assert(x isSize: 1)
		Assert(x[0] is: [from: '', to: '~', n: 1,
			cksum: Adler32(Adler32().
				Update(r.name).
				Update(r.path).
				Update(r.text).
				Update(String(r.lib_committed)).
				Value().Hex())])

		// older versions should not affect checksums
		QueryOutput(tbl, [name: 'One', text: '1', lib_committed: #20210801])
		SvcSyncServer.ResetCache()
		Assert(SvcSyncServer(tbl) is: x)

		QueryOutput(tbl, [name: 'Two', text: '2', lib_committed: #20210901.1400])
		QueryOutput(tbl, [name: 'Three', text: '3', lib_committed: #20210901.1500])

		SvcSyncServer.ResetCache()
		x = SvcSyncServer(tbl)
		Assert(x isSize: 3)
		Assert(x[0] is: #(from: "", to: "Three", n: 1, cksum: 145359356))
		Assert(x[1] is: #(from: "Three", to: "Two", n: 1, cksum: 165085744))
		Assert(x[2] is: #(from: "Two", to: "~", n: 1, cksum: 142672327))

		SvcSyncServer.ResetCache()
		x = SvcSyncServer { NSPLIT: 2 }(tbl)
		Assert(x isSize: 2)
		Assert(x[0] is: #(from: '', to: 'Two', n: 2, cksum: 576259115))
		Assert(x[1] is: #(from: 'Two', to: '~', n: 1, cksum: 142672327))

		x = SvcSyncServer(tbl, '', 'Two')
		Assert(x isSize: 2)
		Assert(x[0] is: #(from: "", to: "Three", n: 1, cksum: 145359356))
		Assert(x[1] is: #(from: "Three", to: "Two", n: 1, cksum: 165085744))
		}
	}
