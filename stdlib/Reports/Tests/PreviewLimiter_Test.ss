// Copyright (C) 2021 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	Test_main()
		{
		cl = PreviewLimiter
			{
			PreviewLimiter_maxWindows: 3
			CallClass()
				{
				new this
				}
			New()
				{
				.PreviewLimiter_queue = Object()
				.PreviewLimiter_destroyed = Object()
				.PreviewLimiter_disabledWindows = Object()
				}
			PreviewLimiter_windowEnabled?(hwnd)
				{
				return not .PreviewLimiter_disabledWindows.Has?(hwnd)
				}
			PreviewLimiter_alert() {}
			PreviewLimiter_destroyWindow(hwnd)
				{
				.PreviewLimiter_destroyed.Add(hwnd)
				}
			Remove(hwnd)
				{
				super.Remove(hwnd)
				.PreviewLimiter_destroyed.Add(hwnd)
				}
			}
		ins = cl()
		ins.PreviewLimiter_disabledWindows.Add(444)
		// test adding to one below limit
		ins.Add(111)
		Assert(ins.BeforeOpen())
		Assert(ins.PreviewLimiter_queue is: #(111))
		Assert(ins.PreviewLimiter_destroyed is: #())
		ins.Add(222)
		Assert(ins.BeforeOpen())
		Assert(ins.PreviewLimiter_queue is: #(111, 222))
		Assert(ins.PreviewLimiter_destroyed is: #())
		ins.Add(111)
		Assert(ins.BeforeOpen())
		Assert(ins.PreviewLimiter_queue is: #(222, 111))
		Assert(ins.PreviewLimiter_destroyed is: #())

		// test adding to the limit
		ins.Add(444)
		Assert(ins.BeforeOpen())
		Assert(ins.PreviewLimiter_destroyed is: #(222))
		Assert(ins.PreviewLimiter_queue is: #(111, 444))
		ins.Add(111)
		Assert(ins.BeforeOpen())
		Assert(ins.PreviewLimiter_destroyed is: #(222))
		Assert(ins.PreviewLimiter_queue is: #(444, 111))
		ins.Add(333)
		Assert(ins.BeforeOpen() is: false)
		Assert(ins.PreviewLimiter_destroyed is: #(222))
		Assert(ins.PreviewLimiter_queue is: #(444, 111, 333))
		// test trying to delete again but the deletion should fail again
		Assert(ins.BeforeOpen() is: false)
		Assert(ins.PreviewLimiter_destroyed is: #(222))
		Assert(ins.PreviewLimiter_queue is: #(444, 111, 333))

		// test deleting a disabled window
		ins.Remove(444)
		Assert(ins.BeforeOpen())
		Assert(ins.PreviewLimiter_destroyed is: #(222, 444))
		Assert(ins.PreviewLimiter_queue is: #(111, 333))

		// test adding one more window
		ins.Add(555)
		Assert(ins.BeforeOpen())
		Assert(ins.PreviewLimiter_destroyed is: #(222, 444, 111))
		Assert(ins.PreviewLimiter_queue is: #(333, 555))
		}
	}