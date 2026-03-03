// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_Controls()
		{
		mock = .mockControl()
		mock.When.Controls().CallThrough()

		// Empty data / graph
		mock.BulletGraphsControl_data = #()
		mock.BulletGraphsControl_heading = ''
		mock.BulletGraphsControl_controlContainer = 'Vert'
		mock.BulletGraphsControl_initDisplayDetails = false
		Assert(mock.Controls() is: Object('Vert', name: 'graphs'))

		// Minimal Vertical Graph
		mock.BulletGraphsControl_data = #(
			#20260101: 0,
			#20260201: 200,
			#20260301: 400,
			#20260401: 1000,
			#20260501: 300)
		Assert(mock.Controls() is: #('Vert',
			#('Vert', 'Skip',
				#('Horz',
					#('BulletGraph', 0, satisfactory: 0, vertical:, height: 100,
						outside: 0, axisFormat: '#,###.##', name: #20260101, good: 0,
						target: false, range: #(0, 1100), color: 0x4d4d4d, width: 25,
						axis:, axisDensity: 3, selectedColor: 16711680)
					#('BulletGraph', 200, satisfactory: 190, vertical:, height: 100,
						outside: 0, axisFormat: '#,###.##', name: #20260201, good: 380,
						target: false, range: #(0, 1100), color: 0x666622, width: 25,
						axis: false, axisDensity: 3, selectedColor: 16711680),
					#('BulletGraph', 400, satisfactory: 190, vertical:, height: 100,
						outside: 0, axisFormat: '#,###.##', name: #20260301, good: 380,
						target: false, range: #(0, 1100), color: 0x666622, width: 25,
						axis: false, axisDensity: 3, selectedColor: 16711680),
					#('BulletGraph', 1000, satisfactory: 190, vertical:, height: 100,
						outside: 0, axisFormat: '#,###.##', name: #20260401, good: 380,
						target: false, range: #(0, 1100), color: 0x666622, width: 25,
						axis: false, axisDensity: 3, selectedColor: 16711680),
					#('BulletGraph', 300, satisfactory: 190, vertical:, height: 100,
						outside: 0, axisFormat: '#,###.##', name: #20260501, good: 380,
						target: false, range: #(0, 1100), color: 0x666622, width: 25,
						axis: false, axisDensity: 3, selectedColor: 16711680),
					),
					#('Horz', 'Fill', #('Static', ' ', justify: 'RIGHT',
						name: 'displayDetails'))), name: 'graphs'))
		mock.Verify.Never().displayDetailsSend([anyArgs:])

		// Horizontal Graph with heading and initDisplayDetails
		mock.BulletGraphsControl_heading = 'Test Heading'
		mock.BulletGraphsControl_controlContainer = 'Horz'
		mock.BulletGraphsControl_graphsContainer = 'Vert'
		mock.BulletGraphsControl_displayDetailsSpacer = 'Skip'
		mock.BulletGraphsControl_initDisplayDetails = 'Last'
		mock.When.displayDetailsSend([anyArgs:]).Return('Display Details')
		mock.BulletGraphsControl_data = #(
			#20260101: -10,
			#20260201: 0,
			#20260301: 0,
			#20260401: 0,
			#20260501: 0)
		Assert(mock.Controls() is: #('Vert',
			#('Heading', 'Test Heading'),
			#('Horz', 'Skip',
				#('Vert',
					#('BulletGraph', 0, satisfactory: 0, vertical:, height: 100,
						outside: 0, axisFormat: '#,###.##', name: #20260101, good: 0,
						target: false, range: #(0, 1), color: 0x5e3838, width: 25,
						axis:, axisDensity: 1, selectedColor: 16711680)
					#('BulletGraph', 0, satisfactory: 0, vertical:, height: 100,
						outside: 0, axisFormat: '#,###.##', name: #20260201, good: 0,
						target: false, range: #(0, 1), color: 0x4d4d4d, width: 25,
						axis: false, axisDensity: 1, selectedColor: 16711680),
					#('BulletGraph', 0, satisfactory: 0, vertical:, height: 100,
						outside: 0, axisFormat: '#,###.##', name: #20260301, good: 0,
						target: false, range: #(0, 1), color: 0x4d4d4d, width: 25,
						axis: false, axisDensity: 1, selectedColor: 16711680),
					#('BulletGraph', 0, satisfactory: 0, vertical:, height: 100,
						outside: 0, axisFormat: '#,###.##', name: #20260401, good: 0,
						target: false, range: #(0, 1), color: 0x4d4d4d, width: 25,
						axis: false, axisDensity: 1, selectedColor: 16711680),
					#('BulletGraph', 0, satisfactory: 0, vertical:, height: 100,
						outside: 0, axisFormat: '#,###.##', name: #20260501, good: 0,
						target: false, range: #(0, 1), color: 0x4d4d4d, width: 25,
						axis: false, axisDensity: 1, selectedColor: 16711680),
					),
					#('Horz', 'Skip', #('Static', 'Display Details', justify: 'RIGHT',
						name: 'displayDetails'))), name: 'graphs'))
		mock.Verify.displayDetailsSend('construct', #20260501)
		}

	mockControl()
		{
		mock = Mock(BulletGraphsControl)
		mock.BulletGraphsControl_target = false
		mock.BulletGraphsControl_min = 0
		mock.BulletGraphsControl_vertical = true
		mock.BulletGraphsControl_height = 100
		mock.BulletGraphsControl_width = 25
		mock.BulletGraphsControl_axisDensity = 3
		mock.BulletGraphsControl_axisFormat = '#,###.##'
		mock.BulletGraphsControl_graphsContainer = 'Horz'
		mock.BulletGraphsControl_displayDetailsSpacer = 'Fill'
		mock.BulletGraphsControl_good = 0x226322
		mock.BulletGraphsControl_satisfactory = 0x666622
		mock.BulletGraphsControl_bad = 0x883322
		mock.BulletGraphsControl_inactive = 0x4d4d4d
		mock.BulletGraphsControl_negative = 0x5e3838
		return mock
		}

	Test_graphsBase()
		{
		mock = .mockControl()
		mock.When.graphsBase([anyArgs:]).CallThrough()

		data = [a: 10, b: 10]
		graphsBase = mock.graphsBase(false, data, 10, 3)
		Assert(graphsBase.range is: #(10, 11))
		Assert(graphsBase.good is: 0)
		Assert(graphsBase.satisfactory is: 0)
		Assert(graphsBase.color is: 0x4d4d4d)
		Assert(graphsBase.axisDensity is: 1)

		graphsBase = mock.graphsBase(false, data, 0, 3)
		Assert(graphsBase.range is: #(0, 20))
		Assert(graphsBase.good is: 10)
		Assert(graphsBase.satisfactory is: 5)
		Assert(graphsBase.color is: 0x226322)
		Assert(graphsBase.axisDensity is: 3)

		data = [a: 10, b: 10, c: 100]
		graphsBase = mock.graphsBase(false, data, 0, 3)
		Assert(graphsBase.range is: #(0, 110))
		Assert(graphsBase.good is: 40)
		Assert(graphsBase.satisfactory is: 20)
		Assert(graphsBase.color is: 0x883322)
		Assert(graphsBase.axisDensity is: 3)
		}

	Test_calcMax()
		{
		m = BulletGraphsControl.BulletGraphsControl_calcMax

		Assert(m(false, [a: 501, b: 701, c: 1001], 0) is: 1100)
		Assert(m(false, [a: 501, b: 9901, c: 701], 0) is: 10000)
		Assert(m(10001, [a: 501, b: 9901, c: 701], 0) is: 11000)
		Assert(m(10000, [a: 10000], 0) is: 11000)
		Assert(m(10000, [a: 10000], 10000) is: false)
		Assert(m(0, [a: 0], 0) is: false)
		Assert(m(0, [a: 0], 1) is: false)
		}

	Test_determineColor()
		{
		mock = .mockControl()
		mock.When.determineColor([anyArgs:]).CallThrough()

		data = [a: 10, b: 20, c: 30, d: 40, e: 50]
		Assert(mock.determineColor(data, 10, 0) is: 0x226322)
		Assert(mock.determineColor(data, 35, 20) is: 0x666622)
		Assert(mock.determineColor(data, 40, 35) is: 0x883322)
		}

	Test_Set()
		{
		mock = .mockControl()
		mock.When.Set([anyArgs:]).CallThrough()
		mock.BulletGraphsControl_graphs = graphsMock = Mock()
		graphsMock.When.RemoveAll().Do({ })
		graphsMock.When.AppendAll().Do({ })

		// Graphs container is empty, and no graphs are built
		mock.BulletGraphsControl_selected = 1
		mock.BulletGraphsControl_hover = 2
		graphsMock.When.GetChildren().Return(Object())
		mock.When.graphControls().Return(false)
		mock.Set(Object(/* data */))
		Assert(mock.BulletGraphsControl_selected is: false)
		Assert(mock.BulletGraphsControl_hover is: false)
		graphsMock.Verify.Never().RemoveAll()
		graphsMock.Verify.Never().AppendAll([anyArgs:])

		// Graphs container has controls, and graphs are built
		mock.BulletGraphsControl_selected = 3
		mock.BulletGraphsControl_hover = 4
		graphsMock.When.GetChildren().Return(Object('fake children'))
		mock.When.graphControls().Return(expectedGraphControls = Object('fake controls'))
		mock.Set(Object(/* data */))
		Assert(mock.BulletGraphsControl_selected is: false)
		Assert(mock.BulletGraphsControl_hover is: false)
		graphsMock.Verify.RemoveAll()
		graphsMock.Verify.AppendAll(expectedGraphControls)
		}

	Test_BulletGraph_Hover()
		{
		mock, ctrlMock = .userMocks('hover')

		mock.BulletGraph_Hover([Name: 'label0'])
		mock.Verify.setDisplayDetails('hover', 'label0')
		ctrlMock.Verify.Never().Set([anyArgs:])

		mock.BulletGraph_Hover([Name: 'label1'])
		mock.Verify.setDisplayDetails('hover', 'label1')
		ctrlMock.Verify.Set('value1')

		mock.BulletGraphsControl_hover = 'label2'
		mock.BulletGraph_Hover([Name: 'label2'])
		mock.Verify.Never().setDisplayDetails('hover', 'label2')
		ctrlMock.Verify.Never().Set('value2')
		}

	userMocks(event)
		{
		mock = Mock(BulletGraphsControl)
		mock.When['BulletGraph_' $ event.Capitalize()]([anyArgs:]).CallThrough()
		ctrlMock = Mock()
		ctrlMock.When.Set([anyArgs:]).Do({ })
		mock.When.FindControl('displayDetails').Return(ctrlMock)

		mock.BulletGraphsControl_data = [label0: 1, label1: 5, label2: 9]
		mock.When.Send('BulletGraphs_DisplayDetails', event, 'label0', 1).Return(0)
		mock.When.Send('BulletGraphs_DisplayDetails', event, 'label1', 5).Return('value1')
		mock.When.Send('BulletGraphs_DisplayDetails', event, 'label2', 9).Return('value2')
		return mock, ctrlMock
		}

	Test_BulletGraph_Click()
		{
		mock, ctrlMock = .userMocks('click')
		mock.When.FindControl('label0').Return(graphMock0 = .graphMock())
		mock.When.FindControl('label1').Return(graphMock1 = .graphMock())
		mock.When.FindControl('label2').Return(graphMock2 = .graphMock())

		mock.BulletGraph_Click([Name: 'label0'])
		mock.Verify.setDisplayDetails('click', 'label0')
		graphMock0.Verify.Never().Selected([anyArgs:])
		ctrlMock.Verify.Never().Set([anyArgs:])

		mock.BulletGraph_Click([Name: 'label1'])
		mock.Verify.setDisplayDetails('click', 'label1')
		graphMock1.Verify.Selected(true)
		ctrlMock.Verify.Set('value1')

		mock.BulletGraphsControl_selected = 'label2'
		mock.BulletGraph_Click([Name: 'label2'])
		mock.Verify.Never().setDisplayDetails('click', 'label2')
		graphMock2.Verify.Never().Selected([anyArgs:])
		ctrlMock.Verify.Never().Set('value2')
		}

	graphMock()
		{
		graphMock = Mock()
		graphMock.When.Selected([anyArgs:]).Do({})
		return graphMock
		}
	}