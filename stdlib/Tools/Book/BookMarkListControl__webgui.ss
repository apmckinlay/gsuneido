// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
Control
	{
	Name:          #BookMarkList
	ComponentName: #BookMarkList
	ComponentArgs: ()
	New(book)
		{
		.marks = BookMarks(book)
		.Act(#Sync, .marks.ListAll())
		}

	UpdateColor(mark)
		{
		if false isnt idx = .marks.Find(mark.path)
			.marks.SetColor(idx, mark.color)
		}

	Click(path)
		{
		if false is idx = .marks.Find(path)
			return
		if not .marks.Valid?(bookmark = .marks.Get(idx).path)
			{
			.RemoveMark(bookmark)
			.AlertInfo("Bookmark Removed",
				"Page " $ bookmark $ " not found.\n" $
					"Page may have been deleted or renamed.\nBookmark has been removed")
			}
		else
			.Send(#Goto, bookmark)
		}

	DoubleClick()
		{
		.AddMark(.Controller.CurrentPage())
		}

	RemoveMark(path)
		{
		if false is idx = .marks.Find(path)
			{
			.AlertInfo("Remove Mark", "There is no bookmark currently selected.")
			return false
			}
		.marks.Remove(idx)
		.Act(#RemoveMark, path)
		}

	AddMark(path)
		{
		if path is "" or .marks.Find(path) isnt false
			return
		.Act(#AddMark, .marks.Add(path))
		.GotoPath(path)
		}

	GotoPath(path)
		{
		.Act(#GotoPath, path)
		}

	UpdateMarks(marks)
		{
		.marks.Set(marks)
		}
	}
