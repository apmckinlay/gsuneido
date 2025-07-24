// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
class
	{
	EnsureTable()
		{
		Database("ensure user_notes
			(name, text, last_modified, last_modified_by, path, usernote_TS)
			key(name, path)")
		}
	Query(title, cur_book_option)
		{
		return "user_notes where name = " $ Display(title) $
			' and path is ' $ Display(cur_book_option)
		}
	}