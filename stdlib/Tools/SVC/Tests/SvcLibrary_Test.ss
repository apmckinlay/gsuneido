// Copyright (C) 2008 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_CheckRecord()
		{
		svcLibrary = SvcLibrary
			{
			SvcLibrary_checkCode(code/*unused*/, name/*unused*/, lib/*unused*/,
				results = false)
				{
				for x in _checkCode[1]
					results.Add([msg: x])
				return _checkCode[0]
				}
			CheckOnServer(name, lib/*unused*/)
				{
				return name is 'wrongName'
					? 'invalid fizz'
					: ''
				}
			}
		svcTable = SvcTable(table = .MakeLibrary())
		.EnsureLibrary(table)

		svcTable.Output([name: #folder, parent: 0, group: 0])
		svcTable.Output([name: #one, parent: 1])
		svcTable.Output([name: #two, parent: 1])
		svcTable.Output([name: #nested_folder, parent: 1, group: 1])
		svcTable.Output([name: #Test_SuppressedRecord, parent: 0, group: -1,
			text: `function () { /*Do nothing*/ }`])
		svcTable.Output([
			name: table.Capitalize() $ #_CheckLibrarySuppressions,
			parent: 0,
			group: -1,
			text: `#(Test_SuppressedRecord)`])

		svcLib = svcLibrary(table)
		f = {|name, type| svcLib.CheckRecord(name, type) }
		_checkCode = [true, []]
		rec = Query1(table, name: 'one')
		Assert(f(rec, '-') is: #())
		Assert(f(rec, '+') is: #())

		_checkCode = ["", [w = "WARNING: bad foo"]]
		rec = Query1(table, name: 'two')
		Assert(f(rec, '-') is: #())
		Assert(f(rec, '+') is: Object(w))

		_checkCode = [false, [w = "WARNING: bad foo", e = "ERROR: invalid bar"]]
		Assert(f(rec, '-') is: #())
		Assert(f(rec, '+') equalsSet: Object(w, e))

		rec.text = '// BuiltDate > 20990101\r\n\r\n function () { fred = barney }'
		Assert(f(rec, '-') is: #())
		Assert(f(rec, '+') is: #())

		_checkCode = [true, []]
		rec.name = 'wrongName'
		rec.text = 'if true {i =5 }; Print(i)'
		Assert(f(rec, '') is: #('invalid fizz'))

		rec.name = 'Test_SuppressedRecord'
		Assert(f(rec, '') is: #('suppressed'))
		}

	Test_checkName()
		{
		fn = SvcLibrary.SvcLibrary_checkName

		// no message if lib is not in use or doesn't exit
		Assert(fn('FakeRecord', results = Object(), 'fakelib'))
		Assert(results is: #())

		// has ERROR prefix if lib is in use
		Assert(fn('FakeRecord', results = Object(), 'stdlib') is: false)
		Assert(results is: #((msg: "ERROR: can't find: FakeRecord")))

		Assert(fn('FakeRecord', results = Object(#(msg: 'pre-existing msg')), 'stdlib')
			is: false)
		Assert(results
			is: #((msg: 'pre-existing msg'), (msg: "ERROR: can't find: FakeRecord")))

		// check web file names
		Assert(fn('FakeRecord.css', results = Object(), 'stdlib'))
		Assert(results is: #())
		Assert(fn('FakeRecord.js', results = Object(), 'stdlib'))
		Assert(results is: #())

		// web record name when lib not in use or doesn't exist
		Assert(fn('FakeRecord.css', results = Object(), 'fakelib'))
		Assert(results is: #())

		// invalid record name
		Assert(fn('FakeRecord.notweb', results = Object(), 'stdlib') is: false)
		expected = #((msg: "ERROR: can't find: FakeRecord.notweb"))
		Assert(results is: expected)
		}
	}
