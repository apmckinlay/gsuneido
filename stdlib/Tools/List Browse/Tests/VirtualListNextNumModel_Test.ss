// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_NextNum()
		{
		date = Date()
		.SpyOn(GetNextNum.GetNextNum_date).Return(date)
		numTable = .MakeTable('(nextNum, getnextnum_reserved_till) key(nextNum)')
		dataTable = .MakeTable('(numField, data) key (numField)')
		QueryOutput(numTable, [nextNum: 100, getnextnum_reserved_till: Date.Begin()])
		virtualListNextClass = VirtualListNextNumModel
			{
			VirtualListNextNumModel_scheduleNextnumRenew(unused) { }
			}
		nextNum = virtualListNextClass(nextNum: Object(numClass: GetNextNum,
			field: 'nextNum', table: numTable, recField: 'numField'))

		rec = []
		nextNum.ReserveNextNum(rec, dataTable)

		Assert(rec.usingNextNum?)
		Assert(rec.numField is: 100)
		Assert(Query1(numTable, nextNum: 100).getnextnum_reserved_till
			is: date.Plus(seconds: GetNextNum.ReserveSeconds))
		nextNum.ConfirmNextNum(rec)
		Assert(rec.usingNextNum? is: false)
		Assert(Query1(numTable, nextNum: 100) is: false)


		rec = []
		nextNum.ReserveNextNum(rec, dataTable)

		//someone changed their num manually so we should no longer reserve the old number
		nextNum.CheckPutBackNextNum(rec, field: 'numField', newValue: 105)
		Assert(Query1(numTable, nextNum: 101).getnextnum_reserved_till
			is: date.Plus(seconds: -1))

		rec.usingNextNum? = false
		rec.num = 2000
		nextNum.CheckPutBackNextNum(rec, field: 'num', newValue: 2002)
		Assert(Query1(numTable, nextNum: 2002) is: false)
		Assert(Query1(numTable, nextNum: 2000) is: false)
		}
	}