// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_rewindLines_currentLine()
		{
		data = "@ making sure this goes to @, # so the At symbol should be first read"
		pos = data.Size()
		f = FakeFile(data)
		res = RewindLines(f, pos)
		Assert(f.Read(1) is: "@")
		Assert(res is: 0)

		data = "@ making sure this stops at \n# so the hashtag should be first read"
		f = FakeFile(data)
		RewindLines(f, pos)
		Assert(f.Read(1) is: "#")

		data = "@ making sure this stops at \n# so the last part should be \nfirst read"
		f = FakeFile(data)
		RewindLines(f, pos)
		Assert(f.Readline() is: "first read")
		}

	Test_rewindLines_previousLine()
		{
		data ="1test\n2test\n3test\n4test\n6test\n7test\n8test"
		f = FakeFile(data)

		//this will put it in line 4, reading back a line will put it in 3
		RewindLines(f, 18)
		res = RewindLines(f, 18, 2)
		Assert(res is: 12)

		//Sets cursor to start of 8test
		RewindLines(f, data.Size())

		//previous line will be 7test or 30
		res = RewindLines(f, f.Tell(), 2)
		Assert(res is: 30)

		//essentially, each step will be 6 steps before the other
		res = RewindLines(f, f.Tell(), 2)
		Assert(res is: 24)

		res = RewindLines(f, f.Tell(), 2)
		Assert(res is: 18)

		res = RewindLines(f, f.Tell(), 2)
		Assert(res is: 12)

		res = RewindLines(f, f.Tell(), 2)
		Assert(res is: 6)

		// read to start of 1test
		res = RewindLines(f, f.Tell(), 2)
		Assert(res is: 0)

		// attempt to read further back, should simply return 0
		res = RewindLines(f, f.Tell(), 2)
		Assert(res is: 0)
		}

	Test_rewindLines_severalLines()
		{
		data ="1test\n2test\n3test\n4test\n6test\n7test\n8test"
		f = FakeFile(data)

		//Sets cursor to start of 8test
		RewindLines(f, data.Size())

		//previous line will be 6test or 24
		res = RewindLines(f, f.Tell(), 3)
		Assert(res is: 24)

		//essentially, each step will be 12 steps before the other
		res = RewindLines(f, f.Tell(), 3)
		Assert(res is: 12)

		// read to start of 1test
		res = RewindLines(f, f.Tell(), 3)
		Assert(res is: 0)

		// attempt to read further back, should simply return 0
		res = RewindLines(f, f.Tell(), 3)
		Assert(res is: 0)
		}

	Test_prevChar()
		{
		data = "test1"
		f = FakeFile(data)
		for (i = 0; i < data.Size(); i++)
			{
			res = RewindLines.RewindLines_prevChar(f, data.Size() - i)
			Assert(res is: data[data.Size() - 1 - i])
			}
		}
	}