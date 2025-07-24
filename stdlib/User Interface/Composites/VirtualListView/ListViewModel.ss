// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(query, reverse = false)
		{
		.reverse = reverse
		.trans = Transaction(read:)
		try
			{
			.q_beg = .trans.SeekQuery(query)
			.q_end = .trans.Query(query)
			}
		catch (x)
			{
			.Close()
			.fields = Object()
			.end = -1
			throw (x)
			}
		.fields = .q_beg.Columns()
		.beg = .beg_pos = 0
		n = .trans.QueryCount(query)
		.end = .end_pos = (n > 0 ? n - 1 : 0)
		.midpoint = (.end / 2).Int()
		.current_record1 = .qbeg_next()
		.current_record2 = .qend_prev()

		//check for empty query and set .end to -1 so Getnumrows returns 0
		if (.current_record1 is false)
			.end = -1
		.recnum = 0
		.rec = .Getrecord(.recnum)
		}
	Getcolumns()
		{
		return .fields
		}
	Getnumrows()
		{
		return .end + 1
		}
	Getitem(recnum, fld)
		{
		if (.recnum isnt recnum)
			.rec = .Getrecord(.recnum = recnum)
		if (.rec is Object() or .rec is false)
			return ""
		return .rec[fld]
		}
	Getrecord(recnum)
		{
		if .trans is false or .trans.Ended?() or recnum < .beg or recnum > .end
			return Record()
		if (recnum <= .midpoint)
			{
			// use beg_pos / q_beg
			if (recnum < (.beg_pos - recnum).Abs())
				{
				.q_beg.Rewind()
				.beg_pos = 0
				.current_record1 = .qbeg_next()
				}
			diff = .beg_pos - recnum
			if (diff > 0)
				{
				for (i = 0; i < diff; ++i)
					{
					--.beg_pos
					Assert(false isnt (x = .qbeg_prev()))
					.current_record1 = x
					}
				return .current_record1
				}
			else if (diff < 0)
				{
				diff *= -1
				for (i = 0; i < diff; ++i)
					{
					++.beg_pos
					Assert(false isnt (x = .qbeg_next()))
					.current_record1 = x
					}
				return .current_record1
				}
			else
				{
				return .current_record1
				}
			}
		else
			{
			// use end_pos / q_end
			if (.end - recnum < (.end_pos - recnum).Abs())
				{
				.q_end.Rewind()
				.end_pos = .end
				.current_record2 = .qend_prev()
				}
			diff = .end_pos - recnum
			if (diff > 0)
				{
				for (i = 0; i < diff; ++i)
					{
					--.end_pos
					Assert(false isnt (x = .qend_prev()))
					.current_record2 = x
					}
				return .current_record2
				}
			else if (diff < 0)
				{
				diff *= -1
				for (i = 0; i < diff; ++i)
					{
					++.end_pos
					Assert (false isnt (x = .qend_next()))
					.current_record2 = x
					}
				return .current_record2
				}
			else
				{
				return .current_record2
				}
			}
		}
	qbeg_next() { return .reverse ? .q_beg.Prev() : .q_beg.Next() }
	qbeg_prev() { return .reverse ? .q_beg.Next() : .q_beg.Prev() }
	qend_next() { return .reverse ? .q_end.Prev() : .q_end.Next() }
	qend_prev() { return .reverse ? .q_end.Next() : .q_end.Prev() }
	Seek(field, prefix)
		{
		if .trans is false or .trans.Ended?()
			return false

		.beg_pos = .q_beg.Seek(field, prefix)
		.current_record1 = .q_beg.Next()
		return .beg_pos
		}
	q_beg: false
	q_end: false
	Close()
		{
		if .q_beg isnt false
			.q_beg.Close()
		if .q_end isnt false
			.q_end.Close()
		if (.trans isnt false and not .trans.Ended?())
			.trans.Complete()
		.trans = .rec = .q_beg = .q_end = .current_record1 = .current_record2 = false
		}
	}
