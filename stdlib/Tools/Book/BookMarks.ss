// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(.book, .user = false)
		{
		.key = .Key(.book)
		.Sync()
		}

	Key(book)
		{
		return "BookMarks-" $ book
		}

	Sync()
		{
		.marks = .load()
		.modifyMarks({ it.DeleteIf({|idx| not .Valid?(it[idx].path) }) })
		}

	load()
		{
		return UserSettings.Get(.key, Object(), user: .user)
		}

	modifyMarks(block)
		{
		block(.marks)
		.save(.marks)
		}

	save(marks)
		{
		return UserSettings.Put(.key, marks, user: .user)
		}

	Valid?(path)
		{
		return not QueryEmpty?(.book,
			path: path.BeforeLast('/'), name: path.AfterLast('/'))
		}

	Add(path)
		{
		if false isnt idx = .Find(path)
			return .marks[idx]
		mark = Object(:path, color: CLR.YELLOW)
		.modifyMarks({ it.Add(mark) })
		return mark
		}

	Remove(idx) // Cannot name: "Delete" as that overlaps with the built-in class methods
		{
		.modifyMarks({ it.Delete(idx) })
		}

	SetColor(idx, color)
		{
		if .marks.Member?(idx)
			.modifyMarks({ it[idx].color = color })
		}

	SetPath(idx, path)
		{
		if .marks.Member?(idx)
			.modifyMarks({ it[idx].path = path })
		}

	Swap(idx1, idx2)
		{
		if .marks.Member?(idx1) and .marks.Member?(idx2)
			.modifyMarks({ it.Swap(idx1, idx2) })
		}

	Set(marks)
		{
		if not Object?(marks)
			return
		.save(marks)
		.Sync()
		}

	Get(idx)
		{
		return .marks.GetDefault(idx, false)
		}

	Find(path)
		{
		return .marks.FindIf({ it.path is path })
		}

	Count() // Cannot name: "Size" as that overlaps with the built-in class methods
		{
		return .marks.Size()
		}

	ListAll() // Used by suneido.js to set the web components
		{
		return .marks
		}
	}
