// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Assert(TextBestFit(5, "Hello World", .measure) is: "Hello ")
		// if no full word will fit, returns as much as possible
		Assert(TextBestFit(4, "Hello World", .measure) is: "Hell")
		Assert(TextBestFit(4, "o World", .measure) is: "o ")
		Assert(TextBestFit(50, "Hello World", .measure) is: "Hello World")
		// if possible, doesn't split on a word, returns the last full word that will fit
		Assert(TextBestFit(14, "Hello World How Are You", .measure)
			is: "Hello World ")

		// Is line based. If a whole line will fit, it will only return that line
		Assert(TextBestFit(500, "Hello World\r\nHow Are You", .measure)
			is: "Hello World\r\n")

		report = class
			{
			PlainText?()
				{
				return true
				}
			}
		Assert(TextBestFit(500, "Hello World\r\nHow Are You", .measure, :report)
			is: "Hello World\r\nHow Are You")

		}
	measure(str)
		{
		return str.Size()
		}
	}
