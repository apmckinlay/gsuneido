// Copyright (C) 2002 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(query)
		{
		.q = Cursor(query)
		.query = query
		}
	qbefore: false
	qafter: false
	field: false
	Next(tran)
		{
		if .q is false and false is .q = .qafter
			return .logError('Next')
		x = .q.Next(tran)
		if x is false and .q is .qbefore
			{
			.q = .qafter
			x = .q.Next(tran)
			}
		return x
		}
	Prev(tran)
		{
		if .q is false and false is .q = .qbefore
			return .logError('Prev')
		x = .q.Prev(tran)
		if x is false and .q is .qafter
			{
			.q = .qbefore
			x = .q.Prev(tran)
			}
		return x
		}
	logError(method)
		{
		position = method is "Next" ? "qafter" : "qbefore"
		SuneidoLog("ERROR: SeekCursor." $ method $ " expected " $ position $
			" to contain next Cursor but was false", calls:)
		return false
		}
	Rewind()
		{
		if .qafter is false
			.q.Rewind()
		else
			{
			// this loses seek field order - shouldn't rewind keep it?
			.qafter.Close()
			.qbefore.Close()
			.qafter = .qbefore = false
			.q = Cursor(.query)
			.field = false
			}
		}
	Seek(field, prefix)
		{
		query = QueryStripSort(.query) $ ' sort ' $ field
		qafter = Cursor(QueryAddWhere(query,
			' where ' $ field $ ' >= ' $ Display(prefix)))
		qbefore = Cursor(QueryAddWhere(query,
			' where ' $ field $ ' < ' $ Display(prefix)))
		// only set these if above succeeds
		.Close()
		.query = query
		.qafter = qafter
		.qbefore = qbefore
		.q = false
		.field = field
		}
	Order()
		{
		return .field // apm - but .field is not necessarily the order ???
		// it's even possible that qbefore and qafter could have different orders
		}
	Columns()
		{
		return false is .q ? .qafter.Columns() : .q.Columns()
		}
	Keys()
		{
		return false is .q ? .qafter.Keys() : .q.Keys()
		}
	Output(t, x)
		{
		(false is .q ? .qafter : .q).Output(t, x)
		}
	Strategy()
		{
		return .q is false ? .qafter.Strategy() : .q.Strategy()
		}
	Close()
		{
		if .qafter is false and .q isnt false
			.q.Close()
		else
			for cursor in Object(.qbefore, .qafter)
				if cursor isnt false
					cursor.Close()
		.q = .qafter = .qbefore = .field = false
		}
	}
