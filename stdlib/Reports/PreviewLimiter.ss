// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Singleton
	{
	maxWindows: 6
	New()
		{
		.queue = Object()
		}

	Add(hwnd)
		{
		if .queue.Has?(hwnd)
			{
			.queue.Remove(hwnd)
			.queue.Add(hwnd)
			return
			}

		.queue.Add(hwnd)
		}

	Remove(hwnd)
		{
		.queue.Remove(hwnd)
		}

	sizeOverLimit?()
		{
		return .queue.Size() >= .maxWindows
		}

	BeforeOpen()
		{
		if not .sizeOverLimit?()
			return true

		hwnd = .queue.PopFirst()
		if not .windowEnabled?(hwnd)
			{
			.queue.Add(hwnd, at: 0)
			.alert()
			return false
			}
		.destroyWindow(hwnd)
		return true
		}

	// extracted for testing
	windowEnabled?(hwnd)
		{
		return IsWindowEnabled(hwnd)
		}

	// extracted for testing
	alert()
		{
		Alert("Maximum preview windows reached.\n\n" $
			"Please close some of your preview windows.",
			title: 'Report Preview', flags: MB.ICONWARNING)
		}

	// extracted for testing
	destroyWindow(hwnd)
		{
		DestroyWindow(hwnd)
		}
	}
