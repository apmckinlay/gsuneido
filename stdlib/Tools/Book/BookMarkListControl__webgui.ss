// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name: "BookMarkList"
	ComponentName: "BookMarkList"
	ComponentArgs: #()
	New()
		{
		.marks = Object()
		}

	GetState()
		{
		return .marks
		}

	SetState(stateobject)
		{
		.marks = Object()
		stateobject.Each()
			{
			if not String?(it)
				.marks.Add(it)
			else
				.marks.Add(Object(path: it, color: CLR.YELLOW))
			}
		.Act(#SetState, .marks)
		.GotoPath(.Controller.CurrentPage())
		}

	UpdateColor(mark)
		{
		if false is i = .findMark(mark.path)
			return
		.marks[i].color = mark.color
		}

	Click(path)
		{
		if false is i = .findMark(path)
			return
		bookmark = .marks[i].path
		if false is .Send('ValidBookmark?', bookmark)
			{
			.RemoveMark(bookmark)
			.AlertInfo('Bookmark Removed', 'Page ' $ bookmark $ ' not found.\n' $
				'Page may have been deleted or renamed.\nBookmark has been removed')
			}
		else
			.Send('Goto', bookmark)
		}

	DoubleClick()
		{
		.AddMark(.Controller.CurrentPage())
		}

	RemoveMark(path)
		{
		if false is x = .findMark(path)
			{
			.AlertInfo("Remove Mark", "There is no bookmark currently selected.")
			return false
			}
		.marks.Delete(x)
		.Act(#RemoveMark, path)
		}

	AddMark(path)
		{
		if path is "" or .findMark(path) isnt false
			return
		mark = Object(:path, color: CLR.YELLOW)
		.marks.Add(mark)
		.Act(#AddMark, mark)
		.GotoPath(path)
		}

	GotoPath(path)
		{
		.Act(#GotoPath, path)
		}

	findMark(path)
		{
		return .marks.FindIf({ it.path is path })
		}

	UpdateMarks(marks)
		{
		.marks = marks.DeepCopy()
		}
	}