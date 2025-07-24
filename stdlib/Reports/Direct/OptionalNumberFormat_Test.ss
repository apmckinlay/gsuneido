// Copyright (C) 2016 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_ExportCSV()
		{
		f = new OptionalNumberFormat()
		Assert(f.ExportCSV(20.15) is: '20.15')

		f = new OptionalNumberFormat(data: 20000.175, mask: '###,###.#')
		// TextFormat always use fmt.Data
		Assert(f.ExportCSV(10.222) is: '20000.2')

		f = new OptionalNumberFormat(mask: '###,###.#')
		Assert(f.ExportCSV([]) is: '.0')

		f = new OptionalNumberFormat()
		Assert(f.ExportCSV(20.1334) is: '20.1334')

		f = new OptionalNumberFormat(mask: '###,###.#')
		Assert(f.ExportCSV('') is: '')

		f = new OptionalNumberFormat(mask: '###,###.#')
		Assert(f.ExportCSV(20000.17) is: '20000.2')

		f = new OptionalNumberFormat(mask: '###,###.#')
		Assert(f.ExportCSV(1e999) is: '')

		// not sure if this is intentional
		f = new OptionalNumberFormat()
		Assert(f.ExportCSV([]) is: 'false')
		}
	}
