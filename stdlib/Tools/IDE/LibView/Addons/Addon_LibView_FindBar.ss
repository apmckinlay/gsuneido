// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
LibViewAddon
	{
	Commands(cmds)
		{
		cmds.Add(
			#(Find, 'Ctrl+F', 'Find text in the current item', seq: 1.0),
			#(Find_Next, 'F3', 'Find the next occurrence in the current item'),
			#(Find_Previous, 'Shift+F3',
				'Find the previous occurrence in the current item'),
			#(Find_Next_Selected, 'Ctrl+F3',
				'Find the next occurrence of the selected text'),
			#(Find_Prev_Selected, 'Shift+Ctrl+F3',
				'Find the previous occurrence of the selected text')
			)
		}

	saveFindData: #()
	Init()
		{
		.Redir('On_Find_Next', .Editor)
		.Redir('On_Find_Previous', .Editor)
		.Redir('On_Find_Next_Selected', .Editor)
		.Redir('On_Find_Prev_Selected', .Editor)
		}

	Getter_Initialized()
		{ return .View isnt false and .Editor isnt false and .container isnt false }

	getter_container()
		{ return .View.FindControl(.View.BottomLeft) }

	tabsChanging?: false
	Explorer_DeselectTab(oldView)
		{
		oldCtrl = .getFindBar(oldView)
		newCtrl = .getFindBar(.container)
		if oldCtrl is false and newCtrl is false
			return
		else if oldCtrl is false
			.View.Send(#On_FindBar_Close, skipSave?:)
		else
			{
			.tabsChanging? = true
			.setCtrl(oldCtrl, .View.Send(#Collect, #FindBar_OpenFindBar)[0])
			}
		}

	setCtrl(oldCtrl, newCtrl)
		{
		.setFindData(oldCtrl.Data.Get(), newCtrl)
		newCtrl.SetText(oldCtrl.GetText())
		}

	setFindData(data, setCtrl = false)
		{
		.saveFindData = data.Copy()
		for field in data.Members().Remove(#replace) // Allows replacebar to set itself
			{
			if setCtrl isnt false
				setCtrl.Data.SetField(field, data[field])
			.Editor.FindReplaceData()[field] = data[field]
			}
		}

	getFindBar(control)
		{ return control.FindControl(#FindBar) }

	Find_Change()
		{
		if .tabsChanging?
			{
			.tabsChanging? = false
			return
			}
		.Editor.SetSelect(.Editor.GetSelect().cpMin)
		if false isnt findbar = .getFindBar(.container)
			findbar.SetStatus(.Editor.On_Find_Next())
		}

	Find_Return()
		{ .Editor.SetFocus() }

	FindBar_OpenFindBar()
		{
		if false is ctrl = .getFindBar(.container)
			{
			ctrl = .container.Append([#FindBar, .Editor.FindReplaceData()])
			.setFindData(.saveFindData.Empty?()
				? .Editor.FindReplaceData()
				: .saveFindData)
			}
		s = .Editor.GetSelText()
		if s > "" and not s.Has?('\n')
			.Editor.FindReplaceData().find = s
		ctrl.Select()
		return ctrl
		}

	On_Find()
		{ .FindBar_OpenFindBar() }

	On_FindBar_Clear()
		{
		.Editor.ClearFindMarks()
		.Editor.SendToAddons(#Overview_Reset)
		}

	On_FindBar_Close(skipSave? = false)
		{
		if not skipSave? // Preserve the old settings
			.saveFindData = .Editor.FindReplaceData().Copy()
		.closeCtrl(ReplaceBarControl)
		.closeCtrl(FindBarControl)
		.On_FindBar_Clear()
		.Editor.SetFocus()
		}

	closeCtrl(ctrl)
		{
		if false isnt idx = .container.GetChildren().FindIf({ it.Base?(ctrl) })
			.container.Remove(idx)
		}

	On_FindBar_Mark()
		{
		.Editor.MarkAll()
		.Editor.SendToAddons(#Overview_Reset)
		}

	On_FindBar_Next()
		{
		.Editor.SetFocus()
		.Editor.On_Find_Next()
		}

	On_FindBar_Previous()
		{
		.Editor.SetFocus()
		.Editor.On_Find_Previous()
		}

	UpdateOccurrence(num = false, count = false)
		{
		if false is findbar = .getFindBar(.container)
			return
		findbar.UpdateOccurrenceInfo(num, count)
		findbar.UpdateOccurrenceMsg()
		}
	}
