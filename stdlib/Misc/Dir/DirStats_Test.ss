// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		_fs = #(
			'folder/': (
				(name: 'sub_folder_1/'),
				(name: 'sub_folder_2/'),
				(name: 'file_1', size: 10, date: #20240917),
				(name: 'file_2', size: 20, date: #20240916)
				(name: 'file_3', size: 5, date: #20240918))
			'folder/sub_folder_1/': (),
			'folder/sub_folder_2/': (
				(name: 'sub_folder_2/')),
			'folder/sub_folder_2/sub_folder_2/': (
				(name: 'file_1', size: 10, date: #20240917))
			)
		cl = DirStats
			{
			DirStats_dir(path, block)
				{
				_fs[path].Each(block)
				}
			}

		Assert(cl('folder/') is: #(
			(fileN: 3,
				largest: (name: "file_2", size: 20, date: #20240916),
				mostRecent: (name: "file_3", size: 5, date: #20240918),
				leastRecent: (name: "file_2", size: 20, date: #20240916),
				size: 45),
			"sub_folder_1/": ((fileN: 0, largest: false, mostRecent: false,
				leastRecent: false, size: 0)),
			"sub_folder_2/": ((fileN: 0, largest: false, mostRecent: false,
				leastRecent: false, size: 10),
				"sub_folder_2/": ((fileN: 1,
					largest: (name: "file_1", size: 10, date: #20240917),
					mostRecent: (name: "file_1", size: 10, date: #20240917),
					leastRecent: (name: "file_1", size: 10, date: #20240917),
					size: 10)))))
		}
	}