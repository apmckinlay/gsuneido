// Copyright (C) 2004 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.prefix, .n = 100, .lib = 'configlib', .checkOnly = false)
		{
		EnsureConfigLib(.lib)
		ServerEval('Use', .lib)
		}
	NeedFunc(text)
		{
		if .checkOnly
			return .getFunc(text)
		RetryTransaction()
			{ |t|
			return .needFunc(text, t)
			}
		}
	getFunc(text)
		{
		Transaction(read:)
			{ |t|
			r = .getRange(t)
			return .find(text, r, t)
			}
		}
	needFunc(text, t)
		{
		r = .getRange(t)
		if false is name = .find(text, r, t)
			name = .create(text, r, t)
		return name
		}
	getRange(t)
		{
		r = t.QueryRange(.query(), 'num')
		if r is false
			r = Object(min: 0, max: 0)
		return r
		}
	topOfStack: 30
	find(text, r, t)
		{
		names = Object()
		// find the Top of the stack
		t.QueryApply(.query() $ ' sort reverse num', readonly:)
			{
			names.Add(it.name)
			if names.Size() >= .topOfStack
				break
			}
		t.QueryApply(.query() $ ' and Adler32(text) is ' $ Adler32(text))
			{ |x|
			if x.text isnt text
				continue

			// want to push the record back to the "top of the stack"
			// so that it does not inadvertantly get deleted by other
			// formulas in the same report, BUT -
			// don't push to the top if it is already there
			if names.Has?(x.name)
				return x.name
			if not .checkOnly
				{
				name = .buildName(r.max, t)
				x.name = name
				x.num = t.QueryMax(.lib, 'num', 0) + 1
				x.Update()
				return x.name
				}
			}
		return false
		}
	recycle(t)
		{
		if t.QueryCount(.query()) < .n
			return

		i = 0
		t.QueryApply(.query() $ ' sort reverse num') // keep n number records
			{
			if not it.name.Prefix?(.prefix) or ++i < .n
				continue

			it.Delete()
			Unload(it.name)
			}
		}
	create(text, r, t)
		{
		.recycle(t)
		name = .buildName(r.max, t)
		t.QueryOutput(.lib,
			Record(num: t.QueryMax(.lib, 'num', 0) + 1, group: -1, parent: 0, :name,
				:text))
		return name
		}
	query()
		{
		return .lib $ ' where name > ' $ Display(.prefix) $
			' and name <= ' $ Display(.prefix $ String(.max))
		}
	max: 999999
	nameWidth: 6
	buildName(i, t)
		{
		rec = t.Query1(.lib, num: i)
		if rec is false
			return .prefix $ 1.Pad(.nameWidth)

		number = Number(rec.name.AfterFirst(.prefix)) + 1
		name = .prefix $ (number).Pad(.nameWidth)
		while (not QueryEmpty?(.lib, :name))
			name = .prefix $ (++number).Pad(.nameWidth)

		if Number(name.AfterFirst(.prefix)) >= .max
			name = .prefix $ '000001'

		return name
		}
	}