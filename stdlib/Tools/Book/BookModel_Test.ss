// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		Suneido.Delete('BookModels')
		x = BookModel('no_table')
		x.BookModel_children = [
			'': [.data[0], .data[4]],
			'/Business': [.data[1], .data[2], .data[3]],
			'/Payables': [.data[5]]]
		Assert(x.Get('/Business/one') is: .data[1])
		Assert(x.Children('/Business') is: .data[1 :: 3])
		Assert(x.Children?(''))
		Assert(x.Children?('/Business'))
		Assert(x.Children?('/Business/two') is: false)
		for i in .. .data.Size()
			{
			Assert(x.Prev(.data[i])
				is: i > 0 ? .data[i - 1] : false)
			Assert(x.Next(.data[i])
				is: i + 1 < .data.Size() ? .data[i + 1] : false)
			Assert(x.Parent(.data[i]) is: .data.GetDefault(.p[i], false))
			}
		}
	Test_sort()
		{
		Suneido.Delete('BookModels')
		toc = Object().Set_default(#())
		for x in .data.Copy().Sort!(By(#num))
			toc[x.path].Add(x)
		x = BookModel('no_table')
		x.BookModel_children = toc
		Assert(.data is: x.BookModel_toc)
		}
	p: (false, 0, 0, 0, false, 4)
	data: (
		{path: '', name: 'Business', num: 3},
		{path: '/Business', name: 'one', order: 1, num: 2},
		{path: '/Business', name: 'two', order: 2, num: 5},
		{path: '/Business', name: 'three', order: 3, num: 6},
		{path: '', name: 'Payables', num: 4},
		{path: '/Payables', name: 'p1', num: 1},
		)
	}