// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(t1, t2)
		{
		d = new LibCompare(t1, t2)
		return d.Result
		}
	New(t1, t2)
		{
		.Result = Object()
		// call function with tables as args
		.libcompare(t1, t2)
		}
	book: false
	libcompare(t1, t2)
		{
		// check if comparing books
		.book = BookTable?(t1) and BookTable?(t2)
		Transaction(read:)
			{|t|
			q_string = (.book) ? " sort name" : " where group = -1 sort name"
			q1 = t.Query(t1 $ q_string)
			q2 = t.Query(t2 $ q_string)

			.x = q1.Next()
			.y = q2.Next()
			.addbookinfo()
			while (.x isnt false or .y isnt false)
				{
				if (.y isnt false and (.x is false or .y.name < .x.name))
					// record added
					{
					.Result.Add(Object("+", .y.num, .y.name))
					.y = q2.Next()
					}
				else if (.y is false or .x.name < .y.name)
					// record removed
					{
					.Result.Add(Object("-", .x.num, .x.name))
					.x = q1.Next()
					}
				else if (.x.name is .y.name)
					{
					if (.x.text.Trim() isnt .y.text.Trim())
						{
						.Result.Add(Object("#", .x.num, .x.name))
						}
					.x = q1.Next()
					.y = q2.Next()
					}
				.addbookinfo()
				}
			q1.Close()
			q2.Close()
			}
		}
	addbookinfo()
		{
		if (not .book)
			return
		pathSize = 4
		if .x isnt false and .x.text[.. pathSize] isnt "path"
			.x.text = "path:" $ .x.path $ ", order:" $ .x.order $ "\r\n" $ .x.text
		if .y isnt false and .y.text[.. pathSize] isnt "path"
			.y.text = "path:" $ .y.path $ ", order:" $ .y.order $ "\r\n" $ .y.text
		}
	}
