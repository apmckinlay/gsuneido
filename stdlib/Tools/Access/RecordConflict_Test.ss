// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(RecordConflict?(
			[a: 'a', a_TS: Timestamp(), bizuser_user_modified: 'a'],
			[a: 'b', a_TS: Timestamp(), bizuser_user_modified: 'b'],
			#(a, a_TS), 0, quiet?:))

		t = Timestamp()
		Assert(RecordConflict?(
			[a: 'a', a_TS: t], [a: 'b', a_TS: t] #(a, a_TS), 0, quiet?:)
			is: false)

		t = Timestamp()
		Assert(RecordConflict?(
			[a: 'a'], [a: 'a'] #(a), 0, quiet?:)
			is: false)

		Assert(RecordConflict?(
			[a: 'a'], [a: 'b'] #(a), 0, quiet?:))
		}

	Test_changesMsg()
		{
		func = RecordConflict?['RecordConflict?_getChanges']
		Assert(func([], [], []) is: '')

		cur = [name: 'bob']
		prev = [name: 'bob']
		Assert(func(cur, prev, []) is: '')
		Assert(func(cur, prev, #(name)) is: '')
		prev = [name: 'george']
		Assert(func(cur, prev, #(name)) is: 'Name changed from "george" to "bob"\n')

		cur.age = 25
		prev.age = 24
		prev.name = 'bob'
		Assert(func(cur, prev, #(name)) is: '')
		Assert(func(cur, prev, #(name, age)) is: 'age changed from 24 to 25\n')
		prev.name = 'george'
		Assert(func(cur, prev, #(name, age))
			is: 'Name changed from "george" to "bob"\nage changed from 24 to 25\n')
		}
	}