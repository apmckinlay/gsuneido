// Copyright (C) 2019 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	New()
		{
		super(.FixCommands())
		.Controller.Redir('On_Add_Bookmark', this)
		.Controller.Redir('On_Delete_Current_Bookmark', this)
		.static = .FindControl('Static')
		}

	Controls:
		(Vert
			(EtchedLine before: 0, after: 1 xstretch: 1)
			(HorzEqualHeight
				(EnhancedButton command: 'Add_Bookmark', image: 'plus.emf',
					imagePadding: 0.1, mouseEffect:, tip: 'Add a new bookmark')
				(EnhancedButton command: 'Delete_Current_Bookmark',
					image: 'minus.emf', imagePadding: 0.1, mouseEffect:,
					tip: 'Delete current bookmark')
				Skip
				(Static ' Bookmarks')
				Skip
			)
		)

	Resize(x, y, w, h)
		{
		.static.SetVisible(w > .Xmin)
		super.Resize(x, y, w, h)
		}

	FixCommands()
		{
		_parent.Window.SetupCommands(.Base())
		return .Controls
		}

	On_Add_Bookmark()
		{ .Parent.AddMark(.Controller.CurrentPage()) }

	On_Delete_Current_Bookmark()
		{ .Parent.RemoveMark(.Controller.CurrentPage()) }
	}