// Copyright (C) 2010 Suneido Software Corp. All rights reserved worldwide.
RepeatControl
	{
	PerRow: 5
	New(type = '', reorderOnly = false, hideLabel? = false)
		{
		super(.layout(type, reorderOnly, hideLabel?), no_minus:, noPlus: reorderOnly)
		}
	layout(type, reorderOnly, hideLabel?)
		{
		imgCtrl = reorderOnly ? 'OpenImage' : 'OpenImageWithLabels'
		return Object('Horz'
			Object('Pane' Object(imgCtrl, :hideLabel?, :type, :reorderOnly,
				name: 'attachment0')),
			Object('Pane' Object(imgCtrl, :hideLabel?, :type, :reorderOnly,
				name: 'attachment1'))
			Object('Pane' Object(imgCtrl, :hideLabel?, :type, :reorderOnly,
				name: 'attachment2'))
			Object('Pane' Object(imgCtrl, :hideLabel?, :type, :reorderOnly,
				name: 'attachment3'))
			Object('Pane' Object(imgCtrl, :hideLabel?, :type, :reorderOnly,
				name: 'attachment4'))
			)
		}

	GetRowHeight(reorderOnly = false, hideLabel? = false)
		{
		height = reorderOnly
			? OpenImageControl.Xmin
			: OpenImageWithLabelsControl.ImageHeight +
				(hideLabel?
					? 0 :
					OpenImageWithLabelsControl.TextHeight + 2/*=EtchedLine*/)
		return height + (1 + GetSystemMetrics(SM.CYEDGE)) * 2/*from PaneControl*/
		}

	ImageDropFileList(files, source)
		{
		addLabels = false
		data = .getRowsData()
		dest = .getIndex(source)
		for file in files
			{
			src = data.Size() * .PerRow
			if false is (fileProc = OpenImageSelect.ResultFile(file,
				useSubFolder: OpenImageWithLabelsControl.UseSubFolder))
				continue // don't advance destination slot if copy failed
			if "" isnt OpenImageWithLabelsControl.SplitFile(.get(data, dest))
				.shiftAttachments(src, dest, data)
			else
				addLabels = true
			if addLabels
				fileProc = source.ProcessValue(fileProc)
			.put(data, dest, fileProc)
			currentSource = .sourceFromIndex(dest)
			currentSource.QueueDeleteAttachment(currentSource.FullPath(), '')
			dest++
			}
		return true
		}

	getRowsData()
		{
		return .GetRows().Map(#Get)
		}

	get(data, i)
		{
		row = (i / .PerRow).Int()
		if row >= data.Size()
			return ''
		col = i % .PerRow
		return data[row]['attachment' $ col]
		}
	put(data, i, value)
		{
		row = (i / .PerRow).Int()
		if row >= data.Size()
			{
			data.Add(rec = [])
			.AppendRow(rec)
			}
		col = i % .PerRow
		data[row]['attachment' $ col] = value
		}

	dragging: false
	switchedImageControls: false
	src: false
	ImageStartDrag(source)
		{
		.src = .getIndex(source)
		}

	// from suneido.js ImageComponent
	ImageEndDrag()
		{
		.src = false
		}

	ImageFinishDrag(source)
		{
		.dragging = false
		dest = .getIndex(source)
		if .src is false or dest is .src
			return
		data = .getRowsData()
		file = .get(data, .src)
		.put(data, .src, "")
		if "" isnt OpenImageWithLabelsControl.SplitFile(.get(data, dest))
			.shiftAttachments(.src, dest, data)
		else if "" is OpenImageWithLabelsControl.SplitLabel(file)["labels"]
			file = source.ProcessValue(file)
		.put(data, dest, file)
		}

	shiftAttachments(src, dest, data)
		{
		// Shifts up to the first empty slot between src and dest starting from
		// dest, if no empty slots then shifts up to src
		// Treats empty spaces with labels as empty
		emptySlot = dest
		if dest < src
			{
			emptySlot = .firstEmptySlot(data, dest)
			for(i = emptySlot; i > dest; i--)
				.put(data, i, .get(data, i - 1))
			}
		else
			{
			emptySlot = .firstEmptySlot(data, dest, true)
			for(i = emptySlot; i < dest; i++)
				.put(data, i, .get(data, i + 1))
			}
		}

	firstEmptySlot(data, start = 0, descending = false)
		{
		emptySlot = start
		increment = descending is false ? 1 : -1
		while not OpenImageWithLabelsControl.SplitFile(.get(data, emptySlot)).Blank?()
			{
			emptySlot = emptySlot + increment
			Assert(emptySlot >= 0, 'firstEmptySlot result less than 0')
			}
		return emptySlot
		}

	getIndex(source)
		{
		cur = source.Parent
		horzRow = cur.Parent
		col = horzRow.GetChildren().Find(cur)
		repeatRow = horzRow.Parent.Parent
		row = .GetRows().Find(repeatRow)
		dest = row * .PerRow + col
		return dest
		}

	ImageMouseLeave(.dragging)
		{
		.switchedImageControls = true
		}

	ImageMouseMove(source)
		{
		if .switchedImageControls
			{
			source.FindControl('image').Dragging = .dragging
			.switchedImageControls = false
			}
		}

	ImageDropFile(file, source)
		{
		data = .getRowsData()
		dest = .getIndex(source)
		src = data.Size() * .PerRow
		if "" isnt OpenImageWithLabelsControl.SplitFile(.get(data, dest))
			.shiftAttachments(src, dest, data)
		else
			file = source.ProcessValue(file) // adds existing labels

		.put(data, dest, file)
		source.QueueDeleteAttachment(source.FullPath(), '')
		return true
		}

	ScanAttachment()
		{
		if .GetReadOnly()
			{
			.AlertInfo("Scan Attachment", "This field is currently protected")
			return
			}
		data = .getRowsData()
		dest = .firstEmptySlot(data)
		.sourceFromIndex(dest).On_Context_Scan_Attachment()
		return
		}

	sourceFromIndex(index)
		{
		row = (index / .PerRow).RoundDown(0)
		if row >= (rows = .GetRows()).Size()
			{
			data = .getRowsData()
			.put(data, index, '') // Construct an extra row
			}
		col = index - row * 5
		sourceRow = rows[row]
		source = sourceRow.FindControl("attachment" $ col)
		return source
		}

	ScanAs()
		{
		if .GetReadOnly()
			{
			.AlertInfo("Scan Attachment", "This field is currently protected")
			return
			}
		data = .getRowsData()
		dest = .firstEmptySlot(data)
		.sourceFromIndex(dest).On_Context_Scan_Attachment_As()
		return
		}

	ScanSelectSource()
		{
		Scanning().SelectScanner(.Window.Hwnd)
		}

	AttachmentFieldName()
		{
		return .Name
		}
	}
