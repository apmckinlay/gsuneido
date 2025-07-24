// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getItem()
		{
		cl = ChooseManyFieldControl
			{
			ChooseManyFieldControl_list: ('p[abc]', 'pDEF', 'aaa', 'aab', 'aac')
			ChooseManyFieldControl_allowOther: false
			ChooseManyFieldControl_setValid(@unused) {}
			}
		m = cl.ChooseManyFieldControl_getItem

		// invalid no allowOther
		Assert(m(#(), 'zzz') is: 'zzz')

		// valid exact match
		Assert(m(#(aaa bbb zzz), 'zzz') is: 'zzz')

		// test prefix match (include regex char to ensure no exception)
		Assert(m(#(), 'p[') is: 'p[abc]')

		// test mixed case match
		Assert(m(#(), 'pd') is: 'pDEF')

		// prefix matches multiple, orig entry should be returned
		Assert(m(#(), 'aa') is: 'aa')
		}

	Test_valid()
		{
		cl = ChooseManyFieldControl
			{
			ChooseManyFieldControl_list: ('p[abc]', 'pDEF')
			ChooseManyFieldControl_allowOther: false
			ChooseManyFieldControl_setValid(valid) { throw String(valid) }
			}
		// do not allow other
		Assert({ cl.ChooseManyFieldControl_getItem(#(), 'aaa') } throws: 'false' )

		cl = ChooseManyFieldControl
			{
			ChooseManyFieldControl_list: ('p[abc]', 'p[DEF]')
			ChooseManyFieldControl_allowOther: true
			ChooseManyFieldControl_setValid(valid) { throw String(valid) }
			}
		m = cl.ChooseManyFieldControl_getItem
		// allow other
		Assert({ m(#(), 'aaa') } throws: 'true' )
		// multiple matches
		Assert({ m(#(), 'p[') } throws: 'false' )
		}
	}