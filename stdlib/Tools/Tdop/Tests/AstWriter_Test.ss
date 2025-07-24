// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	src1: '// test function
	function (a /*test*/, b = false)
		{
		// test statement
		a + 1
		return b // test return
		}'
	src2: 'function ()
		{
		Print(a)
		c = a + b
		return c
		}'
	Test_RewriteByString()
		{
		manager = AstWriteManager(.src1)
		writer = manager.GetNewWriter()
		root = writer.GetRoot()

		// test ADD
		writer.Add(root[5], 1, 'b = 0\r\n')
		writer.Add(root[5], 1, 'b = a + 1\r\n')
		writer.Add(root[2], 2, ', c = 3')
		Assert({ writer.Add(root, 1, 'test') } throws:)
		Assert(writer.ToString() like: '// test function
	function (a /*test*/, b = false, c = 3)
		{
		// test statement
		a + 1
b = 0
b = a + 1
		return b // test return
		}')

		// test REPLACE
		writer.Replace(root[5][1][2][0][0], 'b + c')
		writer.Replace(root[5][0], 'a += 1')
		writer.Replace(root[2][0][1], 'a = 1')
		Assert(writer.ToString() like: '// test function
	function (a = 1 /*test*/, b = false, c = 3)
		{
		// test statement
		a += 1
b = 0
b = a + 1
		return b + c // test return
		}')

		// test REMOVE
		writer.Remove(root[2][0])
		writer.Remove(root[5][0])
		Assert(writer.ToString() like: '// test function
	function ( b = false, c = 3)
		{
		// test statement

b = 0
b = a + 1
		return b + c // test return
		}')
		writer.Remove(root[5][1])
		Assert(writer.ToString() like: '// test function
	function ( b = false, c = 3)
		{
		// test statement

b = 0
b = a + 1
		 // test return
		}')
		}

	Test_RewriteByWriter()
		{
		manager1 = AstWriteManager(.src1)
		writer1 = manager1.GetNewWriter()
		root1 = writer1.GetRoot()

		root2 = Tdop(.src2)
		manager2 = AstWriteManager(.src2, root2)
		writer2 = manager2.GetNewWriter(root2[5])

		Assert(writer2.ToString() like: 'Print(a)
		c = a + b
		return c')
		Assert(writer2.Length() is: 33)

		writer1.Replace(root1[5], writer2)
		Assert(writer1.ToString() like: '// test function
	function (a /*test*/, b = false)
		{
		// test statement
		Print(a)
		c = a + b
		return c // test return
		}')

		writer1.Remove(root1[5])
		Assert(writer1.ToString() like: '// test function
	function (a /*test*/, b = false)
		{
		// test statement
		 // test return
		}')

		writer3 = manager1.GetNewWriter(root1[2][1])
		writer3.Replace(root1[2][1][1], ', c')
		writer1.Add(root1[2], 2, writer3)
		Assert(writer1.ToString() like: '// test function
	function (a /*test*/, b = false, c = false)
		{
		// test statement
		 // test return
		}')
		}

	Test_ReWrite()
		{
		src = `TestTest
	{
	member1: +1
	member2: "a" $ 'b'
	Method1(args)
		{
		if ! args.a and not args[#b] or args.c or (args.d is #abc) and
			(args.e isnt #20171010)
			do
				{
				GlobalName.Method(@+1 args)
				}
			while args.f << 1 is 1
		return ;
		}
	Method2: function(@args)
		{
		try
			{
			rec = #(1, is: true, isnt: (1, 2, 3), -1)
			rec.Each()
				{ |item|
				Print(item)
				}
			}
		catch (e, 'test')
			args.Add(rec)
		}
	}`
		rewrite = AstWriteManager(src).GetNewWriter().ReWrite()
		Assert(rewrite is: src)
		}
	}