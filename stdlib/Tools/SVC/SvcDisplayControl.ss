// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
/*
SvcDisplayControl handles the display of code in the bottom half of the window based on
the selection. The current selection is passed in to Display(), and this class reasons
what code should be displayed in the ScintillaControl(s) based on the type of change, date
of the change, which window it is in, and if there are more changes to the record or not.

SvcModel is accessed to get server information about records.

DisplayChange() is called from Display() to display a change that isn't an image or an
add. It calls getDiffData(), which then calls the appropriate display function to get the
required records to show, as well the title and comment, based on the specifics of the
selection.
*/
PassthruController
	{
	Name: "svc_display"
	Controls: #(Vert (Pane Fill ystretch: 10) name: 'bottom')
	New()
		{
		.bottom = .FindControl("bottom")
		.Redir('On_Use_Master')
		.Redir('On_Merge')
		}
	SetModel(.model) { }
	RemoveAll()
		{
		.bottom.RemoveAll()
		.removeInMemory()
		}
	toRemove: false
	removeInMemory()
		{
		if .toRemove is false
			return
		InMemory.Remove(.toRemove)
		.toRemove = false
		}
	Reset()
		{
		.RemoveAll()
		.bottom.Insert(0, #(Pane Fill ystretch: 1 ymin: 400))
		}

	// the type is treated as master action
	// 		- local add is passed as type "-" (delete)
	// 		- local delete is passed as type "+" (add)
	Display(.name, .table, .type, showComment = false, masterNewer? = false,
		lib_committed = false)
		{
		.RemoveAll()
		curFocus = GetFocus()
		if name =~ '/res/.*(?i)[.](png|jpg|jpeg|gif|svg)$'
			.displayImage()
		else if name =~ '/res/.*(?i)[.](emf|wmf|ico)$'
			.displayMetaImage()
		else if type is '+'
			.displayNewRecord(lib_committed, showComment)
		else if type is '-'
			.displayDeleted(lib_committed, showComment)
		else
			.displayChange(showComment, :masterNewer?, :lib_committed)
		SetFocus(curFocus)
		}

	displayImage()
		{
		localStyle = .style('suneido:/' $ .table $ .name)

		master = .model.GetMasterRec(.table, .name)
		if .type is '-' or master is false
			{
			.bottom.Insert(0, Object('Mshtml' style: localStyle))
			return
			}

		// keep the extension for the file type
		url = .toRemove = InMemory.Add(master.text) $ '.' $ .name.AfterLast('.')
		masterStyle = .style(url)
		if .type is '+'
			.bottom.Insert(0, Object('Mshtml' style: masterStyle))
		else
			.bottom.Insert(0,
				Object('DiffImage' '', '', style1: localStyle, style2: masterStyle))
		}
	style(url)
		{
		return 'body {
			background-color: #eee;
			background-image: url(' $ url.Replace(' ', '%20') $ ');
			background-size: contain;
			background-repeat: no-repeat;
			}'
		}

	displayMetaImage()
		{
		img_name = .name.AfterLast('res/')
		master = .model.GetMasterRec(.table, .name)
		if .type is '-' or master is false
			.bottom.Insert(0, Object('Scroll'
				Object('Image', .table $ '%' $ img_name) noEdge:))
		else if .type is '+'
			.bottom.Insert(0, Object('Scroll', Object('Image', master.text) noEdge:))
		else
			.bottom.Insert(0, Object('DiffImage', .table $ '%' $ img_name, master.text))
		}

	displayNewRecord(lib_committed, showComment)
		{
		masterRec = .model.GetMasterRec(.table, .name, lib_committed)
		preDef = .model.GetPrevDefinition(.name, .table)
		if preDef isnt false
			.displayChange(showComment?:,
				customDiff: [
					newRec: preDef,
					newTitle: 'LOCAL ' $ preDef.table $ ' DEFINITION <Npath>',
					newTitleRight: '<Nmdate>  <Ndate>',
					oldRec: masterRec,
					oldTitle: 'MASTER <Opath>', oldTitleRight: '<Odate>',
					commentRec: masterRec])
		else
			.displayOne(masterRec, showComment)
		}

	displayDeleted(lib_committed, showComment)
		{
		x = .model.GetMasterRec(.table, .name, lib_committed, delete?:)
		old = lib_committed isnt false and
			(prevRow = .model.GetPrevMasterRec(.table, .name, lib_committed)) isnt false
			? .model.GetMasterRec(.table, .name, prevRow.lib_committed)
			: .model.GetLocalRec(.table, .name)
		if x is false
			x = old
		else
			x.text = old.text
		if false is preDef = .model.GetPrevDefinition(.name, .table)
			.displayOne(x, showComment, deleted:)
		else
			.displayChange(showComment,
				customDiff: [oldRec: preDef,
					oldTitle: preDef.table.Upper() $ " <Opath>", oldTitleRight: "<Odate>",
					commentRec: x])
		}

	displayOne(x, showComment?, deleted = false)
		{
		if x is false
			return
		vertOb = Object('Vert', #(DisplayCode, ymin: 400, ystretch: 1))
		if false isnt rec = .model.GetLocalRec(.table, .name, :deleted)
			vertOb.Add(Object('StaticText', text: 'LOCAL ' $ rec.path $ '/' $ .name,
				size: '+2'), at: 1)
		.bottom.Insert(0, vertOb)
		.bottom.FindControl('Editor').Set(x.text)
		if showComment?
			{
			desc = x.lib_committed.ShortDateTime() $ " " $ x.id $ " - " $ x.comment
			.bottom.Insert(0, Diff2Control.Comment(desc))
			}
		}

	flip: false
	flipButtons(type)
		{
		return 	Object('HorzEqual'
			Object('Button', 'Show ' $ type),
			#Fill
			#(Button 'Merge')
			#Skip
			#(Button 'Use Master')
			#Fill
			pad: 0
			)
		}
	displayChange(showComment?, masterNewer? = false, lib_committed = false,
		customDiff = #())
		{
		merge = .type is "#" or .type is "%"
		data = .getDiffData(merge, masterNewer?, lib_committed, customDiff)

		control = Object()
		if data.oldRec is false or data.newRec is false
			{
			text = Object?(data.oldRec) ? data.oldRec.text : data.newRec.text
			control = Object('DisplayCode', ymin: 400, ystretch: 1, set: text)
			}
		else if .type is "%"
			control = .displayConflictChange(data, showComment?, masterNewer?)
		else
			control = .displayRegularChange(data, showComment?, masterNewer?)

		.bottom.Insert(0, control)
		if false isnt .flip = .FindControl('Flip')
			{
			mode = UserSettings.Get('VersionControl-OrigMerge', def: 0)
			.flip.SetCurrent(mode is false ? 0 : mode)
			}
		}

	displayConflictChange(data, showComment?, masterNewer?)
		{
		newTitle = .makeTitle(data.newTitle, data.oldRec, data.newRec, data.base)
		oldTitle = .makeTitle(data.oldTitle, data.oldRec, data.newRec, data.base)
		comment = .makeComment(data.commentRec, showComment?)
		base = data.base is false ? false : data.base.text
		return Object('Vert',
			Object('Flip'
				Object('Vert', .flipButtons('Original')
				Object("MergeDiff", data.newRec.text, data.oldRec.text, .table, .name,
					newTitle, oldTitle,
					:base, :comment, tags: "ML", gotoButton:, newOnRight?: masterNewer?)),
				Object('Vert', .flipButtons('Merged')
				Object("Diff2", data.newRec.text, data.oldRec.text, .table, .name,
					"Merged", "Original", base, :comment, tags: "ML",
					gotoButton:, newOnRight?: masterNewer?))))
		}

	displayRegularChange(data, showComment?, masterNewer?)
		{
		newTitle = .makeTitle(data.newTitle, data.oldRec, data.newRec, data.base)
		oldTitle = .makeTitle(data.oldTitle, data.oldRec, data.newRec, data.base)
		titleNewRight = titleOldRight = ''
		if data.Member?('newTitleRight')
			titleNewRight = .makeRightTitle(data.newTitleRight, data.oldRec, data.newRec)
		if data.Member?('oldTitleRight')
			titleOldRight = .makeRightTitle(data.oldTitleRight, data.oldRec, data.newRec)
		comment = .makeComment(data.commentRec, showComment?)
		base = data.base is false ? false : data.base.text
		return Object("Diff2", data.newRec.text, data.oldRec.text, .table, .name,
			newTitle, oldTitle, base, :comment, tags: "ML",
			gotoButton:, newOnRight?: masterNewer?, :titleNewRight, :titleOldRight)
		}

	ShowHideLineDiff(value, source)
		{
		if .flip is false
			return

		// prevent function call overflow
		if source isnt merge = .FindControl('MergeDiff')
			merge.ToggleShowLineDiff(value)
		if source isnt diff = .FindControl('Diff')
			diff.ToggleShowLineDiff(value)
		}

	getDiffData(conflict?, masterNewer?, lib_committed, customDiff)
		{
		return masterNewer?
			? conflict?
				? lib_committed isnt false
					? .displayMergedConflict(lib_committed)
					: .displayConflict()
				: .displayMasterChange(lib_committed)
			: .displayLocalChange(customDiff)
		}

	displayLocalChange(customDiff)
		{
		diff = customDiff.Copy()
		diff.GetInit('oldRec', 		{ .model.GetMasterRec(.table, .name)	})
		diff.GetInit('newRec', 		{ .model.GetLocalRec(.table, .name) 	})
		diff.GetInit('oldTitle', 	{ "MASTER <Opath>" })
		diff.GetInit('newTitle', 	{ "LOCAL <Npath>" })
		diff.GetInit('oldTitleRight', { "<Odate>" })
		diff.GetInit('newTitleRight', { "<Nmdate>  <Ndate>" })
		diff.GetInit('commentRec', 	{ false })
		diff.Set_default(false)
		return diff
		}

	displayConflict()
		{
		oldRec = .model.GetMasterRec(.table, .name)
		if false is newRec = .model.GetLocalRec(.table, .name)
			newRec = .model.GetLocalRec(.table, .name, deleted:)
		// Makes sure to handle deletion or unsent record on local
		base = newRec isnt false and newRec.lib_committed is ""
			? []
			: .model.GetMasterFromLocal(.table, .name)

		newTitle = "LOCAL <Npath>"
		oldTitle = "MASTER <Opath>"
		newTitleRight = '<Nmdate>  <Ndate>'
		oldTitleRight = '<Odate>'
		commentRec = oldRec
		return Object(:newRec, :oldRec, :newTitle, :oldTitle, :base, :commentRec,
			:newTitleRight, :oldTitleRight).
			Set_default(false)
		}

	displayMasterChange(lib_committed)
		{
		prevRec = .model.GetPrevMasterRec(.table, .name, lib_committed)
		local = .model.GetLocalRec(.table, .name)
		oldRec = prevRec is false
			? lib_committed is false and local.lib_modified is ""
				? local
				: .model.GetMasterFromLocal(.table, .name)
			: prevRec
		newRec = .model.GetMasterRec(.table, .name, lib_committed)
		newTitle = prevRec is false
			? "MASTER <Npath>"
			: "AS OF <Ndate> <Npath>"
		oldTitle = prevRec is false
			? "LOCAL <Opath>"
			: "PREVIOUS <Odate> <Opath>"
		newTitleRight = prevRec is false ? '<Ndate>' : ''
		oldTitleRight = prevRec is false ? '<Omdate>  <Odate>' : ''
		commentRec = newRec
		return Object(:newRec, :oldRec, :newTitle, :oldTitle, :commentRec,
			:newTitleRight, :oldTitleRight).Set_default(false)
		}

	displayMergedConflict(lib_committed)
		{
		prevRec = .model.GetPrevMasterRec(.table, .name, lib_committed)
		Assert(prevRec isnt: false)
		oldRec = prevRec
		merged = .model.GetMergedRec(.table, .name, .model.GetLocalRec(.table, .name),
			.model.GetMasterRec(.table, .name))
		newRec = [text: merged.merged, :lib_committed]
		oldTitle = "PREVIOUS <Odate> <Opath>"
		newTitle = "MERGED"
		return Object(:newRec, :oldRec, :newTitle, :oldTitle).Set_default(false)
		}

	makeTitle(titleTemplate, oldRec, newRec, base)
		{
		if base is false
			base = []
		return titleTemplate.
			Replace('\<<Npath>\>', .buildPath(newRec)).
			Replace('\<<Opath>\>', .buildPath(oldRec)).
			Replace('\<<Bpath>\>', .buildPath(base)).
			Replace('\<<Ndate>\>', .replacementDate(newRec, 'lib_committed')).
			Replace('\<<Odate>\>', .replacementDate(oldRec, 'lib_committed'))
		}

	replacementDate(rec, dateMember, prefix = 'Com')
		{
		return rec.Member?(dateMember)
			? prefix $ ': ' $ rec[dateMember].ShortDateTime()
			: ''
		}

	buildPath(rec)
		{
		if false is .model.Library?(.table) and rec.name.Prefix?(rec.path)
			return rec.name
		return rec.path $ '/' $ rec.name
		}

	makeRightTitle(titleRightTemplate, oldRec, newRec)
		{
		return titleRightTemplate.
			Replace('\<<Nmdate>\>', .replacementDate(newRec, 'lib_modified', 'Mod')).
			Replace('\<<Ndate>\>', .replacementDate(newRec, 'lib_committed')).
			Replace('\<<Omdate>\>',  .replacementDate(oldRec, 'lib_modified', 'Mod')).
			Replace('\<<Odate>\>',  .replacementDate(oldRec, 'lib_committed'))
		}

	makeComment(rec, showComment?)
		{
		return showComment? and rec isnt false
			? rec.lib_committed.ShortDateTime() $ " " $ rec.id $ " - " $ rec.comment
			: false
		}

	GetGoToLine()
		{
		return .bottom.GetChildren()[0].Base?(Diff2Control)
			? .bottom.Diff.CurLine - .bottom.Diff.LineOffset()
			: 0
		}

	On_Show_Merged()
		{
		.onFlip()
		}
	On_Show_Original()
		{
		.onFlip()
		}

	onFlip()
		{
		if .flip is false
			return
		.flip.Flip()
		UserSettings.Put('VersionControl-OrigMerge', .flip.GetCurrent())
		}

	Destroy()
		{
		if .flip isnt false
			UserSettings.Put('VersionControl-OrigMerge', .flip.GetCurrent())
		.removeInMemory()
		super.Destroy()
		}
	}
