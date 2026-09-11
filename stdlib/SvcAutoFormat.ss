// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
/*
Client side auto formatting for version control sends.

Formatting changes and code changes should never be mixed in one commit. Because a
modified local record still holds lib_before_text (the text as of the last get or send)
we can "time travel" even after the code has been changed we can go back, format the
original text, and commit that by itself. The real commit is then a diff between two
formatted versions i.e. only the actual code change.

BeforeSend is called from Svc.SendLocalChanges for each change before it is Put:

	1. format the original text and commit it as "autoformatted"
	2. format the modified text, which Svc.SendLocalChanges then commits as usual

Once the original is formatted the first step does nothing on later sends because
formatting is idempotent, so a record only ever gets one extra commit.

A record is left completely alone when:

	- it is not a library record, or it has an unresolved merge, or it was moved or
		renamed while already committed, because none of those are a plain text change
	- it has not been committed yet, since a new record is all "your change" and there
		is nothing to separate out, so it is simply formatted in the one commit
	- the change disappears once both sides are formatted, i.e. it is a layout only
		edit, so there is nothing to keep separate and no extra commit is worth making

Nothing is formatted unless AstFmtEquals? confirms the formatting changes the layout
and not the code. If a record can't be compiled it is quietly sent as is. If it can be
compiled but not formatted safely the user is told and it is still sent as is a
formatting problem must never block a real commit.
*/
class
	{
	Desc: #autoformatted

	// returns the asof for the real Put, or false if the auto format commit failed
	BeforeSend(svc, svcTable, name, id, asof, feedbackob = false)
		{
		if not .autoFormat?(svcTable, name)
			return asof
		if false is rec = svcTable.Get(name)
			return asof
		newRecord? = rec.lib_committed is ""
		if not .formattable?(rec, newRecord?)
			return asof
		if false is formattedModified = .Format(name, rec.text)
			return asof
		if newRecord?
			{
			.Apply(svcTable, name, formattedModified)
			return asof
			}
		// a committed record with no before text was moved or renamed rather than
		// edited, and formatting it would mix layout into that commit
		if "" is original = rec.lib_before_text
			return asof
		if false is formattedOriginal = .Format(name, original)
			return asof
		if formattedOriginal is formattedModified
			return asof
		if formattedOriginal is original
			{
			.Apply(svcTable, name, formattedModified)
			return asof
			}
		return .commitOriginal(svc, svcTable, name, id, asof, feedbackob,
			Object(:formattedOriginal, :formattedModified, modified: rec.text))
		}

	autoFormat?(svcTable, name)
		{
		return svcTable.Type is #lib
			? LibraryTags.GetTagFromName(name) not in (#__trial, #__alpha)
			: false
		}

	formattable?(rec, newRecord?)
		{
		if rec.GetDefault(#lib_invalid_text, "") isnt ""
			return false
		return newRecord? or rec.lib_before_path is ""
		}

	commitOriginal(svc, svcTable, name, id, asof, feedbackob, texts)
		{
		.Apply(svcTable, name, texts.formattedOriginal)
		try
			result = svc.Put(svcTable, name, id, .Desc, asof)
		catch (e)
			{
			.Apply(svcTable, name, texts.modified)
			throw e
			}
		if result is false
			{
			.Apply(svcTable, name, texts.modified)
			return false
			}
		.Apply(svcTable, name, texts.formattedModified)
		if feedbackob isnt false
			feedbackob.Add(Object(changeType: ' ', lib: svcTable.Table(),
				name: name $ " - " $ .Desc, prefix: ">>>"))
		return result
		}

	Apply(svcTable, name, text)
		{
		Transaction(update:)
			{|t|
			if false isnt rec = svcTable.Get(name, t)
				if rec.text isnt text
					svcTable.Update(rec, :t, newText: text)
			}
		}

	// returns the formatted text, or false if it can't be formatted safely
	Format(name, text)
		{
		if not Compilable?(text)
			return false
		try
			{
			formatted = AstFormatter(text).Replace("\r?\n", "\r\n")
			layoutOnly? = AstFmtEquals?(text, formatted)
			}
		catch (e)
			return .cantFormat(name, e)
		return layoutOnly?
			? formatted
			: .cantFormat(name, "it would change the code, not just the layout")
		}

	cantFormat(name, reason)
		{
		SuneidoLog.Once("ERROR: SvcAutoFormat: can't format " $ name $ ": " $ reason)
		.alertCantFormat(name, reason)
		return false
		}

	alertCantFormat(name, reason)
		{
		if TestRunner.RunningTests?()
			return
		AlertError("Unable to auto format " $ name $ ":\r\n\t" $ reason $
			"\r\n\r\nIt will be sent unformatted.")
		}
	}
