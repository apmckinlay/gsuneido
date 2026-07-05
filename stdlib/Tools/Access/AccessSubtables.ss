// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	linkedSelectFields: false
	subTableMgrs: false
	// csf - Current Select fields
	// 		- The Select Fields based on the columns as they appear on the CURRENT record
	// casf - Current Available Select Fields
	//		- The Select Fields Based Soley on the Sub-Table Columns themselves
	// ssf - Saved Select Fields
	//		- The Select Fields for the SAVED condition we are trying to apply
	// sasf - Saved Available Select Fields
	//		- SHOULD be identical to casf, BUT IS NOT!
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
			{
			// be aware, the other places calling .handleDynamic, the conditions
			// are NOT the conditions from the selectMgr, which is why these are two
			// different parameters
			.handleDynamic(selectMgr.Select_vals(), saveName, .updateChecked, selectMgr)
			}

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
			// ignore the Current Filter
			if table is saveName
				continue
			// need to look up the select fields at this point as well
			result = .linkedSelectFields[table]
			ssf = result.sf
			asf = result.asf
			savedConditions = .subTableMgrs[table].Select_vals()

			for sCondition in savedConditions
				{
				block(sCondition, ssf, asf)
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
		result = .linkedSelectFields[saveName]
		csf = result.sf
		casf = result.asf
		.forEachSavedCondition(saveName)
			{ |sCondition, ssf, sasf|
			updateFn(sCondition, currentConditions, csf, casf, ssf, sasf, selectMgr)
			}
		}

	// situation: we HAVE a select, have navigated to a new record that is a different Dynamic Type
	// and re-opened the select. Want to make sure that our "select" gets added
	// sCondition: condition read from the "saved" conditions (will have the "checked")
	// currentConditions: condition for the record we are currently on (will NOT have the checked)

	// this should NOT modify sCondition, only currentConditions and/or selectMgr
	updateChecked(sCondition, currentConditions, csf, casf, ssf, sasf, selectMgr)
		{
		// We can end up on a Dynamic Type for which the current Select is
		// technically Invalid (have a select on a field that is NOT applicable
		// to the Dynamic Type) This code handles ensuring that the Conditions
		// ends up with a value that will work for the current Dynamic Type.
		// Step 1: Check the Current Select Fields (csf) for a matching Field
		// 		this handles if the field DOES exist, but has just been renamed differntly
		// Step 2: If a matching Field is NOT found, need to check the Select Fields
		// for the Table itsef (casf)
		// Once we have a matching field, Step 3: check the Current Conditions to
		// see if it already exists, if it does, ensure it is marked checked
		// if it is not, Then insert it

		if sCondition.check isnt true
			return

		// don't set this back to sCondition, as that would cause a side-effect
		// this function SHOULD NOT change the passed in sCondition
		rCondition = .renamedCondition(sCondition, ssf, sasf, csf, casf)

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
	// First we try based on the prompt, we look for a field in csf/casf that has
	// the same prompt, and use it.
	// if a matching prompt is NOT found, we need to add the field to csf ourselves
	renamedCondition(sCondition, ssf, sasf /*unused*/, csf, casf)
		{
		newCondition = Record()
		sPrompt = ssf.FieldToPrompt(sCondition.condition_field)

		// Handle if the select HAS found a dynamic record that does NOT
		// contain condition_field in the sub table (i.e. "is empty" filter)
		// temporarily add the condition to the Current Select Fields so
		// it will show up on the select.

		if false is cField = csf.PromptToField(sPrompt)
			{
			if false is cField = casf.PromptToField(sPrompt)
				{
				// At this point it's fair to assume that the field from
				// sCondition has been renamed in a way that we do not have it in casf
				// need to manually add it
				if ssf.NameAbbrev?(sCondition.condition_field)
					{
					type = sCondition.condition_field.Has?('_name') ? '_name' : '_abbrev'
					suffix = sCondition.condition_field.AfterFirst(type)
					prefix = sCondition.condition_field.BeforeFirst(type)
					csf.AddNumField(prefix $ '_num', casf.FieldToPrompt(prefix $ '_num'))
					cField = prefix $ type
					}
				else
					{
					// If we ever get here there is still the potential for bugs.
					csf.AddField(sCondition.condition_field, sPrompt)
					cField = sCondition.condition_field
					}
				}
			else //cField came from casf, not csf, need to add it to csf
				{
				if ssf.NameAbbrev?(cField)
					{
					type = cField.Has?('_name') ? '_name' : '_abbrev'
					suffix = cField.AfterFirst(type)
					prefix = cField.BeforeFirst(type)
					csf.AddNumField(prefix $ '_num' $ suffix,
						casf.FieldToPrompt(prefix $ '_num' $ suffix))
					}
				else
					{
					csf.AddField(cField, sPrompt)
					}
				}
			} // else cField came from csf, just add it.
		newCondition.condition_field = cField

		newCondition[newCondition.condition_field] =
			sCondition[sCondition.condition_field]
		newCondition.check = sCondition.check
		return newCondition
		}

	// this should NOT modify currentConditions, only sCondition
	updateUnchecked(sCondition, currentConditions, csf, casf, ssf, sasf,
		selectMgr /*unused*/ = false)
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
			rCondition = .renamedCondition(cCondition, csf, casf, ssf, sasf)
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

	// conditions is the conditions from the select repeat control
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

	SetSubtableSelectFields(columns, availableColumns, saveName)
		{
		if false isnt sf = .getSubtableSelectFields(saveName)
			return sf
		result = .initSelectFields(columns, availableColumns, :saveName)
		.linkedSelectFields[saveName] = result
		return result.sf
		}

	initSelectFields(columns, availableColumns)
		{
		sf = SelectFields(columns, includeMasterNum:)
		// DO NOT EXPOSE asf, just sf
		asf = SelectFields(availableColumns, includeMasterNum:)
		return Object(:sf, :asf)
		}

	getSubtableSelectFields(saveName)
		{
		if .linkedSelectFields.Member?(saveName)
			return .linkedSelectFields[saveName].sf
		return false
		}
	Clear()
		{
		for selectMgr in .subTableMgrs
			{
			selectMgr.Select_vals().Each({ it.check = false })
			}
		}

	Select_vals()
		{
		x = Object()
		for name in .subTableMgrs.Members()
			x[name] = Record(conditions: .subTableMgrs[name].Select_vals())
		return x
		}

	// HELPER METHODS
	// selectControls is for the scenario when we are trying to add the selects from
	// the List Mode in MultiViewControl. Neither the AccessSubtables, nor the AccessControl
	// have a way to know what the VirtualList has for the Selects. It currently handles
	// this by having the VirtualList send a message WITH the selectControls in it
	// and we handle the selects here.
	Where(access, selectControls = false)
		{
		linkedBrowses = access.GetLinkedBrowseTabs()

		lineItemWhere = ''
		dynamicIgnoreHeader = ""
		for idx in linkedBrowses.Members().Sort!()
			{
			linked = linkedBrowses[idx]
			linkedBrowse = linked.browse

			// ignoreLinked Explaination:
			// On DynamicTypes we need to know if we are potentially adding a semi-join
			// to records that have "fake" entries.
			// i.e. no actual linkedbrowse on the screen, but a record in the subtable
			// (e.g. Credit Card Payments on Payables > Checks
			// they do NOT have a linked browse (unlike OTHER AP Checks), but still put a
			// fake record in the ap_checklines table, causing the semi-join to find a
			// Credit Card Payment when filtering on AP Check Lines anyway.)
			// ignoreLinked adds an extra where clause to ignore those types
			// (e.g. 'where apchk_type isnt "CC"' in the above example)
			// need to make sure this ONLY gets added if there sub table selects
			if false isnt (dynamicIgnore =
				access.IgnoreLinkedBrowseType(linkedBrowse.Name))
				dynamicIgnoreHeader $= ' ' $ dynamicIgnore

			// if linkedBrowse is actualy a virtual list we CANNOT use .GetQuery
			// here, as that has an aditional where on it that limits records
			// to the current record. (BrowseControl.GetQuery does NOT have that issue)
			// need to use .GetLinkedQuery as it returns the query without the
			// additional where. BrowseControl does NOT have .GetLinkedQuery()
			q = linkedBrowse.Base?(BrowseControl)
				? linkedBrowse.GetQuery()
				: linkedBrowse.GetLinkedQuery()

			q = QueryHelper.StripSort(q)

			overrides = Object()
			subTablesConfig = access.SubTablesConfig()
			if Object?(subTablesConfig)
				if subTablesConfig.Member?(linkedBrowse.Name)
					overrides = subTablesConfig[linkedBrowse.Name]

			// see comments in .layoutSubtables
			// this is to make sure the correct rename gets used in the
			// semijoin when we have to override the base queries rename
			selectCols = linkedBrowse.GetColumns().Copy()
			extraRename = ""
			for field in overrides.Members()
				{
				if not selectCols.Has?(field)
					continue
				extraRename $= field $ " to " $ overrides[field] $ ' '
				}
			if extraRename isnt ""
				q $= ' rename ' $ extraRename


			saveName = .subTableFilterName(linkedBrowse)

			// Select Repeat Control
			if selectControls is false
				{
				// THIS is for when we are toggling MultiView Control
				if not .linkedSelectFields.Member?(saveName)
					continue
				lineItemSf = .linkedSelectFields[saveName].sf
				// BEWARE: Select_vals on the AccessSelectMgr is NOT a Getter_
				conditions = .subTableMgrs[saveName].Select_vals()
				where = SelectRepeatControl.BuildWhere(lineItemSf, conditions)
				}
			else
				{
				// This is for when we are just on access control
				subTableSelectCtrl = selectControls[saveName]
				lineItemSf = subTableSelectCtrl.FieldPrompt_GetSelectFields()
				// this NEEDS to be asf, not sf
				asf = .linkedSelectFields[saveName].asf
				conditions = subTableSelectCtrl.Get().conditions
				if false is where = SelectRepeatControl.BuildWhere(asf, conditions)
					continue
				}

			// NEED to ensure we call SetSubTableSelectVals on the AccessSubtables
			// on the access itself, as we could be in an un-instaciated class
			access.SetSubTableSelectVals(saveName, conditions)

			if where.where isnt ''
				{

				lineItemWhere $= .semiJoinWhere(linkedBrowse,
					lineItemSf.Joins(where.joinflds) $ where.where, access, q)
				}
			}

		if lineItemWhere is ""
			return ''
		return dynamicIgnoreHeader $ '\r\n' $ lineItemWhere
		}
	// this name is the saved name in userselects
	subTableFilterName(linkedBrowse)
		{
		return linkedBrowse.GetColumnsSaveName() $ ' Filter'
		}

	semiJoinWhere(linkedBrowse, where, access, subTableQuery)
		{
		hdrQuery = access.GetQuery()
		linkedField = linkedBrowse.GetLinkField()
		if Object?(linkedField)
			{
			// for cases where we are semijoining to a view, or the linked field is
			// NOT explicitly a foreign key on the sub-table
			joinField = linkedField[0]
			linkedField = linkedField[1]
			}
		else if false is joinField = .findJoinField(hdrQuery, subTableQuery, access)
			{
			throw 'No Join Field provided. Please ensure .GetlinkField provides both ' $
				'sides of the rename for the linked field'
			}
		// joinField should be the name from the SubTable that we are joining on
		semiJoin = ' semijoin by (' $ joinField $ ')'
		extraRenames = .extraRenames(joinField, linkedField)
		q = QueryHelper.StripSort(subTableQuery) $ extraRenames
		return semiJoin $ ' (' $ QueryHelper.AddWhere(q, where) $ ')'
		}

	// this does NOT handle when the subTable is a view, or if the linked field
	// is not a foreign key
	findJoinField(hdrQuery, subQuery, access)
		{
		hdrKeys = access.GetKeys()
		hdrTable = QueryGetTable(hdrQuery)
		subTable = QueryGetTable(subQuery)
		for fKey in QueryList("indexes where table is " $
			Display(subTable) $ " and fktable is " $ Display(hdrTable), "fkcolumns")
			{
			if false isnt hdrKeyMem = hdrKeys.FindIf({ it.Prefix?(fKey) })
				{
				return hdrKeys[hdrKeyMem]
				}
			}
		return false
		}

	extraRenames(joinField, linkField)
		{
		extraRenames = ''
		if linkField isnt joinField
			extraRenames $= ' rename ' $ linkField $ ' to ' $ joinField
		return extraRenames
		}

	Layout(access, block)
		{
		linkedBrowses = access.GetLinkedBrowseTabs()
		subTablesConfig = access.SubTablesConfig()

		// ensure that the SelectRepeats are added in the order
		// they are shown on the screen, not the order they were
		// added to linkedBrowses (which depends on the order the tabs were constructed)
		for idx in linkedBrowses.Members().Sort!()
			{
			linked = linkedBrowses[idx]

			// Might Need a way to exclude linkedBrowses
			// e.g. Recent Trucking Transactions on Business Partners
			linkedBrowse = linked.browse
			saveName = .subTableFilterName(linkedBrowse)

			// Need to filter the columns by the columns in the linkedBrowse Query
			selectCols = linkedBrowse.GetColumns().Copy()
			availableCols = .getAvailableColumns(linkedBrowse, selectCols,
				subTablesConfig)
			columns = selectCols.Intersect(availableCols)

			sf = access.SetSubtableSelectFields(columns, availableCols, saveName)
			// we cannot build these ahead of time as the tabs need to be constructed
			selectMgr = access.SetSubtableSelectMgr(sf.Fields, saveName)
			selectVals = selectMgr.Select_vals()

			control = Object('Expand', linked.name $ ' Filters',
				Object('SelectRepeat', sf, selectVals, linked.name),
					saveExpandName: saveName $ ' Filters')
			block(control, linked.name, saveName)
			}
		}

	getAvailableColumns(linkedBrowse, selectCols, subTablesConfig)
		{
		// we do NOT need to worry about the difference between GetQuery and
		// GetLinkedQuery here as we ONLY use the query to get the columns, not data
		colQuery = linkedBrowse.GetQuery()

		// Have an issue where sometimes the control we put on the SelectRepeat NEEDS
		// do be different that the control that exists in the List.
		// using this to override renames in the query so we get the correct
		// control
		overrides = Object()
		if Object?(subTablesConfig)
			if subTablesConfig.Member?(linkedBrowse.Name)
				overrides = subTablesConfig[linkedBrowse.Name]

		extraRename = ""
		if overrides.Member?('renames')
			{
			renames = overrides.renames
			for field in renames.Members()
				{
				if not selectCols.Has?(field)
					continue
				extraRename $= field $ " to " $ renames[field] $ ' '
				selectCols.Replace(field, renames[field])
				}
			}
		// Have an issue where sometimes the base query for the line items
		// has fields with duplicate select promts, this is normally fine if
		// one of the fields does NOT get used on the screen, but we now need
		// to run SelectFields on this query, so we need to exclude the duplicate field
		// make sure to exclude one that DOES NOT show on the screen
		extraRemoves = ""
		if overrides.Member?('removes')
			{
			extraRemoves = " remove " $ overrides.removes.Join(', ')
			}
		colQuery =  QueryHelper.StripSort(colQuery)
		if extraRename isnt ""
			colQuery =  colQuery $ ' rename ' $ extraRename
		colQuery $= extraRemoves

		// this will NOT have an effect on the fields that the user can select
		// this is ONLY to handle going from a Dynamic type that has the field renamed
		// to a dynamic type that does NOT have the field extended
		// this is JUST for the where on the semijoin
		return QueryColumns(colQuery)
		}

	}