// Copyright (C) 2000 Suneido Software Corp. All rights reserved worldwide.
// called by browser for suneido: url's
class
	{
	FileRegx: "^suneido:/[^/]+/res/.*\.[\w]+$"
	CallClass(url)
		{
		url = .decode(url)
		// NOTE: may not need to embed images into html for printing anymore if gSuneido
		// handles it properly
		embed? = url.Suffix?('?print')
		url = url.RemoveSuffix('?print')
		if url.Prefix?("suneido:/eval?")
			{
			.evalCode(url, url[14 ..] /*= strip "suneido:/eval?" */,
				"ERROR: SuneidoAPP - ")
			return ""
			}
		if url.Prefix?("suneido:/from?")
			return .evalCode(url, url[14 ..] /*= strip "suneido:/from?" */,
				"ERROR: (CAUGHT) SuneidoAPP - ", returnMsg: "Resource Execution Failure")
		if url.Prefix?("suneido:/inmemory/")
			return InMemory.Get(url)
		return not .going?(url)
			? .app(url, embed?)
			: ''
		}

	decode(url)
		{
		if BuiltDate() >= #20250307
			url = url.BeforeFirst('#')
		if false is pos = url.FindLast('?')
			pos = url.Size()
		return Url.Decode(url[..pos]) $ url[pos..]
		}

	evalCode(url, code, logPrefix, returnMsg = "")
		{
		code = Url.Decode(code)
		if not .okayToEval?(code)
			{
			SuneidoLog("ERROR: SuneidoAPP - Code is not okay to eval",
				params: Object(:url))
			return ""
			}
		try
			return code.Eval()
		catch (err)
			{
			SuneidoLog(logPrefix $ err, params: Object(:url),
				caughtMsg: Opt("returned: ", returnMsg))
			return returnMsg
			}
		}

	okayToEval?(code)
		{
		allowed = #(NotificationsHtml OpenBook ShellExecute AccessGoTo Exec_DrillDown,
			LoginHtml, ChangePasswordDialog, 'SystemInfo.ShowScript')
		return allowed.Has?(code.BeforeFirst('('))
		}

	going?(url)
		{
		if url !~ .FileRegx
			try
				return PubSub.Publish(.event(.hwnd(url)), 'Going', :url) isnt false
			catch (err)
				{
				if err isnt 'no return value'
					SuneidoLog('ERROR: (CAUGHT) ' $ err, params: Object(:url),
						caughtMsg: 'url navigation failed')
				return true
				}
		return false
		}

	hwnd(url)
		{
		if false is hwnd = .browserHwnd(url, browserLoads = .browserLoads())
			return false
		browserLoads[hwnd].Delete(url)
		SetActiveWindow(hwnd)
		return hwnd
		}

	browserLoads()
		{
		return Suneido.GetDefault('BrowserRedir_Loads', false)
		}

	browserHwnd(url, browserLoads)
		{
		browserHwnd = false
		oldest = Date.End()
		if Object?(browserLoads)
			for hwnd, loads in browserLoads
				for load, date in loads
					if load is url and date < oldest
						{
						browserHwnd = hwnd
						oldest = date
						}
		return browserHwnd
		}

	event(hwnd = false)
		{
		return 'BrowserRedir_' $ (hwnd is false ? GetActiveWindow() : hwnd)
		}

	app(url, embed?)
		{
		url = url[8 ..] /*= strip "suneido:" */
		if not url.Prefix?("/")
			return ''
		rec = .GetBookRec(url)
		if String?(rec)
			{
			SuneidoLog('ERROR: (CAUGHT) SuneidoAPP - ' $ rec, calls:,
				params: Object(:url), caughtMsg: 'returned: ' $ rec)
			return rec
			}
		if url =~ .Images $ '$'
			return rec.text

		// set dynamic variables so they're accessible by Asup expressions
		_table = rec.table
		_path = rec.path
		_name = rec.name
		return rec.text.Prefix?("<")
			? HtmlWrap(rec.text, rec.table, :embed?)
			: .programPage(rec.text, rec.table, rec.path, rec.name)
		}

	Images: "(?i)[.](gif|jpg|png|bmp|svg)"
	GetBookRec(url)
		{
		table = url.Extract("^/([^/]*)")
		if not TableExists?(table)
			return Display(table) $ ' is not a table'
		name = url.Extract("[^/]*$")
		path = url[1 + table.Size() .. -name.Size() - 1]
		if (not name.Suffix?('?') and (false isnt pos = name.FindLast('?')) and
			name[pos+1] isnt ' '/*? isnt in the middle of something*/)
			name = name[..pos]
		qFn = LastContribution('CacheBookRec?')(table, path, name)
			? Query1Cached
			: Query1
		if false is x = qFn(table, :path, :name)
			return "Page " $ Display(table) $ " > " $
				Display(path) $ " > " $ Display(name) $ " not found." $
				" Page may have been deleted or renamed."
		x.table = table
		return x
		}

	programPage(s, table, path, name)
		{
		try
			x = s.Eval() // needs Eval
		catch (unused, "can't find")
			x = "<p>This page is not available.</p>"
		if (String?(x) and x.Prefix?("<"))
			return HtmlWrap(x, table)
		if ("" isnt (s = SuneidoAPP_Authorize(table, path $ '/' $ name)))
			return HtmlWrap(s)
		PubSub.Publish(.event(), 'Set', control: x)
		return ""
		}
	}
