// Copyright (C) 2026 Axon Development Corporation All rights reserved worldwide.
// BuiltDate > 20260824
class
	{
	RecsPerThread: 500
	MaxThreads:    16
	PollSecs:      2
	TimeoutMins:   60
	SecsPerMin:    60
	Percent:       100
	OneLineMax:    12
	OneLineWidth:  80
	Cols: (lib: 28, count: 6, pct: 3, done: 7)
	DumpFirst: (started, seconds, threads, width, totals)
	NotSuneido: (".css", ".js", ".html", ".htm", ".xml", ".json", ".txt", ".svg")

	CallClass(libs :object = #(), nthreads :number = -1, width :number = 90,
		pretty :boolean = true, maxList :number = 100, checkWarnings :boolean = false,
		output :function|false = false) :object
		{
		started = Date()
		libs = .libraries(libs)
		loaded = .load(libs, output)
		nthreads = .threadCount(nthreads, libs.Size(), .sizeOf(loaded.records))
		verdicts = .sweep(loaded.records, nthreads, width, checkWarnings, output)
		result = .cappedObject(
			.summarize(verdicts, loaded.errors, started, nthreads, width), maxList)
		.dump(result, pretty, maxList, output)
		return result
		}

	libraries(libs :object) :object
		{
		if libs.NotEmpty?()
			return libs
		all = Object()
		all.Append(LibraryTables())
		return all
		}

	emit(out :function|false, s :string)
		{
		if out isnt false
			{
			.tryPrint(out, s)
			return
			}
		try
			ServerPrint(s)
		catch
			.tryPrint(Print, s)
		}

	tryPrint(f :function, s :string)
		{
		try
			f(s)
		catch
			return
		}

	sizeOf(ob :object) :number
		{
		return ob.Size()
		}

	load(libs :object, out :function|false) :object
		{
		.emit(out, "=== AstFormatter sweep ===")
		.emit(out, "loading " $ libs.Size() $ " libraries")
		records = Object()
		errors = Object()
		for lib in libs
			{
			before = records.Size()
			err = .loadLib(lib, records)
			if err isnt false
				errors[lib] = err
			.emit(out,
				"  " $ String(lib).RightFill(.Cols.lib) $
					String(records.Size() - before).LeftFill(.Cols.count) $ " records" $
					(err is false ? "" : "   QUERY FAILED: " $ err))
			}
		.emit(out, "  " $ records.Size() $ " records in " $ libs.Size() $ " libraries")
		return Object(:records, :errors)
		}

	loadLib(lib :string, records :object) :string|false
		{
		try
			QueryApply(lib $ " where group is -1 sort name")
				{|x|
				records.Add(Object(:lib, name: x.name, text: x.text))
				}
		catch (e)
			return String(e)
		return false
		}

	threadCount(nthreads :number, nlibs :number, nrecs :number) :number
		{
		if nthreads > 0
			return Number(Max(1, Min(nthreads, nrecs)))
		byVolume = Number(Max(nlibs, (nrecs / .RecsPerThread).Int()))
		return Number(Max(1, Min(.MaxThreads, nrecs, byVolume)))
		}

	batches(records :object, n :number) :object
		{
		result = Object()
		for (i = 0; i < n; ++i)
			{
			batch = Object()
			for (j = i; j < records.Size(); j += n)
				batch.Add(records[j])
			result.Add(batch)
			}
		return result
		}

	sweep(records :object, nthreads :number, width :number, checkWarnings :boolean,
		out :function|false) :object
		{
		started = Date()
		outputs = Object()
		wg = WaitGroup()
		for batch in .batches(records, nthreads)
			{
			o = Object()
			outputs.Add(o)
			wg.Thread(.check, batch, o, width, checkWarnings, name: #AstFmtSweep)
			}
		.showProgress(wg, outputs, records.Size(), started, out)
		verdicts = Object()
		for o in outputs
			verdicts.Append(o)
		return verdicts
		}

	check(batch :object, results :object, width :number, checkWarnings :boolean)
		{
		for rec in batch
			results.Add(.safeVerdict(rec, width, checkWarnings))
		}

	safeVerdict(rec :object, width :number, checkWarnings :boolean) :object
		{
		try
			return .verdict(rec, width, checkWarnings)
		catch (e)
			return .outcome(rec, #failed, #checkThrew, String(e))
		}

	verdict(rec :object, width :number, checkWarnings :boolean) :object
		{
		text = rec.text
		if not String?(text) or text.Blank?()
			return .outcome(rec, #skipped, #empty, "record has no text")
		name = String(rec.name)
		if .notSuneido?(name)
			return .outcome(rec, #skipped, #notSuneido, "not Suneido source")
		src = RemoveUnderscoreRecordName(name, text)
		err = .compileError(src)
		if err isnt false
			return .outcome(rec, #skipped,
				err.Has?("invalid reference to _") ? #overload : #error, err)
		return .formatVerdict(rec, src, width, checkWarnings)
		}

	formatVerdict(rec :object, src :string, width :number, checkWarnings :boolean) :object
		{
		formatted = false
		try
			formatted = AstFormatter(src, :width)
		catch (e)
			return .outcome(rec, #failed, #threw, String(e))
		if not String?(formatted)
			return .outcome(rec, #failed, #notString,
				"AstFormatter returned " $ Type(formatted))
		err = .compileError(formatted)
		if err isnt false
			return .outcome(rec, #failed, #badOutput, err)
		return .outcome(rec, formatted is src.Tr('\r') ? #ok : #reformatted, "", "",
			checkWarnings ? .extraWarnings(src, formatted, rec) : #(),
			.LongLines(src, formatted, rec))
		}

	outcome(rec :object, status :string, reason :string, detail :string,
		warnings :object = #(), longLines :object|false = false) :object
		{
		return Object(lib: String(rec.lib), name: String(rec.name), :status, :reason,
			:detail, :warnings, :longLines)
		}

	LongLines(before :string, after :string, rec :object) :object
		{
		b = .longLineList(before, rec)
		a = .longLineList(after, rec)
		keys = b.Map({ .lineKey(it.text) })
		return Object(before: b.Size(), after: a.Size(),
			added: a.Filter({ not .wasLong?(keys, .lineKey(it.text)) }))
		}

	lineKey(s :string) :string
		{
		return s.Tr(" \t,#").Tr("'", '"')
		}

	wasLong?(keys :object, key :string) :boolean
		{
		return keys.Any?({ it.Has?(key) or key.Has?(it) })
		}

	longLineList(code :string, rec :object) :object
		{
		rd = [lib: String(rec.lib), recordName: String(rec.name), :code]
		lines = code.Lines()
		result = Object()
		for w in Qc_LineSize(rd, true).lineWarnings
			{
			n = Number(w[0])
			s = String(lines[n-1])
			result.Add(
				Object(line: n, text: s.Trim(), width: s.Detab().RightTrim().Size()))
			}
		return result
		}

	compileError(src :string) :string|false
		{
		try
			Suneido.Compile(src)
		catch (e)
			return String(e)
		return false
		}

	notSuneido?(name :string) :boolean
		{
		lower = name.Lower()
		for ext in .NotSuneido
			if lower.Suffix?(ext)
				return true
		return false
		}

	extraWarnings(before :string, after :string, rec :object) :object
		{
		b = .warningCounts(before, rec)
		a = .warningCounts(after, rec)
		if b is false or a is false
			return #()
		extra = Object()
		for msg in a.Members()
			{
			n = Number(a[msg]) - Number(b.GetDefault(msg, 0))
			for (i = 0; i < n; ++i)
				extra.Add(msg)
			}
		return extra
		}

	warningCounts(code :string, rec :object) :object|false
		{
		found = Object()
		try
			CheckCode(code, rec.name, rec.lib, results: found)
		catch
			return false
		counts = Object()
		for w in found
			{
			msg = String(w.msg).Tr(" \t\r\n", ' ').Trim()
			counts[msg] = Number(counts.GetDefault(msg, 0)) + 1
			}
		return counts
		}

	showProgress(wg, outputs :object, total :number, started :date, out :function|false)
		{
		.emit(out, "")
		.emit(out,
			"sweeping " $ total $ " records with " $ outputs.Size() $ " threads")
		done = false
		last = -1
		while done isnt true
			{
			done = wg.Wait(.PollSecs)
			n = 0
			for o in outputs
				n += .sizeOf(o)
			if n isnt last or done is true
				{
				last = n
				.emit(out, "  " $ .progressLine(n, total, started))
				}
			if Date().MinusSeconds(started) > .TimeoutMins * .SecsPerMin
				{
				.emit(out, "  *** gave up waiting after " $ .TimeoutMins $ " minutes ***")
				break
				}
			}
		}

	progressLine(n :number, total :number, started :date) :string
		{
		secs = Date().MinusSeconds(started)
		pct = total is 0 ? .Percent : (.Percent * n / total).Round(0)
		line = String(pct).LeftFill(.Cols.pct) $ "%  " $ String(n).LeftFill(.Cols.done) $
			'/' $ total $ "   " $ .mmss(secs) $ " elapsed"
		return n is 0 or n >= total
			? line
			: line $ ", ~" $ .mmss(secs * (total - n) / n) $ " left"
		}

	mmss(secs :number) :string
		{
		whole = Number(Max(0, secs)).Round(0)
		return (whole / .SecsPerMin).Int() $ ':' $
			String(whole % .SecsPerMin).LeftFill(2, '0')
		}

	summarize(verdicts :object, loadErrors :object, started :date, nthreads :number,
		width :number) :object
		{
		totals = .newTally()
		libs = Object()
		failed = Object()
		reformatted = Object()
		warned = Object()
		longer = Object()
		for v in verdicts
			{
			id = .idOf(v)
			.tally(totals, v)
			if not libs.Member?(v.lib)
				libs[v.lib] = .newTally()
			.tally(libs[v.lib], v)
			if .longer?(v)
				longer.Add(Object(:id, before: v.longLines.before,
					after: v.longLines.after, added: v.longLines.added))
			if v.status is #failed
				failed.Add(
					Object(:id, reason: String(v.reason), detail: String(v.detail)))
			else if v.status is #reformatted
				reformatted.Add(id)
			if .sizeOf(v.warnings) > 0
				warned.Add(Object(:id, msgs: v.warnings))
			}
		failed.Sort!({|a, b| a.id < b.id })
		warned.Sort!({|a, b| a.id < b.id })
		longer.Sort!({|a, b| a.id < b.id })
		reformatted.Sort!()
		return Object(:started, seconds: Date().MinusSeconds(started).Round(0),
			threads: nthreads, :width, :totals, :libs, :failed,
			skipped: .collectSkips(verdicts), :reformatted, :warned, :longer,
			:loadErrors)
		}

	longer?(v :object) :boolean
		{
		return v.longLines isnt false and v.longLines.after > v.longLines.before
		}

	idOf(v :object) :string
		{
		return String(v.lib) $ ':' $ String(v.name)
		}

	collectSkips(verdicts :object) :object
		{
		error = Object()
		overload = Object()
		notSuneido = Object()
		empty = Object()
		for v in verdicts
			{
			if v.status isnt #skipped
				continue
			id = .idOf(v)
			if v.reason is #error
				error.Add(Object(:id, detail: String(v.detail)))
			else if v.reason is #overload
				overload.Add(id)
			else if v.reason is #notSuneido
				notSuneido.Add(id)
			else
				empty.Add(id)
			}
		error.Sort!({|a, b| a.id < b.id })
		overload.Sort!()
		notSuneido.Sort!()
		empty.Sort!()
		return Object(:error, :overload, :notSuneido, :empty)
		}

	newTally() :object
		{
		return Object(records: 0, ok: 0, reformatted: 0, skipped: 0, failed: 0,
			longBefore: 0, longAfter: 0, longer: 0)
		}

	tally(t :object, v :object)
		{
		t.records = Number(t.records) + 1
		t[v.status] = Number(t[v.status]) + 1
		if v.longLines is false
			return
		t.longBefore = Number(t.longBefore) + Number(v.longLines.before)
		t.longAfter = Number(t.longAfter) + Number(v.longLines.after)
		if .longer?(v)
			t.longer = Number(t.longer) + 1
		}

	dump(result :object, pretty :boolean, maxList :number, out :function|false)
		{
		.emit(out, "")
		.emit(out,
			"--- RESULT (lists elided past maxList " $ maxList $
				"; raise it to see more) ---")
		if pretty
			.prettyPrint(result, "", "", out)
		else
			.emit(out, .safeDisplay(result))
		}

	safeDisplay(x) :string
		{
		try
			return Display(x)
		catch
			return '\xab' $ "too large to display" $ '\xbb'
		}

	capped(x, max :number)
		{
		return Object?(x) ? .cappedObject(x, max) : x
		}

	cappedObject(ob :object, max :number) :object
		{
		res = Object()
		nlist = ob.Size(list:)
		shown = Number(Min(nlist, max))
		for (i = 0; i < shown; ++i)
			res.Add(.capped(ob[i], max))
		if nlist > shown
			res.Add("... " $ (nlist - shown) $ " more")
		for m in ob.Members(named:)
			res[m] = .capped(ob[m], max)
		return res
		}

	prettyPrint(val, indent :string, label :string, out :function|false)
		{
		head = label is "" ? "" : label $ ": "
		if not Object?(val) or .oneLine?(val)
			{
			.emit(out, indent $ head $ .safeDisplay(val))
			return
			}
		.emit(out, indent $ head $ "Object(")
		for m in .memberOrder(val)
			.prettyPrint(val[m], indent $ "    ", Number?(m) ? "" : String(m), out)
		.emit(out, indent $ ')')
		}

	memberOrder(val :object) :object
		{
		order = val.Members(list:)
		named = val.Members(named:).Sort!()
		for m in .DumpFirst
			if named.Has?(m)
				order.Add(m)
		for m in named
			if not .DumpFirst.Has?(m)
				order.Add(m)
		return order
		}

	oneLine?(val :object) :boolean
		{
		if val.Size() > .OneLineMax
			return false
		for m in val.Members()
			if Object?(val[m])
				return false
		return .safeDisplay(val).Size() <= .OneLineWidth
		}
	}
