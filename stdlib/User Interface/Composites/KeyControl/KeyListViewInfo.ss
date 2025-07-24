// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
// TODO: rename table, class, fields to something more generic
class
	{
	Save(id, window_info)
		{
		.tryCatchNotAuthorized()
			{
			RetryTransaction()
				{ |t|
				t.QueryDo('delete ' $ .query(id))
				t.QueryOutput('keylistview_info',
					[user: Suneido.User, query: id, :window_info])
				}
			}
		}
	Get(id)
		{
		.tryCatchNotAuthorized()
			{
			return Query1(.query(id))
			}
		}
	DeleteRecord(id)
		{
		.tryCatchNotAuthorized()
			{
			QueryDo('delete ' $ .query(id))
			}
		}
	tryCatchNotAuthorized(block)
		{
		try
			block()
		catch (unused, "not authorized (from server)")
			{
			return false
			}
		}
	query(id)
		{
		Database('ensure keylistview_info
			(user, query, window_info, keylistview_TS)
			key(user, query)')
		return 'keylistview_info
			where user = ' $ Display(Suneido.User) $
			' and query = ' $ Display(id)
		}
	}