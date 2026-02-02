// Copyright (C) 2016 Axon Development Corporation All rights reserved worldwide.
Test
	{
	cmp(target, line)
		{
		if line is false or line.Has?("$%")
			return false
		if line is target
			return 0
		if line < target
			return -1
		if line > target
			return 1
		return false
		}

	Test_hiLoValueBeforeInvalidSect()
		{
		data = "0test\n1test\n$%$%$%$%\n$%$%$%$%\n$%$%$%$%\n2test"
		lo = 0
		hi = data.Size()
		mid = ((hi + lo)/ 2).Ceiling()
		f = (FakeFile)(data)
		hiF = FileBinarySearch.FileBinarySearch_highestValueBeforeInvalidSect
		loF = FileBinarySearch.FileBinarySearch_lowestValueAfterInvalidSect

		Assert(hiF(f, { .cmp("1test", it) }, mid) is: 6)

		Assert(loF(f, { .cmp("2test", it) }) is: 39)
		Assert(f.Readline() is: "2test")

		//Ensure program gets to start/end of file despite invalids
		data = "$%$%$%$%\n$%$%$%$%\n$%$%$%$%end"

		hi = data.Size()
		mid = ((hi + lo)/ 2).Ceiling()
		f = (FakeFile)(data)
		Assert(hiF(f, { .cmp("1test", it) }, mid) is: false)

		loF(f, { .cmp("2test", it) })
		Assert(f.Readline() is: false)

		Assert(FileBinarySearch(f, { .cmp("2test", it) }) is: false)
		}

	Test_orderMaintained?()
		{
		data = "1test\n2test\n3test\n4test\n6test\n7test\n8test"
		lo = 0
		hi = data.Size()
		mid = (( hi + lo ) / 2).Ceiling()
		f = (FakeFile)(data)
		orderCmp = function (x,y) { return x <= y }
		res = FileBinarySearch.FileBinarySearch_orderMaintained?(f, orderCmp, lo, hi, mid)
		Assert(res)

		data = "8test\n7test\n6test\n4test\n3test\n2test\n1test"
		f = (FakeFile)(data)
		orderCmp = function (x,y) { return x <= y }
		res = FileBinarySearch.FileBinarySearch_orderMaintained?(f, orderCmp, lo, hi, mid)
		Assert(res is: false)

		data = "1test\n7test\n6test\n4test\n3test\n2test\n8test"
		f = (FakeFile)(data)
		orderCmp = function (x,y) { return x <= y }
		res = FileBinarySearch.FileBinarySearch_orderMaintained?(f, orderCmp, lo, hi, mid)
		Assert(res) // 1test <= 4test <= 8test >>> only considers selected lines

		f = (FakeFile)(data)
		orderCmp = function (x,y) { return x <= y }
		res = FileBinarySearch(f, { .cmp("2test", it) }, orderCmp)
		Assert(res is: false) // after full processing, it is seen as unsorted
		}

	Test_getLine()
		{
		data = "1test\n2test\n3test\n4test\n6test\n7test\n8test"
		f = (FakeFile)(data)
		res = FileBinarySearch.FileBinarySearch_getLine(f, 0)
		Assert(res is: "1test")

		res = FileBinarySearch.FileBinarySearch_getLine(f, 4)
		Assert(res is: "1test")

		res = FileBinarySearch.FileBinarySearch_getLine(f, 48)
		Assert(res is: "8test")

		res = FileBinarySearch.FileBinarySearch_getLine(f, 20)
		Assert(res is: "4test")
		}

	Test_EmptyFile()
		{
		f = (FakeFile)("")
		FileBinarySearch(f, "shouldn't be used")
		Assert(f.Tell() is: 0)
		}

	Test_OnlyLineMatches()
		{
		data = "test1"
		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("test1", it) })
		Assert(f.Tell() is: 0)
		Assert(f.Readline() is: "test1")
		}

	Test_FileWithOneLine()
		{
		data = "test1"
		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("", it) })
		Assert(f.Tell() is: 0)

		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("test2", it) })
		Assert(f.Tell() is: data.Size())
		}

	Test_FileWithMultiLines()
		{
		data = "1test\n2test\n3test\n4test\n6test\n7test\n8test"

		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("1test", it) })
		Assert(f.Tell() is: 0)
		Assert(f.Readline() is: "1test")

		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("2test", it) })
		Assert(f.Tell() is: 6)
		Assert(f.Readline() is: "2test")

		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("3test", it) })
		Assert(f.Tell() is: 12)
		Assert(f.Readline() is: "3test")

		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("4test", it) })
		Assert(f.Tell() is: 18)
		Assert(f.Readline() is: "4test")

		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("6test", it) })
		Assert(f.Tell() is: 24)
		Assert(f.Readline() is: "6test")

		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("7test", it) })
		Assert(f.Tell() is: 30)
		Assert(f.Readline() is: "7test")

		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("8test", it) })
		Assert(f.Tell() is: 36)
		Assert(f.Readline() is: "8test")

		// 9test is > then all the other lines, file should move to very end
		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("9test", it) })
		Assert(f.Tell() is: 41)
		Assert(f.Readline() is: false)

		// 5test is > 4test, should be right after it (6test technically as it is the
		// closest greater/equal value)
		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("5test", it) })
		Assert(f.Tell() is: 24)
		Assert(f.Readline() is: "6test")

		// 10test < 1test so there for the next >= is test1
		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("10test", it) })
		Assert(f.Tell() is: 0)
		Assert(f.Readline() is: "1test")

		// 20test < 2test so there for the next >= is test2
		f = (FakeFile)(data)
		FileBinarySearch(f, { .cmp("20test", it) })
		Assert(f.Tell() is: 6)
		Assert(f.Readline() is: "2test")
		}
	}