// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Valid?()
		{
		mock = Mock(FirstLastNameControl)
		mock.When.Valid?().CallThrough()
		mock.FirstLastNameControl_first = FakeObject(Valid?: true)
		mock.FirstLastNameControl_last = FakeObject(Valid?: true)

		// test <= 512 characters in total
		mock.When.Get().Return('')
		Assert(mock.Valid?())
		mock.When.Get().Return('FirstName, LastName')
		Assert(mock.Valid?())
		first = 'a'.Repeat(512)
		mock.When.Get().Return(first)
		Assert(mock.Valid?())
		last = 'b'.Repeat(510)
		mock.When.Get().Return(last $ ', ')
		Assert(mock.Valid?())
		last = 'b'.Repeat(509)
		mock.When.Get().Return(last $ ', ' $ 'f')	// 511 chars
		Assert(mock.Valid?())

		// test > 512 characters in total
		first = 'a'.Repeat(513)
		mock.When.Get().Return(first)
		Assert(mock.Valid?() is: false)
		last = 'b'.Repeat(512)
		mock.When.Get().Return(last $ ', ')		// 514 chars
		Assert(mock.Valid?() is: false)
		mock.When.Get().Return(last $ ', ' $ first)
		Assert(mock.Valid?() is: false)

		// individual field is invalid
		mock.FirstLastNameControl_first = FakeObject(Valid?: false)
		mock.When.Get().Return('Test, Name')
		Assert(mock.Valid?() is: false)
		mock.FirstLastNameControl_last = FakeObject(Valid?: false)
		mock.When.Get().Return('Test, Name')
		Assert(mock.Valid?() is: false)
		}
	}
