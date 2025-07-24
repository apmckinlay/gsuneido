// Copyright (C) 2006 Suneido Software Corp. All rights reserved worldwide.
// contributed by Claudio Mascioni
// WARNING: not implemented in suneido.js
Controller
	{
	Name: 'BrowseImageName'
	Xmin: 420 // width
	Ymin: 315 // height
	// file may contain only file name or a complete path with file name
	// or a complete path without file name
	CallClass(hwnd = 0, title = "", filter = "", file = false,
		showimage = true, opendirmsg = "")
		{
		return ToolDialog(hwnd, Object(this, title, filter, file,
			showimage, opendirmsg), border: 2)
		}

	New(title = "", filter = "", file = false, showimage = true, opendirmsg = "")
		{
		super(.controls())
		.Title = title is ""
			? TranslateLanguage('Select an image')
			: title
		.Redir('On_Undo', this)
		.set(.initialfile = file)
		.opendirmsg = TranslateLanguage(opendirmsg)
		.chlist = .Vert.Hzl.Vert.ChooseList
		.lbox = .Vert.Hzl.ListBox
		filter = filter isnt ""
			? filter.Lower()
			: "bmp, gif, jpg, jpe, jpeg, ico, emf, wmf, tif, tiff, png, pdf"
		.filterlist = filter.Split(", ").Sort!()
		.filterlist.Add(TranslateLanguage('all') $ ' (*.*)')
		.chlist.SetList(.filterlist)
		// set the initial filter
		if .getFileExt(.filename) isnt ""
			{
			.imagetype = .getFileExt(.filename)
			.chlist.Set(.imagetype)
			}
		else
			{
			.imagetype = .filterlist[0]
			.chlist.SelectItem(0)
			}
		.Vert.Hzb.CheckBox.Set(showimage)
		.imageslist(.filename)
		}
	Commands:
		((Parent_Folder) (Undo) (Find, "Alt+F")(Select, "Alt+S")(Cancel, "Alt+C"))

	controls()
		{
		return Object('Vert'
			Object('Horz'
				Object('Field' status: TranslateLanguage('current directory')
					width: 41 readonly: )
				#('Skip' 3)
				#('Toolbar' Parent_Folder Undo)
				name : 'Hzt')
			#('Skip' 3)
			Object('Horz'
				#('ListBox' sort: true, xmin: 200, ymin: 200 ystretch: 0),
				'Skip'
				Object('Vert'
					#('Pane' ('Image', xmin: 200, ymin: 200), ystretch: 0)
					Object('ChooseList',
						width: 20, status: TranslateLanguage('selected image type')))
				name : 'Hzl')
			#('Skip')
			Object('Horz'
				#('Field' width: 12)
				#(Skip 3)
				#('Button' '&Find' xmin: 50)
				'Fill'
				Object('CheckBox' TranslateLanguage('show image'))
				'Skip'
				#('Button' '&Select' xmin: 50)
				'Skip'
				#('Button' '&Cancel' xmin: 50)
				name : 'Hzb')
			#('Skip')
			#('Status')
			xstretch: 1
			ystretch: 1)
		}

	imageslist(filename = "")
		{
		.Vert.Hzt.Field.Set(.currentdir)
		.lbox.DeleteAll()
		.setImage('')
		.Vert.Status.Set('')
		for file in Dir(.currentdir $ '*.*').Sort!()
			{
			// "/" is a directory not a file
			if file[-1] isnt "/"
				{
				filext = .getFileExt(file)
				if .imagetype is "(*.*)"
					{
					//add all the images
					if .filterlist.Find(filext) isnt false
						.lboxAddItem(file)
					}
				else
					if .imagetype is filext
						.lboxAddItem(file)
				}
			}
		.Vert.Status.Set(.lbox.GetCount() $ TranslateLanguage(' images listed'))
		if .lbox.GetCount() is 0
			enable = false
		else
			{
			enable = true
			selection = (filename isnt "")
				? .lbox.FindString(.filename)
				: 0
			.lbox.SetFocus()
			.lbox.SetCurSel(selection)
			.lbox.SetTopIndex(selection)
			.ListBoxSelect(selection)
			}
		.Vert.Hzb.Select.SetEnabled(enable)
		}
	lboxAddItem(file)
		{
		try
			.lbox.AddItem(file)
		catch
			.Vert.Status.Set(TranslateLanguage('error adding files names'))
		}
	getFileExt(filename)
		{
		ext = filename.AfterLast('.').Lower()
		return (ext is "")
			? ""
			: ext
		}
	setImage(value)
		{
		if .Vert.Hzb.CheckBox.Get()
			.Vert.Hzl.Vert.Pane.Image.Set(value)
		else
			.Vert.Hzl.Vert.Pane.Image.Set("")
		}
	get()
		{ return .currentdir $ .filename	}
	set(file)
		{
		if file is false
			file = ""  // open in computer resources
		else if file is true
			file = GetCurrentDirectory() // open in current directory
		.filename = Paths.Basename(file)
		.currentdir = file.BeforeLast(.filename)
		if .currentdir is ""
			.currentdir = Paths.EnsureTrailingSlash(GetCurrentDirectory())
		}
	On_Parent_Folder()
		{
		selecteddir = BrowseFolderName(hwnd: .Window.Hwnd, title: .opendirmsg,
			initialPath: .currentdir)
		if selecteddir is ""
			return
		.currentdir = Paths.EnsureTrailingSlash(selecteddir)
		.imageslist()
		//.Send("NewValue", .get())
		}

	On_Undo()
		{
		.set(.initialfile)
		.imageslist(.filename)
		}

	On_Select()
		{ .Window.Result(.get()) }

	On_Find()
		{
		sel = .lbox.FindString(.Vert.Hzb.Field.Get())
		if sel isnt -1
			{
			.ListBoxSelect(sel)
			.lbox.SetFocus()
			.lbox.SetCurSel(sel)
			.lbox.SetTopIndex(sel)
			}
		}

	ListBoxDoubleClick(i)
		{
		.ListBoxSelect(i)
		.On_Select()
		}

	ListBoxSelect(i)
		{
		.filename = .lbox.GetText(i)
		.setImage((.filename isnt "") ? .currentdir $ .filename : "")
		}
	NewValue(value, source)
		{
		if source.Name is 'ChooseList'
			{
			.imagetype = (value.Suffix?('(*.*)')) ? '(*.*)' : value
			.imageslist()
			}
		if source.Name is 'CheckBox'
			{
			.setImage(.currentdir $ .filename)
			.lbox.SetFocus()
			}
		}
	}