// Copyright (C) 2012 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Setup()
		{
		super.Setup(fieldNumber?:)
		.cursors = Object()
		}

	prevActionField: 'VirtualListModelCursor_prevAction'
	makeCursor(@args)
		{
		.cursors.Add(c = VirtualListModelCursor(@args))
		return c
		}
	Test_Read()
		{
		.data = Object()
		c = .makeCursor(.data, .VirtualList_Table)
		Assert(c.Pos is: 0)
		Assert(c[.prevActionField] is: 'Next')

		// read back 10 lines
		c.ReadUp(10)
		Assert(.data isSize: 0)
		Assert(c.Pos is: 0)

		.readAndCheck(c, 10, 10, 10)

		.readAndCheck(c, 20, 30, 30)

		.readAndCheck(c, 5, 30, 30)

		// release 20 lines
		c2 = .makeCursor(.data, .VirtualList_Table)
		c2.ReleaseDown(20)
		Assert(.data isSize: 10)
		for i in ..20
			Assert(.data hasntMember: i)
		for i in 20..30
			Assert(.data[i].num is: i)
		Assert(c2.Pos is: 20)

		// read back 10 lines
		c2.ReadUp(10)
		Assert(.data isSize: 20)
		for i in ..10
			Assert(.data hasntMember: i)
		for i in 10..30
			Assert(.data[i].num is: i)
		Assert(c2.Pos is: 10)

		.readAndCheck(c2, 10, 30, 0, up:)

		.readAndCheck(c2, 5, 30, 0, up:)
		}

	readAndCheck(c, number, endSize, endPos, up = false)
		{
		if not up
			c.ReadDown(number)
		else
			c.ReadUp(number)
		Assert(.data isSize: endSize)
		for i in ..number
			Assert(.data[i].num is: i)
		Assert(c.Pos is: endPos)
		}

	Test_Release()
		{
		.data = Object()
		c = .makeCursor(.data, .VirtualList_Table)
		Assert(c.Pos is: 0)
		Assert(c[.prevActionField] is: 'Next')

		// read back 30 lines
		c.ReadDown(30)
		Assert(.data isSize: 30)
		Assert(c.Pos is: 30)
		for i in ..30
			Assert(.data[i].num is: i)

		// release 10 lines
		c.ReleaseUp(10)
		Assert(c[.prevActionField] is: 'Prev')
		Assert(.data isSize: 20)
		Assert(c.Pos is: 20)
		for i in ..20
			Assert(.data[i].num is: i)

		c.ReadDown(10)
		Assert(c[.prevActionField] is: 'Next')
		Assert(.data isSize: 30)
		Assert(c.Pos is: 30)
		for i in ..30
			Assert(.data[i].num is: i)
		}

	Test_Read_reverse()
		{
		.data = Object()
		c = .makeCursor(.data, .VirtualList_Table, startLast:)
		Assert(c.Pos is: 0)
		Assert(c[.prevActionField] is: 'Prev')

		// read 10 lines
		c.ReadUp(10)
		Assert(.data isSize: 10)
		.assertRows(-1, -9)
		Assert(c.Pos is: -10)
		Assert(c[.prevActionField] is: 'Prev')

		// read 20 more lines
		c.ReadUp(20)
		Assert(.data isSize: 30)
		.assertRows(-1, -30)
		Assert(c.Pos is: -30)
		Assert(c[.prevActionField] is: 'Prev')

		// read another 5 more lines
		c.ReadUp(5)
		Assert(.data isSize: 30)
		.assertRows(-1, -30)
		Assert(c.Pos is: -30)
		Assert(c[.prevActionField] is: 'Prev')

		// release 20 lines
		c2 = .makeCursor(.data, .VirtualList_Table, startLast:)
		c2.ReleaseUp(20)
		Assert(.data isSize: 10) // because Prev will return -10
		for(i = -1; i >= -20; i--)
			Assert(.data hasntMember: i)
		.assertRows(-21, -30)
		Assert(c2.Pos is: -20)
		Assert(c2[.prevActionField] is: 'Prev')

		// read back 10 lines
		c2.ReadDown(10)
		Assert(.data isSize: 20)
		.assertRows(-11, -30)
		for(i = -1; i >= -10; i--)
			Assert(.data hasntMember: i)
		Assert(c2.Pos is: -10)
		Assert(c2[.prevActionField] is: 'Next')

		// read back 10 lines
		c2.ReadDown(10)
		Assert(.data isSize: 30)
		.assertRows(-1, -30)
		Assert(c2.Pos is: 0)
		Assert(c2[.prevActionField] is: 'Next')

		// read back another 5 lines
		c2.ReadDown(5)
		Assert(.data isSize: 30)
		.assertRows(-1, -30)
		Assert(c2.Pos is: 0)
		Assert(c2[.prevActionField] is: 'Next')
		}

	assertRows(from, to)
		{
		for(i = from; i > to; i--)
			Assert(.data[i].num is: 30 + i)
		}

	Test_Release_non_indexed_query()
		{
		.data = Object()
		c = .makeCursor(.data, .VirtualList_Table $ ' sort field1')
		Assert(c.Pos is: 0)
		Assert(c[.prevActionField] is: 'Next')

		// read 30 lines
		c.ReadDown(30)
		Assert(.data isSize: 30)
		Assert(c.Pos is: 30)
		for i in ..30
			Assert(.data[i].field1 is: i * i)
		Assert(c[.prevActionField] is: 'Next')

		// release 10 lines
		c.ReleaseUp(10)
		Assert(.data isSize: 20)
		Assert(c.Pos is: 20)
		for i in ..20
			Assert(.data[i].field1 is: i * i)
		Assert(c[.prevActionField] is: 'Prev')
		}

	Test_ReadMultiRow()
		{
		.data = Object()
		c = .makeCursor(.data, .VirtualList_Table, setupRecord: .setupRecord)
		Assert(c.Pos is: 0)
		Assert(c[.prevActionField] is: 'Next')

		// read back 10 lines
		c.ReadUp(10)
		Assert(.data isSize: 0)
		Assert(c.Pos is: 0)

		// Read 10 lines, this will cover: <num>.<idx>
		// [0.0, 0.1, 0.2], [1.0, 1.1, 1.2], [2.0, 2.1, 2.2], [3.0, 3.1, 3.2]
		.readAndCheckMultiRows(c, 10, 0, 12, 12)
		.assertMinMax(0, 3)

		// Read an additional 20 lines, this will include the above combinations and:
		// [4.0, 4.1, 4.2], [5.0, 5.1, 5.2], [6.0, 6.1, 6.2], [7.0, 7.1, 7.2]
		// [8.0, 8.1, 8.2], [9.0, 9.1, 9.2], [10.0, 10.1, 10.2]
		.readAndCheckMultiRows(c, 20, 0, 33, 33)
		.assertMinMax(0, 10)

		// Read the remaining data
		// [11.0, 11.1, 11.2] ... [29.0, 29.1, 29.2]
		.readAndCheckMultiRows(c, 100, 0, 90, 90)
		.assertMinMax(0, 29)

		// release 20 lines
		c2 = .makeCursor(.data, .VirtualList_Table, setupRecord: .setupRecord)
		c2.ReleaseDown(20)
		.readAndCheckMultiRows(c2, 0, 7, 69, 21)
		.assertMinMax(7, 29)

		// read back 10 lines
		.readAndCheckMultiRows(c2, 10, 3, 81, 9, up:)
		.assertMinMax(3, 29)

		// read back 1 lines
		.readAndCheckMultiRows(c2, 1, 2, 84, 6, up:)
		.assertMinMax(2, 29)

		// read back remaining lines
		.readAndCheckMultiRows(c2, 5, 0, 90, 0, up:)
		.assertMinMax(0, 29)
		}

	setupRecord(rec)
		{
		// Outputs three rows for each table record
		data = rec.Copy()
		rec.Delete(all:)
		vl_multi_row_count = 0
		for i in ..3
			{
			row = data.Copy()
			row.idx = i
			rec.Add(row)
			vl_multi_row_count++
			}
		rec.Each({ it.vl_multi_row_count = vl_multi_row_count })
		rec.vl_multi_row_count = vl_multi_row_count
		}

	readAndCheckMultiRows(c, read, num, endSize, endPos, up = false)
		{
		if not up
			c.ReadDown(read)
		else
			c.ReadUp(read)
		idx = 0
		Assert(.data isSize: endSize)
		.data.Members().Sort!().Each()
			{
			Assert(.data[it].num is: num)
			Assert(.data[it].idx is: idx++)
			if idx is 3
				{
				idx = 0
				num++
				}
			}
		Assert(c.Pos is: endPos)
		}

	assertMinMax(min, max)
		{
		members = .data.Members().Sort!()
		minRec = .data[members.First()]
		Assert(minRec.num is: min)
		Assert(minRec.idx is: 0)

		maxRec = .data[members.Last()]
		Assert(maxRec.num is: max)
		Assert(maxRec.idx is: 2)
		}

	Test_ReleaseMultiRow()
		{
		.data = Object()
		c = .makeCursor(.data, .VirtualList_Table, setupRecord: .setupRecord)
		Assert(c.Pos is: 0)
		Assert(c[.prevActionField] is: 'Next')

		// Read all data, and assert its correct
		.readAndCheckMultiRows(c, 90, 0, 90, 90)
		.assertMinMax(0, 29)

		// release 10 lines
		c.ReleaseUp(11)
		Assert(c[.prevActionField] is: 'Prev')
		.readAndCheckMultiRows(c, 0, 0, 78, 78)
		.assertMinMax(0, 25)

		// release 20 lines
		c2 = .makeCursor(.data, .VirtualList_Table, setupRecord: .setupRecord)
		c2.ReleaseDown(10)
		Assert(c2[.prevActionField] is: 'Next')
		.readAndCheckMultiRows(c, 0, 4, 66, 78)
		.assertMinMax(4, 25)
		}

	Test_highCostCursor?()
		{
		mock = Mock(VirtualListModelCursor)
		mock.When.highCostCursor?().CallThrough()
		mock.When.estimateCostPerRec([anyArgs:]).Return('')
		mock.VirtualListModelCursor_query = ''
		mock.VirtualListModelCursor_cursor = false
		Assert(mock.highCostCursor?() is: false)
		mock.Verify.Never().estimateCostPerRec([anyArgs:])

		mock.When.estimateCostPerRec([anyArgs:]).Return(999)
		mock.VirtualListModelCursor_query = 'stdlib sort name'
		Assert(mock.highCostCursor?() is: false)
		mock.Verify.estimateCostPerRec([anyArgs:])

		mock.When.estimateCostPerRec([anyArgs:]).Return(10000, 10001)
		Assert(mock.highCostCursor?() is: false)
		mock.Verify.Times(3).estimateCostPerRec([anyArgs:])

		mock.When.estimateCostPerRec([anyArgs:]).Return(10000, 999)
		Assert(mock.highCostCursor?())
		}

	Test_estimateCostPerRec()
		{
		ob = Object()
		fn = VirtualListModelCursor.VirtualListModelCursor_estimateCostPerRec
		Assert(ob.Eval(fn, .mockStrategy('')) is: 0)
		Assert(ob.VirtualListModelCursor_estimateAvgCost is: 0)

		exp = 'stdlib'
		Assert(ob.Eval(fn, .mockStrategy(exp)) is: 0)
		Assert(ob.VirtualListModelCursor_estimateAvgCost is: 0)

		exp = 'stdlib [nrecs~ 50 cost~ 100]'
		Assert(ob.Eval(fn, .mockStrategy(exp)) is: 2)
		Assert(ob.VirtualListModelCursor_estimateAvgCost is: 2)

		exp = 'stdlib [nrecs~ 0 cost~ 100]'
		Assert(ob.Eval(fn, .mockStrategy(exp)) is: 100)
		Assert(ob.VirtualListModelCursor_estimateAvgCost is: 100)

		exp = 'stdlib [nrecs~ 50 cost~ 0]'
		Assert(ob.Eval(fn, .mockStrategy(exp)) is: 0)
		Assert(ob.VirtualListModelCursor_estimateAvgCost is: 0)
		}

	mockStrategy(exp)
		{
		mock = Mock()
		mock.When.Strategy().Return(exp)
		return mock
		}

	Test_smallQuery?()
		{
		fn = VirtualListModelCursor.VirtualListModelCursor_smallQuery?
		Assert(fn(.mockStrategy(''), 1) is: false)

		exp = 'stdlib'
		Assert(fn(.mockStrategy(exp), 1) is: false)

		exp = 'stdlib [nrecs~ 50 cost~ 100]'
		Assert(fn(.mockStrategy(exp), 50), msg: '50')

		exp = 'stdlib [nrecs~ 51 cost~ 100]'
		Assert(fn(.mockStrategy(exp), 50) is: false)

		exp = 'stdlib [nrecs~ 0 cost~ 100]'
		Assert(fn(.mockStrategy(exp), 1), msg: '1')
		}

	Teardown()
		{
		for c in .cursors
			c.Close()
		super.Teardown()
		}
	}
