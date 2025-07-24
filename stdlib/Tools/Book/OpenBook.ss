// Copyright (C) 2003 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(bookname, pageinfo = false, bookedit? = false, path = false)
		{
		if not TableExists?(bookname)
			return false
		booktype = .open_book(bookname, bookedit?)
		.goto_page(booktype, bookname, bookedit?, pageinfo, path)
		}
	open_book(bookname, bookedit?)
		{
		booktype = bookedit? is true ? 'EditBooks' : 'OpenBooks'

		.PreOpenBook(booktype, bookname)
		if Suneido.Member?(booktype) and Suneido[booktype].Member?(bookname)
			{
			hwnd = Suneido[booktype][bookname].Window.Hwnd
			if not IsWindowEnabled(hwnd)
				EnableWindow(hwnd, true)
			WindowActivate(hwnd)
			}
		else if bookedit? is true
			BookEditControl(bookname)
		else
			PersistentWindow(Object('Book', bookname, bookname, help_book:),
				option: bookname)
		return booktype
		}

	PreOpenBook(@unused) { }

	goto_page(booktype, bookname, bookedit?, pageinfo, path)
		{
		if pageinfo is false and path is false
			return false
		book = Suneido[booktype][bookname]
		// try context sensitive
		if bookedit? is false
			{
			if path isnt false
				book.Goto(path)
			else
				{
				page = BookPageFind(bookname, pageinfo.path, pageinfo.name)
				if page isnt false
					book.Goto(page.path $ "/" $ page.name)
				}
			}
		else
			book.Explorer.GotoPath(pageinfo)
		}
	}
