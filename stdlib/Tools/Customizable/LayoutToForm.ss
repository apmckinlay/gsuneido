// Copyright (C) 2005 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(layout, sf, onlyCustomFields? = false)
		{
		if Object?(layout)
			{
			return layout
			}
		instance = new this
		return instance.Convert(layout, sf, :onlyCustomFields?)
		}
	Convert(layout, sf, .onlyCustomFields? = false)
		{
		layout = layout.Detab()
		.sf = sf
		.col = 0
		.othertextstr = ''
		.form = Object('Form')
		.groupsInfo = Object()
		.fields = Object()
		sf.ScanFormula(layout, .handle_field, .handle_other)
		.addOtherTextToForm()
		.groupsInfo.Sort!({|x,y| x.group < y.group })
		.groupsInfo.RemoveIf({ it.count < 2 })
		.adjustGroups(.form, .groupsInfo)
		return .form
		}
	handle_field(field)
		{
		// only add fields once, don't want duplicate fields in the form
		if .fields.Has?(field)
			return
		// do not convert text to real fields, just keep the text (ex. Bill To)
		if .onlyCustomFields? and not Customizable.CustomField?(field)
			{
			.othertextstr $= .sf.FieldToPrompt(field)
			return
			}
		.addOtherTextToForm()
		fieldElement = Object(field, group: .col)
		.form.Add(fieldElement)
		.addGroupsInfo(fieldElement, .groupsInfo)
		.fields.Add(field)
		.col += .sf.FieldToPrompt(field).Size()
		}
	maxHeight: 4
	EnsureEditorHeightLimit(field, fieldElement)
		{
		dd = Datadict(field)
		ctrl = dd.Control
		ctrlName = ctrl[0].RemoveSuffix('Control')
		if ctrlName is 'ScintillaAddonsEditor'
			fieldElement.height = Min(.maxHeight,
				ctrl.GetDefault('height', ScintillaAddonsEditorControl.Height))
		else if ctrlName is 'ScintillaRichWordAddons'
			fieldElement.height =
				Min(.maxHeight, ctrl.GetDefault('height', ScintillaAddonsControl.Height))
		else if ctrlName is 'Editor'
			fieldElement.height =
				Min(.maxHeight, ctrl.GetDefault('height', EditorControl.DefaultHeight))
		if fieldElement.Member?('height')
			fieldElement.ystretch = 0
		}
	handle_other(other)
		{
		if other.Has?('\n')
			{
			for i in .. other.Tr('^\n').Size()
				{
				.addOtherTextToForm()
				.form.Add('nl')
				}
			.col = 0
			other = other.AfterLast('\n')
			}
		if other isnt ""
			{
			.othertextstr $= other
			.col += other.Size()
			}
		}
	addOtherTextToForm()
		{
		if .othertextstr isnt ''
			.form.Add(Object('Static', .othertextstr))
		.othertextstr = ''
		}

	addGroupsInfo(item, groups)
		{
		if .hasGroup?(item)
			{
			if false is idx = .matchingGroup(groups, item)
				groups.Add([group: item.group, count: 1])
			else
				groups[idx].count++
			}
		}

	hasGroup?(item)
		{
		return Object?(item) and item.Member?('group')
		}

	matchingGroup(groups, item)
		{
		return groups.FindIf({ it.group is item.group })
		}

	adjustGroups(form, groups)
		{
		for item in form
			if .hasGroup?(item)
				{
				if false is idx = .matchingGroup(groups, item)
					item.Delete('group')
				else
					item.group = idx
				}
		}

	Revert(form, sf)
		{
		result = .preProcessGroups(form, sf)
		groupOb = result[0]
		lines = Object()
		(result.lineNum + 1).Times({ lines.Add('') })
		occupiedCols = Object()
		for idx in groupOb.Members().Sort!()
			{
			if idx % 2 is 0
				.handleNoGroup(groupOb[idx], lines, occupiedCols)
			else
				.handleGroup(groupOb[idx], lines, occupiedCols)
			}
		return lines.Join('\r\n')
		}

	handleNoGroup(ob, lines, occupiedCols)
		{
		for line in ob.Members().Sort!()
			{
			for item in ob[line]
				if String?(item)
					lines[line] $= item
				else
					{
					start = lines[line].Size()
					while (occupiedCols.Member?(start))
						start++
					lines[line] = lines[line].RightFill(start) $ item.prompt
					occupiedCols[start] = true
					}
			}
		}

	handleGroup(ob, lines, occupiedCols)
		{
		start = 0
		for line in ob.Members().Sort!()
			{
			if lines[line].Size() is 0
				continue
			start = Max(start, lines[line].Size())
			}
		while (occupiedCols.Member?(start))
			start++
		for line in ob.Members().Sort!()
			lines[line] = lines[line].RightFill(start) $ ob[line][0].prompt
		occupiedCols[start] = true
		}

	tabLen: 4
	preProcessGroups(form, sf)
		{
		fills = Object().Set_default(Object().Set_default(Object()))
		lineNum = 0
		idx = 0
		firstItemProcessed? = false
		startIdx = 0
		for fieldSpec in form
			{
			if fieldSpec is 'Form'
				continue
			if fieldSpec is 'nl'
				{
				lineNum++
				idx = startIdx
				continue
				}
			if not fieldSpec.Member?('group')
				{
				if idx % 2 is 1
					idx++
				if fieldSpec[0] is 'Static'
					fills[idx][lineNum].Add(fieldSpec[1])
				else
					fills[idx][lineNum].Add(.getField(fieldSpec, sf))
				}
			else
				{
				idx = fieldSpec.group * 2 + 1
				fills[idx][lineNum].Add(.getField(fieldSpec, sf))
				if not firstItemProcessed?
					startIdx = 2
				}
			firstItemProcessed? = true
			}
		return [fills, :lineNum]
		}

	getField(fieldSpec, sf)
		{
		field = fieldSpec[0]
		if false is prompt = sf.FieldToPrompt(field)
			if field is prompt = SelectFields.GetFieldPrompt(field)
				prompt = '???'
		return [:field, :prompt]
		}
	}
