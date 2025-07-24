// Copyright (C) 2007 Suneido Software Corp. All rights reserved worldwide.
Controller
	{
	Title: 'Copy and Link Settings'
	New()
		{
		if false isnt settings = .query1()
			.Data.Set(settings)
		}
	Controls: (Record (Vert
		(Static 'Location to copy attachments to')
		(BrowseFolder name: copyto)
		(Static 'Note: This location should be somewhere accessible\n' $
			'to all users/workstations. i.e. a shared network directory')
		Skip
		(CheckBox 'Normally use Copy and Link' name: normally_linkcopy)
		Skip
		(CheckBox 'Delete source file after copying' name: delete_source)
		OkCancel
		))
	On_OK()
		{
		settings = .Data.Get()
		if settings.normally_linkcopy is true and settings.copyto.Blank?()
			{
			.AlertWarn(.Title, "Please enter a location to copy to.")
			return
			}
		if not settings.copyto.Blank?()
			{
			settings.copyto = Paths.EnsureTrailingSlash(Paths.ToStd(settings.copyto))
			if not .writeable?(settings.copyto)
				{
				.AlertWarn(.Title,
					"You do not seem to be able to save files in this directory.")
				return
				}
			}
		Database("ensure openimagesettings
			(copyto, normally_linkcopy, delete_source) key()")
		Transaction(update:)
			{ |t|
			t.QueryDo('delete openimagesettings')
			t.QueryOutput('openimagesettings', settings)
			}
		.Window.Result(settings)
		}
	writeable?(copyto)
		{
		testfile = copyto $ 'test.txt'
		try
			{
			PutFile(testfile, 'test')
			ok = GetFile(testfile) is 'test'
			DeleteFile(testfile)
			}
		catch
			ok = false
		return ok
		}
	Normally_linkcopy?()
		{
		settings = .query1()
		return settings isnt false and
			settings.normally_linkcopy is true and
			settings.copyto isnt ""
		}
	Copyto()
		{
		settings = .query1()
		return settings is false ? "" : settings.copyto
		}
	DeleteSource?()
		{
		settings = .query1()
		return settings is false ? false : (settings.delete_source is true)
		}
	query1()
		{
		LastContribution(#OpenImageSetting)()
		}
	}