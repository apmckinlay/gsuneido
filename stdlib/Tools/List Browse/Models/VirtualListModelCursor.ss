// Copyright (C) 2014 Suneido Software Corp. All rights reserved worldwide.
/* Adding non-table rows to Virtual List

!!! NOTE: This should ONLY be done with readonly lists !!!

1. Declare VirtualList_ExtraSetupRecordFn in the control
2. In the callable it returns have:
	data = rec.Copy() 		// Copies the data for the additional rows
	rec.Delete(all:)		// Clears the base record, allowing us to re-use its reference
	vl_multi_row_count = 0 	// Tracks how many rows are being added
	...
	row = data.Copy()		// Copy data to ensure the row represents the original record
	<modify row as desired>
	rec.Add(row)			// Add row to be displayed via the Cursor
	vl_multi_row_count++ 	// Indicates to the Cursor how many lines / positions to move
	...
	// Once all rows are added, call the below. This is required during Release Up/Down
	rec.Each({ it.vl_multi_row_count = vl_multi_row_count })
	// Lastly add this to the root of rec to ensure Cursor detects the multi rows
	rec.vl_multi_row_count

You can find an example at: VirtualListModelCursor_Test.setupRecord
*/
class
	{
	cursor: 	false
	t:			false
	Seeking: 	false
	keys:		false

	New(.data, .query, .startLast = false, .setupRecord = false, .asof = false,
		useQuery = 'auto')
		{
		.setupCursor(:useQuery)
		.prevAction = not startLast ? 'Next' : 'Prev'
		.Pos = 0
		if startLast
			.cursor.Rewind()
		}

	setupCursor(useQuery)
		{
		if useQuery is true
			{
			.seekQuery()
			return
			}

		if useQuery is false
			{
			.seekCursor()
			return
			}

		try
			if .asof is false
				.seekCursor()
			else
				.seekQuery()
		catch(unused, 'invalid query')
			.seekQuery()
		}

	UsingCursor?()
		{
		return .cursor.Base?(SeekCursor)
		}

	seekCursor()
		{
		.cursor = SeekCursor(.query)
		if .highCostCursor?()
			throw 'invalid query'
		}

	seekQuery()
		{
		if .cursor isnt false
			.cursor.Close()
		if .t is false or .t.Ended?()
			.t = Transaction(read:)
		if .asof isnt false
			.t.Asof(.asof)
		.cursor = .t.SeekQuery(.query)
		}

	highCostCursor?()
		{
		if not .query.Has?('sort')
			return false
		maxAvgCost = 10000 // based on large customer's database
		cursorCostPerRec = .estimateCostPerRec(.cursor)
		if cursorCostPerRec < maxAvgCost
			return false
		Transaction(read:)
			{ |t|
			queryCostPerRec = .estimateCostPerRec(t.SeekQuery(.query))
			}
		if queryCostPerRec >= cursorCostPerRec
			return false
		return true
		}

	estimateCost: false
	estimateNRecs: false
	estimateAvgCost: false
	estimateCostPerRec(queryOrCursor)
		{
		estimates = queryOrCursor.Strategy().AfterLast('[')
		nRecs = Number(estimates.AfterFirst("nrecs~ ").BeforeFirst(" "))
		.estimateCost = Number(estimates.AfterFirst("cost~ ").
			BeforeFirst(" ").BeforeFirst("]"))
		.estimateNRecs = Max(nRecs, 1)
		.estimateAvgCost = (.estimateCost / .estimateNRecs).Round(2)
		return .estimateAvgCost
		}

	ReadDown(lines) // bottom cursor
		{
		if lines <= 0
			return false

		if .closed
			return 0

		if not .Seeking and .startLast and .Pos is 0 // last record
			return 0

		pre_pos = .Pos
		.Pos = .move(lines)
			{ |rec, pos|
			rows = 1
			if .buildRec(rec).Member?('vl_multi_row_count')
				{
				rows = rec.Extract('vl_multi_row_count')
				for row in rec
					.data[pos++] = row
				}
			else
				.data[pos++] = rec
			rows
			}

		return .Pos - pre_pos
		}

	buildRec(rec)
		{
		if .setupRecord isnt false
			(.setupRecord)(rec)
		return rec
		}

	ReadUp(lines) // top cursor
		{
		if lines <= 0
			return false

		if .closed
			return 0

		if not .Seeking and not .startLast and .Pos is 0 // first record
			return 0

		pre_pos = .Pos

		.Pos = .move(-lines)
			{ |rec, pos|
			rows = 1
			if .buildRec(rec).Member?('vl_multi_row_count')
				{
				rows = rec.Extract('vl_multi_row_count')
				for row in rec.Reverse!()
					.data[--pos] = row
				}
			else
				.data[--pos] = rec
			rows
			}

		return .Pos - pre_pos
		}

	ReleaseDown(lines) // top cursor
		{
		if .closed
			return

		.Pos = .move(lines)
			{ |pos|
			rows = .rows(.data[pos])
			for i in ..rows
				.data.Erase(pos + i)
			rows
			}
		}

	rows(rec)
		{ return rec.Member?('vl_multi_row_count') ? rec.vl_multi_row_count : 1 }

	ReleaseUp(lines) // bottom cursor
		{
		if .closed
			return

		.Pos = .move(-lines)
			{ |pos|
			pos--
			rows = .rows(.data[pos])
			for i in ..rows
				.data.Erase(pos - i)
			rows
			}
		}

	move(lines, block)
		{
		if lines is 0
			return .Pos

		t = Timer()
			{
			startPos = .Pos
			DoWithTran(.t)
				{ |t|
				step = lines.Sign()
				pos = .Pos
				end = .Pos + lines
				action = step > 0 ? 'Next' : 'Prev'
				// Have to do this so that consecutive cursor.Prev and cursor.Next calls
				// will return the same record
				if .prevAction isnt action
					(.cursor[.prevAction])(tran: t)
				.prevAction = action
				do	{
					if false is rowsAdded = .getRec(step, t, block, pos)
						break
					pos += step * rowsAdded
					}
					while step * (pos - end) < 0
				}
			}
		.logIfTooSlow(t, lines, startPos)
		return pos
		}

	logIfTooSlow(t, lines, startPos)
		{
		SlowQuery.LogIfTooSlow(t, .query)
			{ |hash|
			logThresholdSeconds = 10
			qLogSize = 99
			if t > logThresholdSeconds and .hostedSystem?()
				SuneidoLog('ERRATIC: (' $ hash $ ')' $
					' VirtualList reading lines takes more than ' $
						logThresholdSeconds $ ' seconds',
					params: [:lines, cursor?: .t is false, :startPos, endPos: .Pos, :t,
						estimateCost: .estimateCost, estimateNRecs: .estimateNRecs,
						estimateAvgCost: .estimateAvgCost,
						queryHead: .query[..qLogSize], queryTail: .query[-qLogSize..]],
					switch_prefix_limit: 3)
			}
		}

	hostedSystem?()
		{
		return OptContribution('HostedSystem?', function () { return false })()
		}

	getRec(step, t, block, pos)
		{
		rec = step > 0 ? .cursor.Next(tran: t) : .cursor.Prev(tran: t)
		if rec is false
			return false
		return block(:rec, :pos) // each block returns the number of rows added
		}

	Seek(field, prefix, data)
		{
		.data = data
		.Pos = 0
		.cursor.Seek(field, prefix)
		.Seeking = true
		}

	Columns()
		{
		return .cursor.Columns()
		}

	EstimatedSmall?(query, limit)
		{
		try
			Cursor(query) { return .smallQuery?(it, limit) }
		catch(err /*unused*/, 'invalid query')
			Transaction(read:) { return .smallQuery?(it.Query(query), limit) }
		return false
		}

	smallQuery?(q, limit)
		{
		nRecs = q.Strategy().AfterLast('[').AfterFirst("nrecs~ ").BeforeFirst(" ")
		if nRecs.Blank?()
			return false
		return Number(nRecs) <= limit
		}

	closed: false
	Close()
		{
		if .closed
			return
		.closed = true
		.cursor.Close()
		if .t isnt false
			.t.Complete()
		}
	}
