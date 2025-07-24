// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(tran, query)
		{
		.t = tran
		.q = .t.Query(query)
		.query = query
		}
	qbefore: false
	qafter: false
	Next()
		{
		if .t is false or .t.Ended?()
			return Record()
		if .q is false
			.q = .qafter
		x = .q.Next()
		if x is false and .q is .qbefore
			{
			.q = .qafter
			x = .q.Next()
			}
		return x
		}
	Prev()
		{
		if .t is false or .t.Ended?()
			return Record()
		if .q is false
			.q = .qbefore
		x = .q.Prev()
		if x is false and .q is .qafter
			{
			.q = .qbefore
			x = .q.Prev()
			}
		return x
		}
	Rewind()
		{
		if .t is false or .t.Ended?()
			return
		if .qafter is false
			.q.Rewind()
		else
			{
			.close_query(.qafter)
			.close_query(.qbefore)
			.qafter = .qbefore = false
			.q = .t.Query(.query)
			}
		}
	Seek(field, prefix)
		{
		if .t is false or .t.Ended?()
			return 0
		.Close()
		.qafter = .t.Query(QueryAddWhere(.query, ' where ' $ field $ ' >= ' $
			Display(prefix)))
		.qbefore = .t.Query(qb = QueryAddWhere(.query, ' where ' $ field $ ' < ' $
			Display(prefix)))
		.q = false
		return .t.QueryCount(qb)
		}
	Columns()
		{
		return false is .q ? .qafter.Columns() : .q.Columns()
		}
	Keys()
		{
		return false is .q ?  .qafter.Keys() : .q.Keys()
		}
	Strategy()
		{
		return .q is false ? .qafter.Strategy() : .q.Strategy()
		}
	Close()
		{
		if .qafter is false
			.close_query(.q)
		else
			{
			.close_query(.qbefore)
			.close_query(.qafter)
			}
		.q = .qafter = .qbefore = false
		}
	close_query(query)
		{
		if query isnt false
			query.Close()
		}
	}
