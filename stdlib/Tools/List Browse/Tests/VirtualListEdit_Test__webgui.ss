// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
_VirtualListEdit_Test
	{
	Test_editNextCell()
		{
		mock = Mock(VirtualListEdit)
		mock.When.editNextCell([anyArgs:]).CallThrough()
		mock.When.editCell([anyArgs:]).Return('')

		mock.VirtualListEdit_model = Mock(VirtualListModel)
		mock.VirtualListEdit_rec = []
		colModel = mock.VirtualListEdit_model.ColModel = Mock(VirtualListColModel)
		grid = mock.VirtualListEdit_parent = Mock()
		grid.When.Send("VirtualListGrid_SaveRecord", []).Return(true)
		grid.When.InsertRow([anyArgs:]).Return(true)

		mock.When.nextColEditable('b', 1, []).Return(true)

		grid.When.GetSelectedRecord().Return([])

		colModel.When.GetColumns().Return(#(a, b, c))
		colModel.When.FindCol('a').Return(0)
		colModel.When.Get(1).Return('b')

		mock.editNextCell('a', 1)
		mock.Verify.editCell('b')

		colModel.When.FindCol('b').Return(1)
		colModel.When.Get(2).Return('c')
		mock.When.nextColEditable('c', 2, []).Return(true)
		mock.editNextCell('b', 1)
		mock.Verify.editCell('c')

		colModel.When.FindCol('c').Return(2)
		colModel.When.Get(0).Return('a')
		grid.When.SelectNextRow(1).Return(true)
		mock.When.nextColEditable('a', 0, []).Return(true)
		mock.When.rowChanged([anyArgs:]).CallThrough()
		mock.editNextCell('c', 1)
		mock.Verify.editCell('a')

		mock.editNextCell('c', -1)
		mock.Verify.Times(2).editCell('b')

		mock.editNextCell('b', -1)
		mock.Verify.Times(2).editCell('a')

		grid.When.SelectNextRow(-1).Return(true)
		mock.When.nextColEditable('a', 0, []).Return(false)
		mock.When.nextColEditable('b', 1, []).Return(false)
		mock.editNextCell('c', -1)
		mock.Verify.Times(2).editCell('c')

		grid.When.SelectNextRow(-1).Return(false)
		mock.editNextCell('c', -1)
		mock.Verify.Times(2).editCell('c') // no change

		grid.When.SelectNextRow(1).Return(false)
		mock.editNextCell('c', 1)
		mock.Verify.Times(2).editCell('c') // no change
		}
	}