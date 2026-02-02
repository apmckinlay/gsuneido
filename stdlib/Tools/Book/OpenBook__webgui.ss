// Copyright (C) 2024 Axon Development Corporation All rights reserved worldwide.
_OpenBook
	{
	PreOpenBook(booktype, bookname)
		{
		if SuRenderBackend().WindowManager.ShowingModalWindow?() and
			Suneido.Member?(booktype) and Suneido[booktype].Member?(bookname)
			{
			window = Suneido[booktype][bookname].Window
			if not Same?(SuRenderBackend().WindowManager.GetActive(), window)
				DestroyWindow(window.Hwnd)
			}
		}
	}