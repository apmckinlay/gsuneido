// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(listOld, listNew, base = false)
		{
		.diffs = base is false
			// Because Diff.Three always compares base(old) to merged(new),
			// Diff.SideBySide should do the same thing to be consistent
			? Diff.SideBySide(listOld, listNew)
			: Diff.Three(base, listOld, listNew)
		.rowIndics = Object()
		}

	Getter_Diffs()
		{
		return .diffs
		}

	GetRowIndics(row)
		{
		if .rowIndics.Member?(row)
			return .rowIndics[row]

		.rowIndics[row] = Object()
		tokensOld = .getTokens(.diffs[row][0])
		tokensNew = .getTokens(.diffs[row][2])

		rowDiffOldToNew = Diff(tokensOld[0], tokensNew[0])
		rowDiffNewToOld = Diff(tokensNew[0], tokensOld[0])
		rowDiffSideBySide = Diff.SideBySide(tokensOld[0], tokensNew[0])

		numChanges = .getNumChanges(rowDiffOldToNew, rowDiffNewToOld)
		hasSimilarWord? = .hasSimilarWord?(rowDiffSideBySide)
		percTokensChanged = .percTokensChanged(rowDiffOldToNew, rowDiffNewToOld,
			tokensOld[0], tokensNew[0])

		.rowIndics[row].OldToNew =
			.createRowIndics(numChanges, hasSimilarWord?, percTokensChanged,
				rowDiffOldToNew, tokensOld)
		.rowIndics[row].NewToOld =
			.createRowIndics(numChanges, hasSimilarWord?, percTokensChanged,
				rowDiffNewToOld, tokensNew)

		return .rowIndics[row]
		}

	getTokens(s)
		{
		scanner = Scanner(s)
		tokens = Object()
		tokensAddr = Object(0)
		size = 0
		while scanner isnt token = scanner.Next()
			{
			if token.Blank?()
				{
				tmp = token.Divide(1)
				tokens.Add(@tmp)
				tokensAddr.Add(@tmp.Map({ size += it.Size() }))
				}
			else
				{
				tokens.Add(token)
				tokensAddr.Add(size += token.Size())
				}
			}
		return Object(tokens, addr: tokensAddr)
		}

	getNumChanges(rowDiff1, rowDiff2)
		{
		numChanges1 = rowDiff1.CountIf({ it[1] is 'D' })
		numChanges2 = rowDiff2.CountIf({ it[1] is 'D' })

		// Return the greater numChanges to avoid highlighting one line (numChanges > 1)
		// and leaving the other blank (numChanges <= 1)
		return Max(numChanges1, numChanges2)
		}

	specialChar: #(`.`, `,`, `/`, `\`, `=`, `;`, `:`, `(`, `)`, `{`, `}`, `[`, `]`)
	hasSimilarWord?(rowDiffSideBySide)
		{
		return rowDiffSideBySide.Any?({ |word|
			word[1] is "" and not word[0].Blank?() and not .specialChar.Has?(word[0]) })
		}

	createRowIndics(numChanges, hasSimilarWord?, percTokensChanged, rowDiffs, tokens)
		{
		// if there is more than one changed word
		// or there is exactly one deletion/insertion
		if numChanges > 1 or rowDiffs.Size() < 2
			{
			// and if the line has at least one similar word
			// and the percentage of word tokens changed <= 40%
			// then highlight the line regularly
			minPercent = 0.4
			return hasSimilarWord? and percTokensChanged <= minPercent
				? .generateRowIndics(rowDiffs, tokens)
				: Object()
			}
		else // only one text token different
			{
			// if there is more than one non-blank token on the line
			if tokens[0].CountIf({ not it.Blank?() }) <= 1
				return Object()

			// and the change was a consecutive insertion/deletion of chrs
			wordDiff = Diff(rowDiffs[1][2], rowDiffs[0][2])
			if .consecChrsOnly?(wordDiff)
				{
				// highlight character-specific
				offset = tokens.addr[rowDiffs[0][1] is 'D'
					? rowDiffs[0][0]
					: rowDiffs[1][0]]
				return .generateRowIndics(wordDiff, :offset)
				}
			// else, just highlight the word
			return .generateRowIndics(rowDiffs, tokens)
			}
		}

	percTokensChanged(rowDiff1, rowDiff2, tokens1, tokens2)
		{
		return Min(.getPercTokensChanged(rowDiff1, tokens1),
			.getPercTokensChanged(rowDiff2, tokens2))
		}

	getPercTokensChanged(rowDiffs, tokens)
		{
		deletedWordCount = rowDiffs.CountIf({
			it[1] is 'D' and .specialChar.Has?(it[2]) is false and not it[2].Blank?()})
		wordTokenCount = tokens.CountIf({
			.specialChar.Has?(it) is false and not it.Blank?() })
		return (deletedWordCount / wordTokenCount)
		}

	consecChrsOnly?(wordDiff)
		{
		return .consecInserts?(wordDiff.Filter({ it[1] is 'I' })) and
			.consecDeletes?(wordDiff.Filter({ it[1] is 'D' }))
		}

	// Consecutive inserts should all have the same index
	consecInserts?(inserts)
		{
		idx = false
		for chr in inserts
			if idx is false
				idx = chr[0]
			else
				{
				if chr[0] isnt idx
					return false
				}
		return true
		}

	// Consecutive deletes should have consecutive indexes
	consecDeletes?(deletes)
		{
		idx = false
		for chr in deletes.Sort!({ |x,y| x[0] < y[0] })
			if idx is false
				idx = chr[0]
			else
				{
				if chr[0] - idx > 1
					return false
				idx = chr[0]
				}
		return true
		}

	generateRowIndics(itemDiffs, items = false, offset = 0)
		{
		itemDiffs.
			Filter({ it[1] is 'D' }).
			Map({ Object(
				pos: (items is false ? it[0] : items.addr[it[0]]) + offset,
				length: it[2].Size()) })
		}
	}