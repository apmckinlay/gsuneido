// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Addon_VirtualListViewBase
	{
	On_Context_Go_To_QueryView()
		{
		.Model.GoToQueryView()
		}

	On_Context_Inspect_Control()
		{
		.Grid.KEYDOWN(VK.F8)
		}

	On_Context_Copy_Field_Name()
		{
		VirtualListDev.CopyFieldName(.GetContextMenu().ContextCol)
		}

	On_Context_Go_To_Field_Definition()
		{
		VirtualListDev.GoToFieldDefinition(.GetContextMenu().ContextCol)
		}

	On_Context_Print()
		{
		ReportFromQueryAndColumns(.Model.GetQuery(), .GetColumns(), hwnd: .GetGridHwnd())
		}

	On_Context_Reporter()
		{
		Reporter()
		}

	On_Context_Reason_Protected()
		{
		rec = .GetContextMenu().ContextRec
		query = .Model.GetQuery()
		ListCustomize.ReasonProtected(
			rec, .Model.EditModel.ProtectField , .GetGridHwnd(), query)
		}

	On_Context_Reset_Columns()
		{
		.Header.ResetColumns(.Model.GetQuery(), .GetPrimarySort())
		}

	On_Context_Customize_Columns()
		{
		.Model.ColModel.CustomizeColumns(
			.Parent, .Model.GetQuery(), .Model.EditModel.Editable?())
		.Grid.ScrollToLeft()
		}

	On_Context_Customize()
		{
		if .SaveFirst() is false
			return true
		if 0 is subTitle = .Send('VirtualList_GetSubTitle')
			subTitle = ''
		if .Model.ColModel.Customize(.Parent, .Model.GetQuery(), :subTitle)
			.Send('BookRefresh')
		}

	On_Context_Customize_Expand()
		{
		if .SaveFirst() is false
			return true
		query = .Model.GetQuery()
		if 0 is defaultExpandLayout = .Send('VirtualList_DefaultExpandLayout')
			defaultExpandLayout = ''
		if .Model.ExpandModel.Customize(query, .ExpandColumns(), defaultExpandLayout,
			.GetAccessCustomKey())
			.Send('BookRefresh')
		}

	On_Context_Global(item)
		{
		recMenu = .GetContextRecordMenu()
		if recMenu isnt false
			recMenu.On_Global(item[0], .GetGridHwnd())
		}

	On_Context_Set_as_Default_Sort()
		{
		.Model.SetDefaultSort()
		}

	On_Context_Reset_Sort_to_System_Default()
		{
		if .SaveFirst() is false
			return
		.ClearSelect()
		.Model.ResetSort()
		.Header.RefreshSort(.GetPrimarySort())
		.Grid.Repaint()
		}
	}
