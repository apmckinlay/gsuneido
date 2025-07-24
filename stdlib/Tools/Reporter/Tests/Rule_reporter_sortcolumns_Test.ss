// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		rec = Record()
		Assert(rec.reporter_sortcolumns is: #())

		rec.design_cols = #('col1' 'col2' 'col3')
		mock = Mock()
		mock.When.Reporter_BuildSortColsList_unsortable?(rec, 'col1').Return(false)
		mock.When.Reporter_BuildSortColsList_unsortable?(rec, 'col2').Return(true)
		mock.When.Reporter_BuildSortColsList_unsortable?(rec, 'col3').Return(false)
		cols = mock.Eval(Reporter_BuildSortColsList, rec)
		Assert(cols is: #('col1', 'col3'))
		}

	Test_sort_column_order()
		{
		rec = Record()
		rec.columns = Object(#(text: col1) #(text: col2) #(text: col3))
		rec.design_cols = Object('col0', 'col3', 'col2', 'col5', 'col1')
		Reporter_BuildSortColsList.
			Reporter_BuildSortColsList_move_designcols_to_start(rec, rec.design_cols)
		Assert(rec.design_cols is: #('col3', 'col2', 'col1', '', 'col0', 'col5'))
		}
	}
