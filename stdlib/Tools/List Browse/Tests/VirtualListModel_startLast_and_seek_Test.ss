// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
VirtualListModelTests
	{
	Test_not_all_read_reverse()
		{
		model = VirtualListModel(.VirtualList_Table, startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)
		Assert(model.Offset is: -19)
		Assert(model.VisibleRows is: 20)

		Assert(model.Seek('num', 25) is: 0)
		Assert(model.GetStartLast() is: false)
		}

	Test_not_all_read()
		{
		model = VirtualListModel(.VirtualList_Table)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(20)
		Assert(model.Offset is: 0)
		Assert(model.VisibleRows is: 20)

		Assert(model.Seek('num', 25) is: 0)
		Assert(model.GetStartLast() is: false)
		}

	Test_all_read_reverse()
		{
		model = VirtualListModel(.VirtualList_Table, startLast:)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(40)
		Assert(model.Offset is: -30)
		Assert(model.VisibleRows is: 40)

		Assert(model.Seek('num', 25) is: -5)
		Assert(model.Seek('num', 20) is: -10)
		Assert(model.Seek('num', 40) is: -1)
		Assert(model.Seek('num', -1) is: -30)
		Assert(model.GetStartLast())
		}

	Test_all_read()
		{
		model = VirtualListModel(.VirtualList_Table)
		.AddTeardownModel(model)
		model.UpdateVisibleRows(40)
		Assert(model.Offset is: 0)
		Assert(model.VisibleRows is: 40)

		Assert(model.Seek('num', 25) is: 25)
		Assert(model.GetStartLast() is: false)
		}

	Test_seekAllRead_oddSizedObjects()
		{
		mock = Mock(VirtualListModel)
		mock.When.seekAllRead([anyArgs:]).CallThrough()
		mock.VirtualListModel_data = Object()

		// Test: Empty .data
		Assert(mock.seekAllRead('test', 10) is: 0)

		// Test: .data is populated, startLast is false (positive indexes)
		mock.VirtualListModel_data = Object(
			[test: 10],
			[test: 20],
			[test: 30],
			[test: 40],
			[test: 50])
		.assertSeek(mock, -100, 0, 10)
		.assertSeek(mock, 5, 0, 10)
		.assertSeek(mock, 10, 0, 10)
		.assertSeek(mock, 12, 1, 20)
		.assertSeek(mock, 20, 1, 20)
		.assertSeek(mock, 23, 2, 30)
		.assertSeek(mock, 30, 2, 30)
		.assertSeek(mock, 35, 3, 40)
		.assertSeek(mock, 40, 3, 40)
		.assertSeek(mock, 44, 4, 50)
		.assertSeek(mock, 50, 4, 50)
		.assertSeek(mock, 54, 4, 50)
		.assertSeek(mock, 60, 4, 50)
		.assertSeek(mock, 100, 4, 50)

		// Test: .data is populated, startLast is true (negative indexes)
		mock.VirtualListModel_data = Object(
			-5: [test: 10],
			-4: [test: 20],
			-3: [test: 30],
			-2: [test: 40],
			-1: [test: 50])
		.assertSeek(mock, -100, -5, 10)
		.assertSeek(mock, 5, -5, 10)
		.assertSeek(mock, 10, -5, 10)
		.assertSeek(mock, 20, -4, 20)
		.assertSeek(mock, 30, -3, 30)
		.assertSeek(mock, 35, -2, 40)
		.assertSeek(mock, 40, -2, 40)
		.assertSeek(mock, 50, -1, 50)
		.assertSeek(mock, 100, -1, 50)
		}

	assertSeek(mock, seek, expectedIdx, expectedValue)
		{
		Assert(mock.seekAllRead('test', seek) is: expectedIdx, msg: msg = 'Seek: ' $ seek)
		Assert(mock.VirtualListModel_data[expectedIdx].test is: expectedValue, :msg)
		}

	Test_seekAllRead_evenSizedObjects()
		{
		mock = Mock(VirtualListModel)
		mock.When.seekAllRead([anyArgs:]).CallThrough()
		mock.VirtualListModel_data = Object()

		// Test: Empty .data
		Assert(mock.seekAllRead('test', 10) is: 0)

		// Test: .data is populated, startLast is false (positive indexes)
		mock.VirtualListModel_data = Object(
			[test: 10],
			[test: 20],
			[test: 30],
			[test: 40],
			[test: 50],
			[test: 60])
		.assertSeek(mock, -100, 0, 10)
		.assertSeek(mock, 5, 0, 10)
		.assertSeek(mock, 10, 0, 10)
		.assertSeek(mock, 20, 1, 20)
		.assertSeek(mock, 30, 2, 30)
		.assertSeek(mock, 35, 3, 40)
		.assertSeek(mock, 40, 3, 40)
		.assertSeek(mock, 45, 4, 50)
		.assertSeek(mock, 50, 4, 50)
		.assertSeek(mock, 60, 5, 60)
		.assertSeek(mock, 100, 5, 60)

		// Test: .data is populated, startLast is true (negative indexes)
		mock.VirtualListModel_data = Object(
			-5: [test: 10],
			-4: [test: 20],
			-3: [test: 30],
			-2: [test: 40],
			-1: [test: 50],
			0: [test: 60])
		.assertSeek(mock, -100, -5, 10)
		.assertSeek(mock, 5, -5, 10)
		.assertSeek(mock, 10, -5, 10)
		.assertSeek(mock, 11, -4, 20)
		.assertSeek(mock, 20, -4, 20)
		.assertSeek(mock, 30, -3, 30)
		.assertSeek(mock, 35, -2, 40)
		.assertSeek(mock, 40, -2, 40)
		.assertSeek(mock, 50, -1, 50)
		.assertSeek(mock, 55, 0, 60)
		.assertSeek(mock, 60, 0, 60)
		.assertSeek(mock, 100, 0, 60)
		}

	Test_seekAllRead_strings()
		{
		mock = Mock(VirtualListModel)
		mock.When.seekAllRead([anyArgs:]).CallThrough()
		mock.VirtualListModel_data = Object()

		mock.VirtualListModel_data = Object([test: 'AAA'])
		Assert(mock.seekAllRead('test', '') is: 0)
		Assert(mock.seekAllRead('test', 'A') is: 0)
		Assert(mock.seekAllRead('test', 'AAA') is: 0)

		mock.VirtualListModel_data = Object(
			[test: 'Border Crossing Loaded'],
			[test: 'D.O.T. Inspections'],
			[test: 'Driver Backhaul'],
			[test: 'Driver Training In Truck'],
			[test: 'Driver Training Out Of Truck'],
			[test: 'Drug Test'],
			[test: 'Extend A-Train Poles'],
			[test: 'Layover'],
			[test: 'Loading  - 0 - 30 kms'],
			[test: 'Loading - 31 - 75 kms'],
			[test: 'Loading - 76 kms & up'],
			[test: 'O/Op - Fuel Sur-Charge - 5axle 37%'],
			[test: 'O/Op - Fuel Sur-Charge - R-Tac 37%'],
			[test: 'O/Op - Fuel Sur-Charge - Sec 37%'],
			[test: 'O/Op - Fuel Sur-Charge - Triaxle 37%'],
			[test: 'O/Op - Split Load'],
			[test: 'O/Op - Unloading - 0-18 Miles'],
			[test: 'O/Op - Unloading - 19-46 Miles'],
			[test: 'O/Op - Unloading - 47-92 Miles'],
			[test: 'O/Op Load Light Deduction'],
			[test: 'O/Op Return load'],
			[test: 'Reload'],
			[test: 'Reset Away from Brandon'],
			[test: 'Split Load'],
			[test: 'Trailer Switch'],
			[test: 'Unloading - 0 - 30 kms'],
			[test: 'Unloading - 31 - 75 kms'],
			[test: 'Unloading - 76 kms & up'],
			[test: 'Used Oil Load'],
			[test: 'Used Oil Unload'],
			[test: 'Waiting Time'],
			[test: 'Washout PTO Pump'],
			[test: 'Washout Trailers']) // Index: 32
		.assertSeek(mock, 'A', 0, 'Border Crossing Loaded')
		.assertSeek(mock, 'Bo', 0, 'Border Crossing Loaded')
		.assertSeek(mock, 'D', 1, 'D.O.T. Inspections')
		.assertSeek(mock, 'Dr', 2, 'Driver Backhaul')
		.assertSeek(mock, 'Driver T', 3, 'Driver Training In Truck')
		.assertSeek(mock, 'Dru', 5, 'Drug Test')
		.assertSeek(mock, 'E', 6, 'Extend A-Train Poles')
		.assertSeek(mock, 'H', 7, 'Layover')
		.assertSeek(mock, 'Loading  ', 8, 'Loading  - 0 - 30 kms')
		.assertSeek(mock, 'Loading -', 9, 'Loading - 31 - 75 kms')
		.assertSeek(mock, 'M', 11, 'O/Op - Fuel Sur-Charge - 5axle 37%')
		.assertSeek(mock, 'O/Op - S', 15, 'O/Op - Split Load')
		.assertSeek(mock, 'P', 21, 'Reload')
		.assertSeek(mock, 'Reload', 21, 'Reload')
		.assertSeek(mock, 'Res', 22, 'Reset Away from Brandon')
		.assertSeek(mock, 'S', 23, 'Split Load')
		.assertSeek(mock, 'T', 24, 'Trailer Switch')
		.assertSeek(mock, 'U', 25, 'Unloading - 0 - 30 kms')
		.assertSeek(mock, 'Un', 25, 'Unloading - 0 - 30 kms')
		.assertSeek(mock, 'Us', 28, 'Used Oil Load')
		.assertSeek(mock, 'Washout T', 32, 'Washout Trailers')
		.assertSeek(mock, 'Washout Ts', 32, 'Washout Trailers')
		.assertSeek(mock, 'Z', 32, 'Washout Trailers')
		}
	}
