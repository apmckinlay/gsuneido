// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Test_init()
		{
		s = StickyFields()
		Assert(s.StickyFields_stickyFields is: false)

		s = StickyFields(#('a', 'b'))
		Assert(s.StickyFields_stickyFields is: [a: '', b: ''])
		}

	Test_UpdateStickyField()
		{
		s = StickyFields(#('a', 'b'))

		s.UpdateStickyField([a: 'hello'], 'a')
		Assert(s.StickyFields_stickyFields is: [a: 'hello', b: ''])

		s.UpdateStickyField([c: 'hello c'], 'c')
		Assert(s.StickyFields_stickyFields is: [a: 'hello', b: ''])

		rec = QueryFirst(.VirtualList_Table $ ' sort num')
		s.UpdateStickyField(rec, 'a')
		Assert(s.StickyFields_stickyFields is: [a: 'hello', b: ''])
		}

	Test_setRecordStickyFields()
		{
		s = StickyFields(#('a', 'b'))
		s.UpdateStickyField([a: 'hello'], 'a')

		rec = Record()
		s.SetRecordStickyFields(rec)
		Assert(rec is: [a: 'hello'])
		}
	}