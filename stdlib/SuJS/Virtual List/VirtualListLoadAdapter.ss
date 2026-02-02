// Copyright (C) 2020 Axon Development Corporation All rights reserved worldwide.
class
	{
	loadedTop: false
	loadedBottom: false
	rowsPerLoad: 10
	New(.grid, .model, .convertFn, .saveAndCollapseRelease)
		{
		}

	InitLoad()
		{
		.loadedTop = false
		.loadedBottom = false
		dataBatch = Object()
		data = .model.GetLoadedData()
		members = data.Members().Sort!()
		for i in members
			dataBatch.Add((.convertFn)(data[i]), at: i)
		if members.NotEmpty?()
			{
			.loadedTop = members.First()
			.loadedBottom = members.Last()
			}

		.grid.Act(#AddBatch, dataBatch,
			newTop: .loadedTop, topEnded?: .model.Begin?(),
			newBottom: .loadedBottom, bottomEnded?: .model.End?(),
			loadOnTop?: .model.GetInitStartLast())
		}

	ForLoadedData(block)
		{
		if .loadedTop is false or .loadedBottom is false
			return
		data = .model.GetLoadedData()
		for (i = .loadedTop; i <= .loadedBottom; i++)
			block(i, data[i])
		}

	ModelIsReset?()
		{
		if .loadedTop is false or .loadedBottom is false
			return false
		data = .model.GetLoadedData()
		if data.Empty?()
			return true
		members = data.Members().Sort!()
		return .loadedTop isnt members.First() or .loadedBottom isnt members.Last()
		}

	LoadFrom(row)
		{
		dataBatch = Object()
		if row is .loadedTop
			{
			offset = -(.rowsPerLoad + .model.Offset - .loadedTop)
			.updateOffset(offset)
			.loadUp(row-1, dataBatch)
			.grid.Act(#AddBatch, dataBatch,
				newTop: .loadedTop, topEnded?: .model.Begin?()
				newBottom: false, bottomEnded?: false
				loadOnTop?:)
			}
		else if row is .loadedBottom
			{
			offset = .rowsPerLoad +
				Max(0, .loadedBottom - .model.Offset - .model.VisibleRows)
			.updateOffset(offset)
			.loadDown(row+1, dataBatch)
			.grid.Act(#AddBatch, dataBatch,
				newTop: false, topEnded?: false,
				newBottom: .loadedBottom, bottomEnded?: .model.End?()
				loadOnTop?: false)
			}
		else
			BookLog('LoadFrom Failed',
				params: [:row, top: .loadedTop, bottom: .loadedBottom])
		}

	loadUp(i, dataBatch)
		{
		data = .model.GetLoadedData()
		while data.Member?(i)
			{
			.loadedTop = i
			dataBatch.Add((.convertFn)(data[i]), at: i)
			i--
			}
		}

	loadDown(i, dataBatch)
		{
		data = .model.GetLoadedData()
		while data.Member?(i)
			{
			.loadedBottom = i
			dataBatch.Add((.convertFn)(data[i]), at: i)
			i++
			}
		}

	EnsureRow(row)
		{
		if row is VirtualListGridBodyComponent.FillRowY
			return

		if row >= .loadedTop and row <= .loadedBottom
			return

		offset = row < .loadedTop
			? row - .model.Offset
			: row - .model.Offset - .model.VisibleRows + 1

		if .updateOffset(offset) is 0
			return

		dataBatch = Object()
		newTop = false
		topEnded? = false
		newBottom = false
		bottomEnded? = false
		loadOnTop? = false
		if row < .loadedTop
			{
			.loadUp(.loadedTop - 1, dataBatch)
			topEnded? = .model.Begin?()
			newTop = .loadedTop
			loadOnTop? = true
			}
		else
			{
			.loadDown(.loadedBottom + 1, dataBatch)
			bottomEnded? = .model.End?()
			newBottom = .loadedBottom
			}

		if dataBatch.NotEmpty?()
			.grid.Act(#AddBatch, dataBatch,
				:newTop, :topEnded?, :newBottom, :bottomEnded?, :loadOnTop?)
		}

	InsertNewRecord(record = false, row_num = false, force = false)
		{
		shiftTop? = .model.GetStartLast()
		// need to manually set row_num (pos = 'end'),
		// otherwise VirtualListModel.InsertNewRecord uses .VisibleRows
		if row_num is false
			{
			bottom = .loadedBottom is false ? -1 : .loadedBottom
			row_num = bottom + 1 - .model.Offset
			}

		if false is newRec = .model.InsertNewRecord(record, row_num, :force)
			return false

		newRowNum = .model.GetRecordRowNum(newRec) + .model.Offset
		.updateTopBottom(shiftTop?)
		.grid.Act(#InsertData, newRowNum, (.convertFn)(newRec), shiftTop?)
		return Object(:newRec, :newRowNum)
		}

	updateTopBottom(shiftTop?)
		{
		if shiftTop? is true
			{
			if .loadedTop is false
				.loadedTop = .loadedBottom = -1
			else
				.loadedTop--
			}
		else
			{
			if .loadedBottom is false
				.loadedTop = .loadedBottom = 0
			else
				.loadedBottom++
			}
		}

	DeleteRecord(rowNum)
		{
		shiftTop? = .model.GetStartLast()
		if shiftTop? is true
			.loadedTop++
		else
			.loadedBottom--
		if .loadedTop > .loadedBottom // empty
			.loadedTop = .loadedBottom = false
		.grid.Act(#DeleteRecord, rowNum, shiftTop?)
		}

	releaseTo: false
	updateOffset(offset)
		{
		releaseTop? = offset > 0
		result = .model.UpdateOffset(offset, .release)
		if .releaseTo isnt false
			{
			.grid.Act(#Release, .releaseTo, releaseTop?)
			if releaseTop? is true
				.loadedTop = .releaseTo + 1
			else
				.loadedBottom = .releaseTo - 1
BookLog('Virtuallist recycle', params: [releaseTo: .releaseTo, loadedTop: .loadedTop,
	loadedBottom: .loadedBottom, modelOffset: .model.Offset, :offset])
			.releaseTo = false
			}
		return result
		}

	release(rec, row_num)
		{
		if ((.saveAndCollapseRelease)(rec, row_num) is false)
			{
			.releaseTo = false
			return false
			}

		.releaseTo = row_num + .model.Offset
		return true
		}
	}
