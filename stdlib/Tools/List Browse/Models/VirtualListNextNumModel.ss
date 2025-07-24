// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
// TODO: fix Timer renewing on none existent reservations
	//(restore into leave screen into table change)
// TODO: find/fix NextNum master reserve (#17000101) getting wipped out
// TODO: VirtualList should share the same NextNum with Access
class
	{
	field: false
	numClass: false
	recField: false
	table: false
	nextnumTimer: false
	New(nextNum)
		{
		if nextNum isnt false
			{
			.field = nextNum.field
			.recField = nextNum.recField
			.numClass = nextNum.numClass
			.table = nextNum.table
			}
		}

	clearNextNumTimer()
		{
		if .nextnumTimer isnt false
			{
			.nextnumTimer.Kill()
			.nextnumTimer = false
			}
		}

	ConfirmNextNum(rec)
		{
		if .listHasNextNum?(rec) and .inReserveTable?(rec[.recField])
			{
			.numClass.Confirm(rec[.recField], .table, .field)
			rec.usingNextNum? = false
			.clearNextNumTimer()
			}
		}

	listHasNextNum?(rec)
		{
		return .checkTable() and rec.GetDefault('usingNextNum?', false) is true
		}

	checkTable()
		{
		if .table is false
			return false
		return not QueryEmpty?(.table)
		}

	inReserveTable?(value, extraWhere = '')
		{
		return not QueryEmpty?(.numClass.NumQuery(.table,.field,value) $ extraWhere)
		}

	CheckPutBackNextNum(rec, field = false, newValue = false)
		{
		if .checkTable() is false
			return
		if field isnt false and field isnt .recField
			return

		if rec.GetDefault('usingNextNum?', false) is true and
			rec[.recField] isnt newValue and
			.inReserveTable?(rec[.recField])
			{
			.numClass.PutBack(rec[.recField], .table, .field)
			.clearNextNumTimer()
			rec.usingNextNum? = false
			}
		}

	ReserveNextNum(rec, query)
		{
		if .checkTable() is false
			return
		i = 0
		maxAttempts = 100
		num = false
		do
			{
			num = .numClass.Reserve(.table, .field)
			if not IsDuplicate(query, .recField, num)
				break
			.numClass.Confirm(num, .table, .field)
			} while (++i < maxAttempts)
		if i is maxAttempts
			{
			SuneidoLog('ERROR: VirtualList failed to Reserve nextNum')
			return
			}
		rec[.recField] = num
		rec.usingNextNum? = true
		.scheduleNextnumRenew(rec)
		}

	scheduleNextnumRenew(rec)
		{
		if .listHasNextNum?(rec) is false or .nextnumTimer isnt false
			return

		seconds = (.numClass.ReserveSeconds / 3).Int() /*= 100s */
		.nextnumTimer = Delay(seconds.SecondsInMs())
			{
			.renewNextNum(rec)
			}
		}

	renewNextNum(rec)
		{
		if .nextnumTimer isnt false and .checkTable() isnt false
			{
			.numClass.Renew(rec[.recField], .table,  .field)
			.nextnumTimer = false
			.scheduleNextnumRenew(rec)
			}
		}

	CheckAndClearNewNextNums(rec)
		{
		if .checkTable() is false
			return
		.clearNextNumTimer()
		if rec.GetDefault('usingNextNum?', false) is true and
			.inReserveTable?(rec[.recField])
			.numClass.PutBack(rec[.recField], .table, .field)
		}
	}