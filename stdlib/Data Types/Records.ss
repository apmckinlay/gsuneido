// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	DoWithTran(block, update = false)
		{
		if false isnt (t = (.Member?(#recordTran) ? this[#recordTran] : false)) or
			false isnt (t = .Transaction())
			block(t)
		else
			Transaction(:update, :block)
		}
	SafeMembers()
		{
		// for Records this needs to copy because even read-only access
		// can modify them due to rules
		// NOTE members created during iteration will NOT be included
		return .Members().Copy()
		}
	Project(@fields)
		{
		// same as Objects but makes a Record not an Object
		// and doesn't set non-existent members
		if fields.Size(list:) is 1 and Object?(fields[0])
			fields = fields[0]
		ob = Record()
		for f in fields
			if .Member?(f)
				ob[f] = this[f]
		return ob
		}
	Query1(@args)
		{
		.DoWithTran({|t| return t.Query1(@args) })
		}
	Query1Cached(@args)
		{
		.DoWithTran({|t| return t.Query1Cached(@args) })
		}
	QueryEmpty?(@args)
		{
		.DoWithTran({|t| return t.QueryEmpty?(@args) })
		}
	QueryFirst(@args)
		{
		.DoWithTran({|t| return t.QueryFirst(@args) })
		}
	QueryLast(@args)
		{
		.DoWithTran({|t| return t.QueryLast(@args) })
		}
	}
