// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
Test
	{
	Test_main()
		{
		contributionBase = class
			{
			New(checkVal)
				{ .checkVal = checkVal }

			Call(data)
				{
				warnings = Object()
				for val in data[.checkVal]
					warnings.Add(val)
				return warnings
				}
			}
		contrib1 = new contributionBase('contrib1')
		contrib2 = new contributionBase('contrib2')
		contrib3 = new contributionBase('contrib3')
		.SpyOn(Contributions).Return([contrib2, contrib1, contrib3])

		data = []
		Assert(WarningStatus(data, 'Warning0:', '') is: '')

		data.contrib3 = [[priority: 50, msg: 'contrib3']]
		Assert(WarningStatus(data, 'Warning1:', '') is: 'Warning1: contrib3')

		data.contrib1 = [[priority: 1, msg: 'contrib1']]
		Assert(WarningStatus(data, 'Warning2:', '')
			is: 'Warning2: contrib1; and contrib3')

		data.contrib2 = [
			[priority: 10, msg: 'contrib2.1'],
			[priority: 100, msg: 'contrib2.2']
			]
		Assert(WarningStatus(data, 'Warning3:', '', lastJoin: 'also')
			is: 'Warning3: contrib1; contrib2.1; contrib3; also contrib2.2')

		data.contrib1 = #()
		Assert(WarningStatus(data, 'Warning4:', '', lastJoin: 'or')
			is: 'Warning4: contrib2.1; contrib3; or contrib2.2')

		data.contrib2 = #()
		Assert(WarningStatus(data, 'Warning5:', '', lastJoin: 'NA')
			is: 'Warning5: contrib3')

		data.contrib3 = #()
		Assert(WarningStatus(data, 'Warning6:', '') is: '')
		}

	Test_NoContributions()
		{
		// Ensuring that there are no issues if no contributions are returned
		.SpyOn(Contributions).Return([])
		Assert(WarningStatus([], 'Warning', '') is: '')
		}
	}