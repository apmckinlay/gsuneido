// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(model, .enableMultiSelect = false)
		{
		.selection = Object()
		.model = model
		}

	NotEmpty?()
		{
		return not .selection.Empty?()
		}

	GetSelectedRecords()
		{
		return .selection.Copy()
		}

	HasSelectedRow?(rec)
		{
		if rec is false
			return false
		return .selection.Has?(rec)
		}

	ClearSelect(rec = false)
		{
		if rec isnt false
			.selection.Remove(rec)
		else
			.selection = Object()
		}

	ReloadRecord(oldrec, newrec)
		{
		if .selection.Has?(oldrec)
			{
			.addRecordToSelection(newrec)
			.selection.Remove(oldrec)
			}
		}

	addRecordToSelection(rec)
		{
		if rec isnt false
			.selection.Add(rec)
		}

	UpdateShiftStart(absRowIndex, rows)
		{
		if .model.Offset < 0
			{
			if .shiftStart isnt false and .shiftStart <= absRowIndex + rows
				.shiftStart -= rows
			}
		else
			{
			if .shiftStart isnt false and .shiftStart > absRowIndex
				.shiftStart += rows
			}
		}

	ClearShiftStart()
		{
		.shiftStart = false
		}

	DecreaseShiftStart(rows)
		{
		if .model.Offset < 0 and .shiftStart isnt false
			.shiftStart -= rows
		}

	shiftStart: false
	SelectRows(ctrl, shift, row)
		{
		// ignore ctrl + shift at the same time which doesn't really make sense
		// PLUS we don't have row status for focused but not selected which would be
		// needed in order to mimic the behavior of file explorer in this case
		if shift is true and ctrl is true
			return
		if .shiftStart is false
			.shiftStart = row
		.handleShift(shift, ctrl, row)
		.handleCtrl(ctrl, shift, row)
		}

	handleShift(shift, ctrl, row)
		{
		if shift is true and .enableMultiSelect is true
			{
			oldKeys = .selection
			if ctrl isnt true
				.selection = Object()

			.selectGroupStartingAt(row, oldKeys)
			}
		if shift is false
			.shiftStart = row
		}

	handleCtrl(ctrl, shift, row)
		{
		rec = .model.GetRecord(row - .model.Offset)
		if ctrl is true and shift is false and .enableMultiSelect is true
			{
			if .selection.Has?(rec)
				.selection.Remove(rec)
			else
				.addRecordToSelection(rec)
			}
		if ((ctrl isnt true and shift isnt true) or .enableMultiSelect is false)
			{
			.selection = Object()
			.addRecordToSelection(rec)
			}
		}

	SelectOneRow(row)
		{
		if false is rec = .model.GetRecord(row)
			return
		.selection = Object()
		.addRecordToSelection(rec)
		}

	selectGroupStartingAt(row, oldKeys)
		{
		i = Min(.shiftStart, row)
		to = i + Abs(.shiftStart - row)
		limit = .model.Limit()
		msg = 'Cannot select more than ' $ limit $ ' records at a time.'
		if limit < Abs(.shiftStart - row)
			{
			.selection = oldKeys
			throw msg // will be caught
			}
		else
			for (; i <= to; ++i)
				{
				if false is rec = .model.GetRecord(i - .model.Offset)
					{
					.selection = oldKeys
					throw 'VirtualList tried to select non-existing record.'
					}
				if rec.vl_expand? isnt true
					.addRecordToSelection(rec)
				}
		}

	PageKey(focusedRow, shift, selectRowFn, up? = false)
		{
		listBottom = .listBottom()
		focusedRowPos = (up? ? focusedRow - .model.Offset : listBottom - focusedRow)
		pageSize = (up? ? -1 : 1) * (.model.VisibleRows - 1)

		if .focusOffScreen?(focusedRowPos, focusedRow, listBottom)
			{
			.pageReposition(up?, focusedRow, pageSize, selectRowFn, shift)
			return
			}

		// when focusedRow is in the middle of screen
		rec = .model.GetRecord(focusedRow - .model.Offset)
		if rec.vl_expanded_rows isnt ''	and
			rec.vl_expanded_rows + focusedRow >= listBottom
			{
			.pageReposition(up?, focusedRow, pageSize, selectRowFn, shift)
			return
			}

		if up?
			selectRowFn(focusedRow - focusedRowPos, :shift)
		else // down
			selectRowFn(listBottom, :shift)
		}

	focusOffScreen?(focusedRowPos, focusedRow, listBottom)
		{
		return focusedRowPos is 0 or focusedRow < .model.Offset or focusedRow > listBottom
		}

	listBottom()
		{
		return .model.End?()
			? .model.GetLastVisibleRowIndex()
			: .model.Offset + .model.VisibleRows - 1
		}

	pageReposition(up?, focusedRow, pageSize, selectRowFn, shift)
		{
		if false is selectRowFn(focusedRow + pageSize, :shift)
			{ // when focusedRow is on the first page or last page
			if focusedRow >= 0
				selectRowFn(up? ? 0 : .model.GetLastVisibleRowIndex(), :shift)
			else // start from last
				selectRowFn(up? ? .model.Offset : -1, :shift)
			}
		}

	AdjustFocusedRow(focusedRow, rowNum)
		{
		if focusedRow is false
			return false

		focusedRow = .model.ValidateRow(.focusRow(focusedRow, rowNum), returnBoundary?:)
		if focusedRow is false
			return false

		while false isnt focusedRec = .model.GetRecord(focusedRow - .model.Offset)
			if focusedRec.Member?('vl_rows')
				focusedRow--
			else
				break
		return focusedRec is false
			? false
			: focusedRow
		}

	focusRow(focusedRow, rowNum)
		{
		rowIdx = rowNum + .model.Offset
		if focusedRow < 0
			{
			if focusedRow <= rowIdx
				focusedRow++
			}
		else
			{
			if focusedRow > rowIdx
				focusedRow--
			}
		return focusedRow
		}
	}