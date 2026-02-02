// Copyright (C) 2025 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(dir1, dir2)
		{
		files = .dirFiles(dir1).MergeUnion(.dirFiles(dir2))
		results = Object(
			missing_dir1: Object(),
			missing_dir2: Object(),
			different: Object().Set_default(Object()),
			matching: 0,
			total: files.Size())
		for file in files
			.compareFiles(file, dir1, dir2, results)
		return results
		}

	dirFiles(dir)
		{
		return Dir(Paths.Combine(dir, '*.*'), files:)
		}

	compareFiles(file, dir1, dir2, results)
		{
		filePath1 = .filePath(file, dir1, results, 'missing_dir1')
		filePath2 = .filePath(file, dir2, results, 'missing_dir2')
		if filePath1 is '' or filePath2 is ''
			return
		if CompareFiles?(filePath1, filePath2)
			results.matching++
		else
			results.different[file.AfterLast('.').Lower()].Add(file)
		}

	filePath(file, dir, results, member)
		{
		if FileExists?(path = Paths.Combine(dir, file))
			return path
		results[member].Add(file)
		return ''
		}
	}