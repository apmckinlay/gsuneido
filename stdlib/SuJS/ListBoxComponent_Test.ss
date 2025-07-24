// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_alphaNumeric()
		{
		fn = ListBoxComponent.ListBoxComponent_alphaNumeric?
		Assert(fn(Object()) is: false)
		Assert(fn('') is: false)
		Assert(fn('test') is: false)
		Assert(fn('Test') is: false)
		Assert(fn('.') is: false)
		Assert(fn('/') is: false)
		Assert(fn(`\`) is: false)
		Assert(fn('~') is: false)
		Assert(fn('?') is: false)
		Assert(fn('10') is: false)
		Assert(fn('a'))
		Assert(fn('j'))
		Assert(fn('t'))
		Assert(fn('T'))
		Assert(fn('K'))
		Assert(fn('L'))
		for i in ..10
			Assert(fn(String(i)))
		}

	Test_findItemPos()
		{
		fn = ListBoxComponent.ListBoxComponent_findItemPos
		mock = Mock(ListBoxComponent)
		mock.When.selectFirstFound([anyArgs:]).CallThrough()
		mock.When.clearSelect().Do({ })
		mock.When.select([anyArgs:]).Do({ })
		mock.ListBoxComponent_items = Object()
		mock.ListBoxComponent_selected = false
		mock.Eval(fn, 'a')
		mock.Verify.Never().clearSelect()
		mock.Verify.Never().select([anyArgs:])

		// nothing found
		mock.ListBoxComponent_items = Object([innerText: 'abc'], [innerText: 'efg'])
		mock.Eval(fn, 'z')
		mock.Verify.Never().clearSelect()
		mock.Verify.Never().select([anyArgs:])

		// search first available
		mock.ListBoxComponent_items = Object([innerText: 'abc'], [innerText: 'efg'])
		mock.Eval(fn, 'e')
		mock.Verify.Times(1).clearSelect()
		mock.Verify.select([innerText: 'efg'])

		// search next available
		mock.ListBoxComponent_items = Object([innerText: 'abc'], [innerText: 'efg'],
			[innerText: 'exy'])
		mock.ListBoxComponent_selected = [innerText: 'efg']
		mock.Eval(fn, 'e')
		mock.Verify.Times(2).clearSelect()
		mock.Verify.select([innerText: 'exy'])

		mock.ListBoxComponent_selected = [innerText: 'exy'] // search around
		mock.Eval(fn, 'e')
		mock.Verify.Times(3).clearSelect()
		mock.Verify.Times(2).select([innerText: 'efg'])

		mock.Eval(fn, 'z') // has current select, but nothing found on next
		mock.Verify.Times(3).clearSelect()
		mock.Verify.Times(2).select([innerText: 'efg'])
		}
	}