// Copyright (C) 2019 Axon Development Corporation All rights reserved worldwide.
class
	{
	New(@args)
		{
		.noDragDrop = args.GetDefault('noDragDrop', false)
		.noHeaderButtons = args.GetDefault('noHeaderButtons', false)
		.headerSelectPrompt = args.GetDefault('headerSelectPrompt', false)
		if args.Member?(0) and Object?(args[0])
			args = args[0]
		.headCols = Object()
		for arg in args.Values(list:)
			.AddItem(arg)
		}

	AddItem(field, width = false, tip = false, sort = false)
		{
		.InsertItem(.GetItemCount(), field, width, tip, sort)
		}

	InsertItem(idx, field, width = false, tip = false, sort = false)
		{
		Assert(Integer?(idx) and 0 <= idx and idx <= .GetItemCount())
		text = .GetHeaderText(field)
		if width is false
			width = .getWidth(field, text)
		.headCols.Add([:text, :width, :field, :tip, :sort], at: idx)
		}

	GetHeaderText(field, _capFieldPrompt = false)
		{
		if .headerSelectPrompt is 'no_prompts' or field is 'listrow_deleted'
			return capFieldPrompt isnt false ? capFieldPrompt : field

		if .headerSelectPrompt is false or field.Prefix?('custom_')
			return Datadict.PromptOrHeading(field)

		return Datadict.SelectPrompt(field, excludeTags: #(Internal))
		}

	headerBorder: 2 // two character wdith
	GetDefaultColumnWidth(field)
		{
		return .getWidth(field, .GetHeaderText(field))
		}
	getWidth(field, text)
		{
		// FIXME: using magical number to convert the text size number to a width in px on browse
		heading_width = text.Size() * 7/*=width per char*/
		format_width = FieldFormatWidth(field, 7/*=width per char*/)
		return Max(heading_width, format_width) + .headerBorder
		}

	GetItem(idx)
		{
		return .headCols[idx]
		}

	GetItemCount()
		{
		return .headCols.Size()
		}

	GetItemWidth(idx)
		{
		return .headCols[idx].width
		}

	SetItemWidth(i, width)
		{
		.headCols[i].width = width
		}

	SetItemSort(i, sort)
		{
		.headCols[i].sort = sort
		}

	SetItemFormat(i, format)
		{
		.headCols[i].format = format
		}

	DeleteItem(idx)
		{
		Assert(Integer?(idx) and 0 <= idx and idx < .GetItemCount())
		.headCols.Delete(idx)
		}

	Reorder(col, newIdx)
		{
		temp = .headCols[col]
		.headCols.Delete(col).Add(temp, at: newIdx)
		}

	Clear()
		{
		.headCols.Delete(all:)
		}

	Get()
		{
		return .headCols.Copy().Filter({ it.width isnt false})
		}
	}