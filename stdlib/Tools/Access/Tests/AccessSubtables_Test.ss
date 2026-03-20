// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Test
	{
	sf: class
		{
		New(.field, .prompt)
			{
			}

		PromptToField(prompt)
			{
			return prompt is .prompt ? .field : 'unknown_field'
			}

		FieldToPrompt(field)
			{
			return field is .field ? .prompt : 'unknown_prompt'
			}
		}

	selectMgr: class
		{
		New(vals = false)
			{
			.select_vals = vals is false ? Object() : vals
			}
		Select_vals()
			{
			return .select_vals
			}
		SetSelectVals(.select_vals, sf/*unused*/)
			{
			}
		}

	makeCondition(field, check, operation = "equals", value = "test", value2 = '')
		{
		cond = Record()
		cond.condition_field = field
		cond[field] = Object(:operation, :value, :value2)
		cond.check = check
		return cond
		}

	Test_updateUnchecked()
		{
		fn = AccessSubtables.AccessSubtables_updateUnchecked

		// Test 1: sets check false when matching unchecked
		csf = .sf("field_a", "Prompt A")
		ssf = .sf("field_a", "Prompt A")

		sCondition = .makeCondition("field_a", true)
		currentConditions = Object(.makeCondition("field_a", false))

		fn(sCondition, currentConditions, csf, ssf)
		Assert(sCondition.check is false,
			msg: "Test 1: sCondition.check should be set to false")

		// Test 2: does not change sCondition when no match
		csf = .sf("field_a", "Prompt A")
		ssf = .sf("field_a", "Prompt A")

		sCondition = .makeCondition("field_a", true)
		currentConditions = Object(.makeCondition("field_b", false))

		fn(sCondition, currentConditions, csf, ssf)
		Assert(sCondition.check is false,
			msg: "Test 2: sCondition.check should get unchecked when there is no match")

		// Test 3: skips when sCondition already false
		csf = .sf("field_a", "Prompt A")
		ssf = .sf("field_a", "Prompt A")

		sCondition = .makeCondition("field_a", false)
		currentConditions = Object(.makeCondition("field_a", false))

		fn(sCondition, currentConditions, csf, ssf)
		Assert(sCondition.check is false,
			msg: "Test 3: sCondition.check should remain false")

		// Test 4: with renamed fields
		csf = .sf("field_y", "Common Prompt")
		ssf = .sf("field_x", "Common Prompt")

		sCondition = .makeCondition("field_x", true)
		currentConditions = Object(.makeCondition("field_y", false))

		fn(sCondition, currentConditions, csf, ssf)
		Assert(sCondition.check is false,
			msg: "Test 4: sCondition.check should be false when renamed fields match")

		// Test 5: ignores checked cConditions
		csf = .sf("field_a", "Prompt A")
		ssf = .sf("field_a", "Prompt A")

		sCondition = .makeCondition("field_a", true)
		currentConditions = Object(.makeCondition("field_a", true))

		fn(sCondition, currentConditions, csf, ssf)
		Assert(sCondition.check is true,
			msg: "Test 5: sCondition.check should remain true when cCondition is checked")

		// Test 6: with empty currentConditions
		csf = .sf("field_a", "Prompt A")
		ssf = .sf("field_a", "Prompt A")

		sCondition = .makeCondition("field_a", true)
		currentConditions = Object()

		fn(sCondition, currentConditions, csf, ssf)
		Assert(sCondition.check isnt true,
			msg: "Test 6: sCondition.check should not remain true " $
				"with empty currentConditions")
		}

	Test_updateChecked()
		{
		fn = AccessSubtables.AccessSubtables_updateChecked

		// Test 1: sets cCondition.check to true when matching sCondition is checked
		csf = .sf("field_a", "Prompt A")
		ssf = .sf("field_a", "Prompt A")

		sCondition = .makeCondition("field_a", true)
		cCondition = .makeCondition("field_a", false)
		currentConditions = Object(cCondition)
		selectMgr = .selectMgr()

		fn(sCondition, currentConditions, csf, ssf, selectMgr)
		Assert(cCondition.check is true,
			msg: "Test 1: cCondition.check should be set to true")

		// Test 2: adds to selectMgr when no matching cCondition found
		csf = .sf("field_a", "Prompt A")
		ssf = .sf("field_a", "Prompt A")

		sCondition = .makeCondition("field_a", true)
		currentConditions = Object(.makeCondition("field_b", false))
		selectMgr = .selectMgr()

		fn(sCondition, currentConditions, csf, ssf, selectMgr)
		Assert(selectMgr.Select_vals().Size() is 1,
			msg: "Test 2: condition should be added to selectMgr when no match")
		added = selectMgr.Select_vals()[0]
		Assert(added.condition_field is "field_a",
			msg: "Test 2: added condition should have correct field")
		Assert(added.check is true,
			msg: "Test 2: added condition should have check true")

		// Test 3: does nothing when sCondition.check is false
		csf = .sf("field_a", "Prompt A")
		ssf = .sf("field_a", "Prompt A")

		sCondition = .makeCondition("field_a", false)
		cCondition = .makeCondition("field_a", false)
		currentConditions = Object(cCondition)
		selectMgr = .selectMgr()

		fn(sCondition, currentConditions, csf, ssf, selectMgr)
		Assert(cCondition.check is false,
			msg: "Test 3: cCondition.check should remain false " $
				"when sCondition is unchecked")
		Assert(selectMgr.Select_vals().Size() is 0,
			msg: "Test 3: nothing should be added when sCondition is unchecked")

		// Test 4: with renamed fields, sets cCondition.check to true
		csf = .sf("field_y", "Common Prompt")
		ssf = .sf("field_x", "Common Prompt")

		sCondition = .makeCondition("field_x", true)
		cCondition = .makeCondition("field_y", false)
		currentConditions = Object(cCondition)
		selectMgr = .selectMgr()

		fn(sCondition, currentConditions, csf, ssf, selectMgr)
		Assert(cCondition.check is true,
			msg: "Test 4: cCondition.check should be true with renamed fields match")

		// Test 5: with renamed fields, adds to selectMgr when no match
		csf = .sf("field_y", "Common Prompt")
		ssf = .sf("field_x", "Common Prompt")

		sCondition = .makeCondition("field_x", true)
		currentConditions = Object(.makeCondition("other_field", false))
		selectMgr = .selectMgr()

		fn(sCondition, currentConditions, csf, ssf, selectMgr)
		Assert(selectMgr.Select_vals().Size() is 1,
			msg: "Test 5: renamed condition should be added when no match")
		added = selectMgr.Select_vals()[0]
		Assert(added.condition_field is "field_y",
			msg: "Test 5: added condition should have renamed field")
		Assert(added.check is true,
			msg: "Test 5: added condition should have check true")

		// Test 6: with empty currentConditions, adds to selectMgr
		csf = .sf("field_a", "Prompt A")
		ssf = .sf("field_a", "Prompt A")

		sCondition = .makeCondition("field_a", true)
		currentConditions = Object()
		selectMgr = .selectMgr()

		fn(sCondition, currentConditions, csf, ssf, selectMgr)
		Assert(selectMgr.Select_vals().Size() is 1,
			msg: "Test 6: condition should be added with empty currentConditions")
		added = selectMgr.Select_vals()[0]
		Assert(added.condition_field is "field_a",
			msg: "Test 6: added condition should have correct field")
		Assert(added.check is true,
			msg: "Test 6: added condition should have check true")

		// Test 7: does not modify sCondition
		csf = .sf("field_a", "Prompt A")
		ssf = .sf("field_a", "Prompt A")

		sCondition = .makeCondition("field_a", true)
		cCondition = .makeCondition("field_a", false)
		currentConditions = Object(cCondition)
		selectMgr = .selectMgr()

		fn(sCondition, currentConditions, csf, ssf, selectMgr)
		Assert(sCondition.check is true,
			msg: "Test 7: sCondition.check should remain unchanged")
		}

	cl: AccessSubtables
		{
		New(dynamic?, .testSelectMgrs, .testSfs)
			{
			super(dynamic?)
			}
		AccessSubtables_initSelectMgr(saveName, fields/*unused*/)
			{
			return .testSelectMgrs[saveName]
			}
		AccessSubtables_initSelectFields(columns/*unused*/, saveName)
			{
			return .testSfs[saveName]
			}
		}

	Test_integration()
		{
		// Setup stub select fields for multiple subtables
		sf1 = .sf("field_a", "Prompt A")
		sf2 = .sf("field_b", "Prompt B")

		// Setup stub select managers with pre-existing saved conditions
		cond1 = .makeCondition("field_a", true)
		cond2 = .makeCondition("field_b", true)

		sm1 = .selectMgr(Object(cond1))
		sm2 = .selectMgr(Object(cond2))

		// Create AccessSubtables instance (non-dynamic)
		testSelectMgrs = Object(
			"subtable1": sm1,
			"subtable2": sm2
			)
		testSfs = Object(
			"subtable1": sf1,
			"subtable2": sf2
			)
		ast = .cl(false, testSelectMgrs, testSfs)

		// Step 1: SetSubtableSelectFields for subtable1
		ast.SetSubtableSelectFields(#("field_a"), "subtable1")
		Assert(ast.AccessSubtables_linkedSelectFields.Member?("subtable1"),
			msg: "subtable1 should be in linkedSelectFields")

		// Step 2: SetSubtableSelectMgr for subtable1
		result = ast.SetSubtableSelectMgr(#("field_a"), "subtable1")
		Assert(result is sm1,
			msg: "SetSubtableSelectMgr should return the selectMgr")

		// Step 3: SetSubtableSelectFields for subtable2
		ast.SetSubtableSelectFields(#("field_b"), "subtable2")
		Assert(ast.AccessSubtables_linkedSelectFields.Member?("subtable2"),
			msg: "subtable2 should be in linkedSelectFields")

		// Step 4: SetSubtableSelectMgr for subtable2
		ast.SetSubtableSelectMgr(#("field_b"), "subtable2")

		// Step 5: SetSubTableSelectVals for subtable1 with new conditions
		newConditions = Object(.makeCondition("field_a", false))
		ast.SetSubTableSelectVals("subtable1", newConditions)

		// Verify the saved condition was unchecked
		Assert(sm1.Select_vals()[0].check is false,
			msg: "saved condition should be unchecked after SetSubTableSelectVals")

		// Verify subtable2's selectMgr was not affected
		Assert(sm2.Select_vals()[0].check is true,
			msg: "subtable2's condition should remain unchanged")
		}

	Test_integration_dynamic()
		{
		// Setup stub select fields - same prompt means renamed fields
		sf1 = .sf("field_a", "Common Prompt")
		sf2 = .sf("field_x", "Common Prompt")

		// Setup stub select manager with a checked saved condition
		savedCond = .makeCondition("field_a", true)
		sm1 = .selectMgr(Object(savedCond))

		// Setup empty select manager for second subtable
		sm2 = .selectMgr()

		// Create AccessSubtables instance (dynamic mode)
		testSelectMgrs = Object(
			"table|linkedBrowse": sm1,
			"table2|linkedBrowse": sm2
			)
		testSfs = Object(
			"table|linkedBrowse": sf1,
			"table2|linkedBrowse": sf2
			)
		ast = .cl(true, testSelectMgrs, testSfs)

		// Step 1: Initialize first subtable
		ast.SetSubtableSelectFields(#("field_a"), "table|linkedBrowse")
		ast.SetSubtableSelectMgr(#("field_a"), "table|linkedBrowse")

		// Step 2: Initialize second subtable (different field, same prompt)
		// This triggers updateChecked via handleDynamic
		ast.SetSubtableSelectFields(#("field_x"), "table2|linkedBrowse")
		ast.SetSubtableSelectMgr(#("field_x"), "table2|linkedBrowse")

		// Verify the condition was added to sm2 (renamed from field_a to field_x)
		Assert(sm2.Select_vals().Size() is 1,
			msg: "condition should be added to second selectMgr in dynamic mode")
		added = sm2.Select_vals()[0]
		Assert(added.condition_field is "field_x",
			msg: "added condition should have renamed field")
		Assert(added.check is true,
			msg: "added condition should have check true")

		// Step 3: SetSubTableSelectVals on second subtable
		// This triggers updateUnchecked
		newConditions = Object(.makeCondition("field_x", false))
		ast.SetSubTableSelectVals("table2|linkedBrowse", newConditions)

		// Verify the saved condition in sm1 was unchecked
		Assert(sm1.Select_vals()[0].check is false,
			msg: "saved condition should be unchecked after SetSubTableSelectVals")
		}

	Test_integration_multipleSubtables()
		{
		// Setup stub select fields for three subtables
		sf1 = .sf("field_a", "Prompt A")
		sf2 = .sf("field_b", "Prompt B")
		sf3 = .sf("field_c", "Prompt C")

		// Setup stub select managers with various conditions
		sm1 = .selectMgr(Object(.makeCondition("field_a", true)))
		sm2 = .selectMgr(Object(.makeCondition("field_b", true)))
		sm3 = .selectMgr(Object(.makeCondition("field_c", false)))

		testSelectMgrs = Object(
			"subtable1": sm1,
			"subtable2": sm2,
			"subtable3": sm3
			)
		testSfs = Object(
			"subtable1": sf1,
			"subtable2": sf2,
			"subtable3": sf3
			)
		ast = .cl(false, testSelectMgrs, testSfs)

		// Initialize all three subtables
		ast.SetSubtableSelectFields(#("field_a"), "subtable1")
		ast.SetSubtableSelectMgr(#("field_a"), "subtable1")

		ast.SetSubtableSelectFields(#("field_b"), "subtable2")
		ast.SetSubtableSelectMgr(#("field_b"), "subtable2")

		ast.SetSubtableSelectFields(#("field_c"), "subtable3")
		ast.SetSubtableSelectMgr(#("field_c"), "subtable3")

		// Verify all subtables are registered
		Assert(ast.AccessSubtables_subTableMgrs.Size() is 3,
			msg: "should have 3 subtables registered")
		Assert(ast.AccessSubtables_linkedSelectFields.Size() is 3,
			msg: "should have 3 select fields registered")

		// SetSubTableSelectVals on subtable2 with unchecked condition
		newConditions = Object(.makeCondition("field_b", false))
		ast.SetSubTableSelectVals("subtable2", newConditions)

		// Verify only subtable2's condition was affected
		Assert(sm1.Select_vals()[0].check is true,
			msg: "subtable1's condition should remain checked")
		Assert(sm2.Select_vals()[0].check is false,
			msg: "subtable2's condition should be unchecked")
		Assert(sm3.Select_vals()[0].check is false,
			msg: "subtable3's condition should remain unchecked")
		}

	}
