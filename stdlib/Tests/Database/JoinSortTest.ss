// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_one_to_one()
		{
		table1 = .MakeTable('(a, b, c) key(a)',
			#(a: 1, b: 11, c: 111),
			#(a: 2, b: 22, c: 222),
			#(a: 3, b: 33, c: 333),
			#(a: 4, b: 44, c: 444),
			#(a: 5, b: 55, c: 555))
		table2 = .MakeTable('(a, d, e) key(a)',
			#(a: 1, d: 11, e: 111),
			#(a: 3, d: 33, e: 333),
			#(a: 5, d: 55, e: 555))

		// next
		WithQuery(table1 $ " join by(a) " $ table2 $ " sort b")
			{ |q|
			Assert(q.Next() is: #(a: 1, b: 11, c: 111, d: 11, e: 111))
			Assert(q.Next() is: #(a: 3, b: 33, c: 333, d: 33, e: 333))
			Assert(q.Next() is: #(a: 5, b: 55, c: 555, d: 55, e: 555))
			Assert(q.Next() is: false)
			}

		// prev
		WithQuery(table1 $ " join by(a) " $ table2 $ " sort d")
			{ |q|
			Assert(q.Prev() is: #(a: 5, b: 55, c: 555, d: 55, e: 555))
			Assert(q.Prev() is: #(a: 3, b: 33, c: 333, d: 33, e: 333))
			Assert(q.Prev() is: #(a: 1, b: 11, c: 111, d: 11, e: 111))
			Assert(q.Prev() is: false)
			}

		// next
		WithQuery(table2 $ " join by(a) " $ table1 $ " sort b")
			{ |q|
			Assert(q.Next() is: #(a: 1, b: 11, c: 111, d: 11, e: 111))
			Assert(q.Next() is: #(a: 3, b: 33, c: 333, d: 33, e: 333))
			Assert(q.Next() is: #(a: 5, b: 55, c: 555, d: 55, e: 555))
			Assert(q.Next() is: false)
			}

		// prev
		WithQuery(table2 $ " join by(a) " $ table1 $ " sort d")
			{ |q|
			Assert(q.Prev() is: #(a: 5, b: 55, c: 555, d: 55, e: 555))
			Assert(q.Prev() is: #(a: 3, b: 33, c: 333, d: 33, e: 333))
			Assert(q.Prev() is: #(a: 1, b: 11, c: 111, d: 11, e: 111))
			Assert(q.Prev() is: false)
			}
		}
	}