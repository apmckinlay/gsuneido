// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		findIn = false
		valToFind = "Test"
		Assert(FindValueInX.Reference?(findIn, valToFind) is: false)

		findIn = false
		valToFind = true
		Assert(FindValueInX.Reference?(findIn, valToFind) is: false)

		findIn = "sample"
		valToFind = "sample"
		Assert(FindValueInX.Reference?(findIn, valToFind))

		findIn = "one,two,three"
		valToFind = "four"
		Assert(FindValueInX.Reference?(findIn, valToFind) is: false)

		findIn = "one,two,three"
		valToFind = "two"
		Assert(FindValueInX.Reference?(findIn, valToFind))

		findIn = Object("one", "two", "three")
		valToFind = "four"
		Assert(FindValueInX.Reference?(findIn, valToFind) is: false)

		findIn = Object("one", "two", "three")
		valToFind = "two"
		Assert(FindValueInX.Reference?(findIn, valToFind))

		findIn = Object(one: 1, two: 2, three: 3)
		valToFind = "four"
		Assert(FindValueInX.Reference?(findIn, valToFind) is: false)

		findIn = Object(one: 1, two: 2, three: 3)
		valToFind = "two"
		Assert(FindValueInX.Reference?(findIn, valToFind))
		}
	}