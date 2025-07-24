// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
class
	{
	Increment(type, date)
		{
		Database('ensure usages (type, date, used) key (type, date)')
		RetryTransaction()
			{ |t|
			if 0 is t.QueryDo('update usages
				where type is ' $ Display(type) $ ' and date is ' $ Display(date) $
				' set used = used + 1')
				t.QueryOutput('usages', [:type, :date, used: 1])
			}
		}

	Get(type, date)
		{
		if not TableExists?('usages') or (false is rec = Query1('usages', :type, :date))
			return 0
		return rec.used
		}

	Remove(type, asof = false)
		{
		if not TableExists?('usages')
			return

		deleteQuery = 'delete usages where type is ' $ Display(type)
		asofWhere = asof is false ? '' : ' where asof <= ' $ Display(asof)
		QueryDo(deleteQuery $ asofWhere)
		}

	GetDailyTotals(type, date)
		{
		if not TableExists?('usages')
			return 0

		return QueryTotal('usages
			where type is ' $ Display(type) $
			' and date >= ' $ Display(date) $
			' and date < ' $ Display(date.Plus(days: 1)), 'used')
		}
	}