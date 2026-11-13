// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Setup()
		{
		super.Setup()

		.book = .MakeBook()
		.rec1 = .MakeBookRecord(.book, #NA)
		.rec2 = .MakeBookRecord(.book, #NA, path: "/Test2")
		.rec3 = .MakeBookRecord(.book, #NA, path: "/Test/Test3")

		.user = .TempName()
		.key = BookMarks.Key(.book)

		.userSettingsExists? = TableExists?(UserSettings.Table)
		}

	Test_main()
		{
		marks = BookMarks(.book, .user)
		Assert(marks.BookMarks_key is: .key)

		// Test / build book paths
		Assert(marks.Valid?(rec1Path = .rec1.path $ '/' $ .rec1.name))
		marks.Add(rec1Path)
		Assert(marks.Valid?(rec2Path = .rec2.path $ '/' $ .rec2.name))
		marks.Add(rec2Path)
		// Add will simply add the path as directed (Valid? is only triggered on Load)
		Assert(marks.Valid?(invalidPath1 = "/NotARealPath/OrRecord1") is: false)
		marks.Add(invalidPath1)
		Assert(marks.Valid?(rec3Path = .rec3.path $ '/' $ .rec3.name))
		marks.Add(rec3Path)
		marks.Add(rec3Path) // Attempt to add a duplicate path, (should be ignored)
		// Add will simply add the path as directed (Valid? is only triggered on Load)
		Assert(marks.Valid?(invalidPath2 = "/NotARealPathOrRecord2") is: false)
		marks.Add(invalidPath2)
		Assert(marks.Count() is: 5)

		// Test "invalid" interactions (shouldn't throw errors or modify the marks)
		marks.SetColor(-1, CLR.BLACK)
		marks.SetColor(10, CLR.BLACK)
		marks.SetPath(10, "/FAKE")
		marks.Swap(0, 10)
		marks.Swap(-1, 3)
		marks.Swap(-1, 10)
		marks.Set(false)
		marks.Remove(-1)
		marks.Remove(10)

		// Verify marks current state prior to adjustments
		Assert(list = marks.ListAll() isSize: 5)
		.assertMarkState(0, list, marks, rec1Path, CLR.YELLOW)
		.assertMarkState(1, list, marks, rec2Path, CLR.YELLOW)
		.assertMarkState(2, list, marks, invalidPath1, CLR.YELLOW)
		.assertMarkState(3, list, marks, rec3Path, CLR.YELLOW)
		.assertMarkState(4, list, marks, invalidPath2, CLR.YELLOW)

		// Test various interaction / adjustment calls
		marks.SetColor(0, CLR.plum) // rec1Path
		marks.SetColor(3, CLR.RED) // rec3Path
		marks.Swap(0, 3) // Swap: rec1Path and rec3Path
		Assert(marks.Add(rec3Path).color is: CLR.RED) // Duplicate path, return original
		marks.Remove(1) // Remove: rec2Path
		marks.SetPath(0, rec2Path) // Replace: rec3Path, with: rec2Path
		Assert(marks.Count() is: 4) // 1 mark is manually removed
		Assert(marks.Add(rec2Path).color is: CLR.RED) // Duplicate path, return original

		// Reload the BookMarks, verifying the save and on-load valid
		marks = BookMarks(.book, .user)
		Assert(marks.BookMarks_key is: .key)
		Assert(marks.Count() is: 2) // 1 mark invalid mark is removed on load
		Assert(list = marks.ListAll() isSize: 2)
		.assertMarkState(0, list, marks, rec2Path, CLR.RED)
		.assertMarkState(1, list, marks, rec1Path, CLR.plum)

		// Clear the list by setting it to and empty object and verify the state
		marks.Set([])
		Assert(marks.Count() is: 0)
		Assert(marks.ListAll() is: #())
		Assert(marks.Get(0) is: false)
		Assert(marks.Find(rec1Path) is: false)

		// Ensure the above Set was successfully saved
		marks = BookMarks(.book, .user)
		Assert(marks.BookMarks_key is: .key)
		Assert(marks.Count() is: 0)
		Assert(marks.ListAll() is: #())
		}

	assertMarkState(idx, list, marks, path, color)
		{
		Assert(list[idx] is: marks.Get(idx))
		Assert(list[idx].path is: path)
		Assert(list[idx].color is: color)
		Assert(marks.Find(path) is: idx)
		}

	Teardown()
		{
		UserSettings.RemoveAllUsers(.key)
		if not .userSettingsExists? and TableExists?(UserSettings.Table)
			Database("destroy " $ UserSettings.Table)
		super.Teardown()
		}
	}
