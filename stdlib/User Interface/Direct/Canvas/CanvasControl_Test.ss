// Copyright (C) 2017 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_moveItem()
		{
		mock = Mock(CanvasControl)

		mock = Mock()
		mock.When.GetWidth().Return(1000, 200)
		mock.When.GetHeight().Return(1000, 200)
		mock.CanvasControl_selected = Object()
		mock.Eval(CanvasControl.MoveSelected, false, false)
		mock.Verify.Never().Send("CanvasChanged")
		mock.Verify.Never().InvalidateClientArea()

		item1 = Mock()
		item2 = Mock()
		item1.When.BoundingRect().Return([x1: 10, x2: 100, y1: 10, y2: 100])
		item2.When.BoundingRect().Return([x1: 110, x2: 200, y1: 110, y2: 200])
		selected = mock.CanvasControl_selected = Object(item1, item2)
		mock.Eval(CanvasControl.MoveSelected, 1, 1)
		selected[0].Verify.Move(1, 1)
		selected[1].Verify.Move(1, 1)
		mock.Verify.Send("CanvasChanged")
		mock.Verify.InvalidateClientArea()

		mock.Eval(CanvasControl.MoveSelected, 1, 1)
		selected[0].Verify.Move(0, 0)
		selected[1].Verify.Move(0, 0)
		mock.Verify.Times(2).Send("CanvasChanged")
		mock.Verify.Times(2).InvalidateClientArea()
		}
	}