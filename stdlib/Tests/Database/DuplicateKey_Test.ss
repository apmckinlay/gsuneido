// Copyright (C) 2013 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_output_normal()
		{
		table = .MakeTable("(a,b) key(a)", #(a: 1, b: 2))
		Assert({
			QueryOutput(table, #(a: 1, b: 3))
			} throws: "duplicate key")
		}
	Test_update_normal()
		{
		table = .MakeTable("(a,b) key(a)", #(a: 1, b: 2), #(a: 3, b: 4))
		Assert({
			Transaction(update:)
				{|t|
				x = t.Query1(table, a: 3)
				x.a = 1
				x.Update()
				}
			} throws: "duplicate key")
		}
	Test_output_conflict()
		{
		table = .MakeTable("(a,b) key(a)")
		Assert({
			Transaction(update:)
				{|t|
				t.QueryOutput(table, #(a: 1, b: 2))
				Transaction(update:)
					{|t2| t2.QueryOutput(table, #(a: 1, b: 3)) }
				}
			} throws: "conflict")
		}
	Test_update_conflict()
		{
		table = .MakeTable("(a,b) key(a)", #(a: 1, b: 2), #(a: 3, b: 4))
		Assert({
			Transaction(update:)
				{|t|
				x = t.Query1(table, a: 1)
				x.a = 9
				x.Update()
				Transaction(update:)
					{|t2|
					x = t2.Query1(table, a: 3)
					x.a = 9
					x.Update()
					}
				}
			} throws: "conflict")
		}
	}
