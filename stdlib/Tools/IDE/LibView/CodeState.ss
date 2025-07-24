// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
/*
This class handles saving invalid code seperate from valid code. This is to prevent errors
when working in "live" code, (ie: LibViewCoreControl, LibTreeModel, etc).

By "invalid", we mean code which has compilation issues. As long as it does not
compile, it will be saved seperate from the loaded / valid code

"Valid" code, has no compilation errors. This code has had "Unload" called on it and
is the code currently being used.
*/
class
	{
	New(.libTreeModel)
		{ }

	hasInvalidHandling?(lib)
		{ return QueryColumns(lib).Has?(#lib_invalid_text) }

	CallClass(lib, rec)
		{
		libTreeModel = new LibTreeModel
		num = libTreeModel.MangleNum(lib = .toggleLib(lib), rec.num)
		if false isnt libRec = libTreeModel.Get(num)
			{
			libRec.text = rec.text
			(new this(libTreeModel)).Save(libRec)
			}
		else
			.log('Failed to look up record. Record not updated', [name: rec.name, :lib])
		}

	log(message, params = '')
		{ SuneidoLog('ERROR: (CAUGHT) CodeState: ' $ message, :params, caughtMsg: 'IDE') }

	toggleLib(lib)
		{
		lib = lib.Tr('()')
		return Libraries().Has?(lib) ? lib : '(' $ lib $ ')'
		}

	// Save should only ever be called on a Library record
	Save(rec)
		{
		if not .hasInvalidHandling?(lib = .libTreeModel.TableName(rec.num).Tr('()'))
			Database('ensure ' $ lib $ ' (lib_invalid_text)')
		if not rec.Member?(#text)
			return .libTreeModel.Update(rec)
		if not .Valid?(lib, rec)
			{
			rec.lib_invalid_text = rec.text
			rec.Delete(#text)
			}
		else
			rec.lib_invalid_text = ''
		rec.Invalidate(#lib_current_text)
		return .libTreeModel.Update(rec)
		}

	Valid?(lib, rec)
		{
		if false is valid? = rec.GetDefault(#valid?, true)
			return false
		try
			valid? = CheckCode(rec.text, .validateName(rec.name), lib) isnt false
		catch
			valid? = false
		return valid?
		}

	validateName(name)
		{ return name.Tr('()') }

	InvalidRec(lib, name)
		{
		if not .hasInvalidHandling?(lib)
			return false
		rec = Query1(lib, name: .validateName(name), group: -1)
		return rec is false or rec.lib_invalid_text is '' ? false : rec
		}

	InvalidRecs(lib)
		{
		if not .hasInvalidHandling?(lib)
			return []
		return QueryAll(lib $ ' where lib_invalid_text isnt "" and group is -1').
			Map!({ lib $ ':' $ it.name })
		}

	// This should only be used to run tests. The purpose of CodeState is to ensure
	// that any invalid code is kept separate from its unloaded version. This function
	// temporarily replaces the unloaded code with the invalid. IF a catastrophic
	// failure occurred the invalid code could be left in place, potentially
	// leading to the exact situation this class was designed to avoid
	// -
	// In short: Use with caution
	RunCurrentCode(lib, name, block)
		{
		if Record?(rec = .InvalidRec(lib, name))
			{
			origText = .imposeCode(lib, name, rec)
			error = false
			try
				block()
			catch (e)
				error = e
			.restoreCode(lib, name, origText)
			if error isnt false
				throw error
			}
		else
			block()
		}

	imposeCode(lib, name, rec)
		{
		origRec = Query1(lib, :name, group: -1)
		QueryDo('update ' $ lib $ ' where name is ' $ Display(name) $ 'and group is -1
			set text = ' $ Display(rec.lib_current_text))
		LibUnload(name)
		return origRec.text
		}

	restoreCode(lib, name, origTex)
		{
		QueryDo('update ' $ lib $ ' where name is ' $ Display(name) $ ' and group is -1
			set text = ' $ Display(origTex))
		LibUnload(name)
		}
	}
