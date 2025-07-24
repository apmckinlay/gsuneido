// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.query, .keyfields, .headerfields = false)
		{
		.query_columns = QueryColumns(.query)
		}

	GetKey()
		{
		return .keyfields
		}
	GetQuery()
		{
		return .query
		}
	GetHeaderFields()
		{
		return .headerfields
		}
	GetFields()
		{
		return .query_columns
		}

	Get(item)
		{
		Transaction(read:)
			{|t|
			item = .Getrecord(t, item)
			}
		return item
		}

	GetEntries()
		{
		Transaction(read:)
			{|t|
			if (false is q = t.Query(.query))
				throw "bad query: " $ .query
			entries = Object()
			while (false isnt entry =  q.Next())
				{
				entry.PreSet("Explorer_PreviousData", entry.Copy())
				entries.Add(entry)
				}
			}
		return entries
		}

	NewRecord()
		{
		return Record()
		}

	Output(item)
		{
		return KeyExceptionTransaction()
			{|t|
			t.QueryOutput(.query, item)
			item.PreSet("Explorer_PreviousData", item.Copy())
			return true
			}
		}

	Update(item)
		{
		return KeyExceptionTransaction()
			{ |t|
			result = true
			if (false is (update_item = .Getrecord(t, item.Explorer_PreviousData)))
				result = false
			else if RecordConflict?(
				item.Explorer_PreviousData, update_item, .query_columns)
				result = false
			else
				{
				update_item.Update(item)
				item.Delete("Explorer_PreviousData")
				item.PreSet("Explorer_PreviousData", item.Copy())
				}
			return result is true ? update_item : result
			}
		}

	DeleteItem(item)
		{
		return KeyExceptionTransaction()
			{|t|
			if false is delete_item = .Getrecord(t, item)
				return false
			delete_item.Delete()
			return true
			}
		}

	Getrecord(t, item)
		{
		q = Object(QueryStripSort(.query))
		.keyfields.Each({ q[it] = DatadictEncode(it, item[it]) })
		return t.Query1(@q)
		}
	ChangeQuery(query)
		{
		.query = query
		.query_columns = QueryColumns(.query)
		}
	Save(t /*unused*/ = false)
		{
		return true
		}
	}