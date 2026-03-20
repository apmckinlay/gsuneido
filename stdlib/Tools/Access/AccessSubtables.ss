// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	linkedSelectFields: false
	subTableMgrs: false
	New(.dynamic? = false)
		{
		.linkedSelectFields = Object()
		.subTableMgrs = Object()
		}

	SetSubtableSelectMgr(fields, saveName)
		{
		// Don't recreate if it already exists
		if .subTableMgrs.Member?(saveName)
			{
			selectMgr = .subTableMgrs[saveName]
			}
		else
			{
			selectMgr = .initSelectMgr(saveName, fields)
			.subTableMgrs[saveName] = selectMgr
			}

		// need to do this on the load, not the save, as the selectMgr for the
		// current dynamic type does NOT exist until the first time we open select
		// for it specifically
		if .dynamic? is true
			.handleDynamic(selectMgr.Select_vals(), saveName, .updateChecked, selectMgr)

		return selectMgr
		}

	initSelectMgr(saveName, fields)
		{
		selectMgr = AccessSelectMgr(#(), saveName)
		selectMgr.LoadSelects(fields, false)
		return selectMgr
		}

	// loops through each saved condition for every dynamic type that matches
	// the current saveName
	forEachSavedCondition(saveName, block)
		{
		filter = { |x|  x.AfterFirst('|') is saveName.AfterFirst('|') }
		for table in .subTableMgrs.MembersIf(filter)
			{
			// need to look up the select fields at this point as well
			ssf = .linkedSelectFields[table]
			savedConditions = .subTableMgrs[table].Select_vals()

			for sCondition in savedConditions
				{
				block(sCondition, ssf)
				}
			}
		}

	conditionsMatch(condition1, condition2)
		{
		return condition1.condition_field is condition2.condition_field and
			condition1[condition1.condition_field] is
				condition2[condition2.condition_field]
		}

	handleDynamic(currentConditions, saveName, updateFn, selectMgr = false)
		{
		// Need to check for filters that have the same linkedBrowse name
		// in saveName it's anything after the |
		csf = .linkedSelectFields[saveName]
		.forEachSavedCondition(saveName)
			{ |sCondition, ssf|
			updateFn(sCondition, currentConditions, csf, ssf, selectMgr)
			}
		}

	// situation: we HAVE a select, have navigated to a new record that is a different Dynamic Type
	// and re-opened the select. Want to make sure that our "select" gets added
	// sCondition: condition read from the "saved" conditions (will have the "checked")
	// currentConditions: condition for the record we are currently on (will NOT have the checked)

	// this should NOT modify sCondition, only currentConditions and/or selectMgr
	updateChecked(sCondition, currentConditions, csf, ssf, selectMgr)
		{
		// the ASSUMPTION is that if we are SITTING on the record, that the saved select
		// condition IS applicable to the current record. (I.E. we shouldn't need
		// to worry about adding a condition that does not match the current available
		// columns)
		// this should be a safe assumption as when we "save" the condition we update
		// all the saved conditions to clear conditions that have been unchecked/unselected

		// however, using locate can get around this, We belive this is NOT an issue
		// as if they end up on a record that does NOT match their select
		// this still works the same as if they've done it with the main select
		// even if they end up on a record where their select no longer shows up
		// and no records exist that match the select, the select control just
		// wipes out the checked conditions anyway when they try to navigate
		if sCondition.check isnt true
			return

		// don't set this back to sCondition, as that would cause a side-effect
		// this function SHOULD NOT change the passed in sCondition
		rCondition = .renamedCondition(sCondition, ssf, csf)

		for cCondition in currentConditions
			{
			// if we have a matching current condition, ensure it is checked
			if .conditionsMatch(rCondition, cCondition)
				{
				cCondition.check = true
				return // Next Saved Condition
				}
			}

		// we have a saved condition without a matching current condition
		// need to add the checked one in
		selectMgr.Select_vals().Add(rCondition)
		}

	// Need to handle when two different Dynamic Types, use the same subtable
	// but have the same column renamed differently
	// using prompts from select fields to track what field should be
	// set in conditions
	// This assumes that the renamed fields have the same prompt
	renamedCondition(sCondition, ssf, csf)
		{
		newCondition = Record()

		newCondition.condition_field = csf.PromptToField(
			ssf.FieldToPrompt(sCondition.condition_field))
		newCondition[newCondition.condition_field] =
			sCondition[sCondition.condition_field]
		newCondition.check = sCondition.check
		return newCondition
		}

	// this should NOT modify currentConditions, only sCondition
	updateUnchecked(sCondition, currentConditions, csf, ssf, selectMgr /*unused*/ = false)
		{
		// Need to check saved conditions for checked values that have been
		// unselected or removed on currentConditions

		// only care about saved conditions that ARE checked
		if sCondition.check is false
			return

		found = false
		for cCondition in currentConditions
			{
			// don't set this back to cCondition, as that would cause a side-effect
			// this function SHOULD NOT change the passed in currentConditions
			rCondition = .renamedCondition(cCondition, csf, ssf)
			// find a matching current condition
			if not .conditionsMatch(sCondition, rCondition)
				continue
			found = true

			// if matching current condition found, and is unchecked,
			// uncheck the saved condition
			if rCondition.check isnt true
				sCondition.check = false

			// matching current condition was found and handled, move on to
			// next saved condition
			return
			}
		// saved condition is true, but we DO NOT HAVE a matching current condition
		// need to uncheck the saved condition
		if not found and sCondition.check is true
			sCondition.check = false
		}

	// conditions is the conditions from the select repeate control
	SetSubTableSelectVals(saveName, conditions)
		{
		if false is sf = .getSubtableSelectFields(saveName)
			throw "SelectFields for " $ saveName $ " used, but not initialized"

		selectMgr = .subTableMgrs[saveName]
		if .dynamic? is true
			.handleDynamic(conditions, saveName, .updateUnchecked)
		// we CANNOT handle adding the check here as not all selectMgrs needed
		// will be constructed

		selectMgr.SetSelectVals(conditions, sf)
		}

	SaveSelects()
		{
		if .subTableMgrs isnt false
			.subTableMgrs.Each({ it.SaveSelects() })
		}

	SetSubtableSelectFields(columns, saveName)
		{
		if false isnt sf = .getSubtableSelectFields(saveName)
			return sf
		sf = .initSelectFields(columns, :saveName)
		.linkedSelectFields[saveName] = sf
		return sf
		}

	initSelectFields(columns)
		{
		return SelectFields(columns, includeMasterNum:)
		}

	getSubtableSelectFields(saveName)
		{
		if .linkedSelectFields.Member?(saveName)
			return .linkedSelectFields[saveName]
		return false
		}
	Clear()
		{
		for selectMgr in .subTableMgrs
			{
			selectMgr.Select_vals().Each({ it.check = false })
			}
		}
	}