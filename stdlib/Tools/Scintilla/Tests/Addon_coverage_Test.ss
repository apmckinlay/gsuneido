// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_rollBackBlanks()
		{
		roll = Addon_coverage.Addon_coverage_rollBackBlanks
		Assert(roll("hello\r\n\t\tworld", 9, 0) is: 7)
		Assert(roll("hello\r\n\t\tworld", 0, 0) is: 0)
		Assert(roll("hello\r\n\t\twor ld", 13, 0) is: 12)
		Assert(roll("hello\r\n\t\twor ld", 3, 0) is: 3)
		Assert(roll("hello\r\n\t\t\t\tworld", 11, 9) is: 9)
		}

	Test_getMarginText()
		{
		getMarginText = Addon_coverage.Addon_coverage_getMarginText
		Assert(getMarginText(0) is: '0')
		Assert(getMarginText(20) is: '20')
		Assert(getMarginText(65534) is: '65534')
		Assert(getMarginText(65535) is: '>=64k')
		Assert(getMarginText(64.Kb()) is: '>=64k')
		}

	Test_startCoverageAndTest()
		{
		_covered = Object()
		addon = Addon_coverage
			{
			Addon_coverage_startCoverage(lib, name)
				{
				_covered.Add(lib $ ':' $ name)
				}
			}
		fn = addon.Addon_coverage_startCoverageAndTest
		lib = 'Test_lib'
		name = .TempName()
		fn(lib, name)
		Assert(_covered is: Object('Test_lib:' $ name))

		.MakeLibraryRecord([:name, text: `class { }`])
		.MakeLibraryRecord([name: name $ '_Test', text: `class { }`])
		_covered = Object()
		fn(lib, name)
		Assert(_covered is: Object('Test_lib:' $ name, 'Test_lib:' $ name $ '_Test'))

		.MakeLibraryRecord([name: name $ 'Test', text: `class { }`])
		_covered = Object()
		fn(lib, name)
		Assert(_covered is:
			Object(lib $ ':' $ name,
				'Test_lib:' $ name $ '_Test',
				'Test_lib:' $ name $ 'Test'))

		_covered = Object()
		fn(lib, name $ 'Test')
		Assert(_covered is:
			Object('Test_lib:' $ name $ 'Test',
				'Test_lib:' $ name))

		_covered = Object()
		fn(lib, name $ '_Test')
		Assert(_covered is:
			Object('Test_lib:' $ name $ '_Test',
				'Test_lib:' $ name))
		}

	Test_skip?()
		{
		fn = Addon_coverage.Addon_coverage_skip?
		Assert(fn('', ''))

		name = .TempName()
		Assert(fn('Test_lib', name))

		.MakeLibraryRecord([:name, text: `class { }`])
		Assert(fn('Test_lib', name) is: false)

		// builtin
		Assert(fn('stdlib', 'Object?'))
		Assert(fn('stdlib', 'Objects') is: false)
		Assert(fn('stdlib', 'CLR')) // object
		Assert(fn('stdlib', 'MAX_PATH')) // value
		}
	}
