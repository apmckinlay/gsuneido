// Copyright (C) 2024 Suneido Software Corp. All rights reserved worldwide.
/* WARNING: Read before modifying this record.
This record is used for the majority of IDE PAINT procedures.
As a result, errors thrown by this record can have catastrophic results for the IDE.
When making changes to this record, proceed with caution.

Usage:
- The settings argument has the following structure:
	Object(
		"SelectObjectValue_n",
		"SelectObjectValue_n + 1",
		...
		"SelectObjectValue_n + n",
		brush: 	<hex color>,
		<Win32BuiltinNames Set Method>: New Value)
- Example:
	Object(
		font, 			// Will be processed through: SelectObject
		brush:			CLR.BLUE,
		SetBkMode: 		OPAQUE,
		SetTextColor:	CLR.RED)

- Call Example:
	WithHdcSettings(hdc, Object(font, SetBkMode: TRANSPARENT))
		{
		<code to be executed with the specified settings>
		}
*/
class
	{
	CallClass(hdc, settings, block)
		{
		.process(hdc, settings, restoreCalls = Object())
		return Finally(block, { restoreCalls.Each({ it() }) })
		}

	process(hdc, settings, restoreCalls)
		{
		try
			{
			.processValues(hdc, settings, restoreCalls)
			.processMembers(hdc, settings, restoreCalls)
			}
		catch (e)
			SuneidoLog.OnceByCallstack('ERROR: WithHdcSettings threw: ' $ e)
		}

	processValues(hdc, settings, restoreCalls)
		{
		settings.Values(list:).Each({ .setCall(hdc, it, restoreCalls, SelectObject) })
		}

	processMembers(hdc, settings, restoreCalls)
		{
		settings.Members(named:).Each()
			{
			if it is #brush
				.setBrush(hdc, settings[it], restoreCalls)
			else
				.setCall(hdc, settings[it], restoreCalls, Global(it))
			}
		}

	setBrush(hdc, color, restoreCalls)
		{
		prevValue = SelectObject(hdc, brush = CreateSolidBrush(color))
		restore = {
			SelectObject(hdc, prevValue)
			DeleteObject(brush)
			}
		restoreCalls.Add(restore)
		}

	setCall(hdc, newValue, restoreCalls, callable)
		{
		prevValue = callable(hdc, newValue)
		restoreCalls.Add({ (callable)(hdc, prevValue) })
		}
	}
