// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_getCustomTabs()
		{
		mock = Mock(CustomizeButtonControl)
		mock.When.getCustomTabs([anyArgs:]).CallThrough()

		mock.Window = FakeObject(FindControl: false)
		Assert(mock.getCustomTabs(Object()) is: false)

		mock.Window = FakeObject(FindControl: Object(TabName: false))
		Assert(mock.getCustomTabs(Object()) is: #(Header))

		mock.Window = FakeObject(FindControl: Object(TabName: #TestName))
		Assert(mock.getCustomTabs(Object()) is: #(TestName))

		tabs = Object(.tabsMock(#General, #Custom1))
		Assert(result = mock.getCustomTabs(tabs) equalsSet: #(TestName, General, Custom1))
		Assert(result[0] is: #TestName)

		mock.Window = FakeObject(FindControl: Object(TabName: false))
		tabs.Add(.tabsMock(#Custom2, #General))
		result = mock.getCustomTabs(tabs)
		Assert(result equalsSet: #(Header, General, Custom1,  Custom2))
		Assert(result[0] is: #Header)

		mock.Window = FakeObject(FindControl: false)
		Assert(mock.getCustomTabs(tabs) equalsSet: #(General, Custom1,  Custom2))

		mock.Window = FakeObject(FindControl: Object(TabName: false))
		tabs.Add(.tabsMock(#Header))
		Assert(result = mock.getCustomTabs(tabs)
			equalsSet: #(Header, General, Custom1,  Custom2))
		Assert(result[0] is: #Header)
		}

	tabsMock(@tabs)
		{
		tabMock = Mock()
		tabMock.When.CustomizableTabs().Return(tabs)
		return tabMock
		}
	}