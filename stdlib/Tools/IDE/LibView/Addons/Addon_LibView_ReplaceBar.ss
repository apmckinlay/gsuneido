// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
LibViewAddon
	{
	Commands(cmds)
		{
		cmds.Add(
			#(Replace, 'Ctrl+H', 'Find and replace text in the current item'),
			#(Replace_Current, 'F8')
			)
		}

	saveReplaceData: #()
	Getter_Initialized()
		{ return .View isnt false and .Editor isnt false and .container isnt false }

	getter_container()
		{ return .View.FindControl(.View.BottomLeft) }

	Explorer_DeselectTab(oldView)
		{
		oldCtrl = .getReplaceBar(oldView)
		newCtrl = .getReplaceBar(.container)
		if oldCtrl is false and newCtrl is false
			return
		else if oldCtrl is false
			.View.Send(#On_FindBar_Close, skipSave?:)
		else
			.setCtrl(oldCtrl, .View.Send(#Collect, #On_Replace)[0])
		}

	getReplaceBar(control)
		{ return control.FindControl(#ReplaceBar) }

	setCtrl(oldCtrl, newCtrl)
		{
		newCtrl.Data.Set(oldCtrl.Data.Get())
		newCtrl.SetText(.Editor.FindReplaceData().replace = oldCtrl.GetText())
		}

	On_Replace()
		{
		.FindBar_OpenFindBar()
		return .openReplaceBar()
		}

	openReplaceBar()
		{
		if false is ctrl = .getReplaceBar(.container)
			{
			.saveReplaceData = .saveReplaceData.Empty?()
				? .Editor.FindReplaceData()
				: .saveReplaceData
			ctrl = .container.Append(['ReplaceBar', .saveReplaceData])
			}
		ctrl.Select()
		return ctrl
		}

	On_FindBar_Close(skipSave? = false)
		{
		if not skipSave? // Preserve the old settings, FindBar prompts close command
			.saveReplaceData = .Editor.FindReplaceData().Copy()
		}

	On_Replace_Current()
		{
		.setReplace()
		if not .Editor.ReplaceOne()
			Beep()
		.Editor.On_Find_Next()
		.Editor.SetFocus()
		}

	setReplace()
		{
		if false isnt replaceBar = .getReplaceBar(.container)
			.Editor.FindReplaceData().replace = replaceBar.GetText()
		}

	On_ReplaceBar_ReplaceAll()
		{ .replace('Entire Text') }

	replace(replaceIn)
		{
		.setReplace()
		.Editor.FindReplaceData().replaceIn = replaceIn
		.Editor.ReplaceAll()
		.Editor.SetFocus()
		}

	On_ReplaceBar_ReplaceCurrent()
		{ .On_Replace_Current() }

	On_ReplaceBar_ReplaceInSelection()
		{ .replace('Selection') }
	}
