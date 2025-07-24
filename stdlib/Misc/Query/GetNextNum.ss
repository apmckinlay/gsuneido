// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
// new version of GetNextNumber
class
	{
	Create(table, field = "nextnum", nextnum = 1)
		{
		Database("create " $ table $ " (" $ field $ ", getnextnum_reserved_till)
			key(" $ field $ ")")
		Transaction(update:)
			{ |t|
			.output(t, table, field, nextnum, Date.Begin())
			}
		}
	Reserve(table, field = "nextnum")
		{
		RetryTransaction()
			{ |t|
			if ((x = t.QueryFirst(table $
				' where getnextnum_reserved_till < ' $ Display(.date()) $
					' sort ' $ field)) is false)
				throw "GetNextNum failed, no records in: " $ table
			nextnum = x[field]
			if x.getnextnum_reserved_till is Date.Begin()
				{
				do // skip reserved numbers
					++x[field]
					while false isnt t.Query1(table $
						' where ' $ field $ ' = ' $ Display(x[field]))

				x.Update()
				.output(t, table, field, nextnum, .expiry())
				}
			else // use a number that's expired or been put back
				.renew(x)
			}
		return nextnum
		}
	Renew(num, table, field = "nextnum", skipLog = false)
		{
		RetryTransaction()
			{ |t|
			if false isnt x = t.Query1(.NumQuery(table, field, num))
				{
				.renew(x)
				return true
				}
			}
		if not skipLog
			.Log(table, field, num, "Renew on non-existent reservation: ")
		return false
		}
	renew(x)
		{
		x.getnextnum_reserved_till = .expiry()
		x.Update()
		}
	Confirm(num, table, field = "nextnum")
		{
		// remove reservation so no one can take it
		RetryTransaction()
			{ |t|
			if 1 isnt t.QueryDo('delete ' $ .NumQuery(table, field, num) $
				' where getnextnum_reserved_till isnt ' $ Display(Date.Begin()))
				{
				.Log(table, field, num, "Confirm on non-existent reservation: ")
				return false
				}
			return true
			}
		}
	Log(table, field, num, msg)
		{
		SuneidoLog(
			"WARNING: GetNextNum: " $ msg $ table $ ", " $ field $ ", " $ num,
			calls:)
		}
	NumQuery(table, field, num)
		{
		return table $ ' where ' $ field $ ' = ' $ Display(num)
		}
	ReserveSeconds: 300
	expiry()
		{
		// two minute reservations (to match LockManager)
		return .date().Plus(seconds: .ReserveSeconds)
		}
	PutBack(num, table, field = "nextnum")
		{
		Assert(num isnt '')
		// output with expired time so it will be re-used
		till = .date().Plus(seconds: -1)
		RetryTransaction()
			{ |t|
			if false isnt x = t.Query1(.NumQuery(table, field, num))
				{
				x.getnextnum_reserved_till = till
				x.Update()
				}
			else
				.output(t, table, field, num, till)
			}
		}
	output(t, table, field, num, till)
		{
		x = Record()
		x[field] = num
		x.getnextnum_reserved_till = till
		t.QueryOutput(table, x)
		}

	ChangeNextNum(table, field, nextnum)
		{
		Transaction(update:)
			{ |t|
			t.QueryDo('delete ' $ table)
			.output(t, table, field, nextnum, Date.Begin())
			}
		}

	date() // so test can override
		{ return Timestamp() }

	Convert(table, field) // convert table from old GetNextNumber
		{
		if QueryColumns(table).Has?('getnextnum_reserved_till')
			return
		Database("ensure " $ table $ " (getnextnum_reserved_till)
			key(" $ field $ ")")
		Database("alter " $ table $ " drop key()")
		QueryDo("update " $ table $
			" set getnextnum_reserved_till = " $ Display(Date.Begin()))
		}
	CallClass(table, field) // for compatibility with old code
		{ // reserve + confirm
		RetryTransaction()
			{ |t|
			if ((x = t.QueryFirst(table $
				' where getnextnum_reserved_till < ' $ Display(.date()) $
					' sort ' $ field)) is false)
				throw "GetNextNum failed, no records in: " $ table
			nextnum = x[field]
			if x.getnextnum_reserved_till is Date.Begin()
				{
				do // skip reserved numbers
					++x[field]
					while false isnt t.Query1(table $
						' where ' $ field $ ' = ' $ Display(x[field]))

				x.Update()
				}
			else // use a number that's expired or been put back
				x.Delete()
			}
		return nextnum
		}

	// called from setup options screens (i.e. Eta_Configuration)
	UnreservedNextNum(table, field)
		{
		return Query1(table, getnextnum_reserved_till: Date.Begin())[field]
		}
	UpdateUnreservedNextNum(t, oldrec, newrec, field, table, table_field)
		{
		if oldrec[field] is newrec[field]
			return
		t.QueryDo('delete ' $ table $
			' where getnextnum_reserved_till < ' $ Display(.date()) $
			' and getnextnum_reserved_till > ' $ Display(Date.Begin()))
		t.QueryDo('update ' $ table $
			' where getnextnum_reserved_till is ' $ Display(Date.Begin()) $
			' set ' $ table_field $ ' = ' $ Display(newrec[field]))
		}
	}
