// Copyright (C) 2025 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one()
		{
		table = .MakeTable("(a,b) key(a)", [a: 1, b: 2])
		test =
			{|first, second|
			expected = first $ " & " $ second $ " on same record"
			t = Transaction(update:)
			x = t.Query1(table)
			y = t.Query1(table)
			switch first
				{
			case 'delete': x.Delete()
			case 'update': x.b = 123; x.Update()
				}
			switch second
				{
			case 'delete':
				Assert({ y.Delete() } throws: expected)
			case 'update':
				y.b = 123;
				Assert({ y.Update() } throws: expected)
				}
			t.Rollback()
			}
		test("delete", "delete")
		test("delete", "update")
		test("update", "delete")
		test("update", "update")
		}
	}