// Copyright (C) 2020 Suneido Software Corp. All rights reserved worldwide.
/* USAGE:
- To define custom markers / indicators for addons:
1. Add this method (or add to this method) in the addon:
	Styling()
		{
		return [
			[level: ###,
				marker: [... marker details here ...],
				indicator: [... indicator details here ...]],
			[level: ###,
				marker: [... marker details here ...],
				indicator: [... indicator details here ...]],
			...
			]
		}
===============  MARKERS  ================
- Parameters are as follows:
-- Required:
--- 0: The type of marker to be used.
---- XPM Example: stdlib:Addon_show_modified_lines
---- SC  Example: stdlib:Addon_check_code

-- Optional:
--- fore: 	Foreground color, (ie: CLR.RED)
--- back:	Background color, (ie: CLR.RED)
--- type: 	Allows for separation of markers (ie: OverviewBar markers)

- NOTE: Type and level need to be a unique combination, or the marker will not be defined
- NOTE: If the number of markers exceed MAXMARKERS, the higher level markers will be
		given precedence
==========================================

===============	INDICATORS ===============
- Indicator parameters are as follows:
-- Required: 0:		The type of indicator to be used (Reference: stdlib:INDIC)
-- Optional: fore: 	Foreground color, (ie: CLR.RED)
==========================================

- If a reference to the marker / indicator is required, call:
-- MarkerIdx(level, type)
-- IndicatorIdx(level)
--- The return value can be used with Scintilla to Add / Remove the marker where required
*/
class
	{
	New(.scintilla)
		{
		.markerColors = []
		.markers = Object().Set_default(Object())
		.indicators = []
		}

	DefineStyles(addons)
		{
		markers = Object().Set_default(Object())
		indicators = Object().Set_default(Object())
		for styleOb in addons.Collect(#Styling).Add(.scintilla.BaseStyling())
			for style in styleOb
				{
				level = style.GetDefault(#level, 0)
				if not style.GetDefault(#marker, Object()).Empty?()
					.processMarker(style.marker, level, markers)
				if not style.GetDefault(#indicator, Object()).Empty?()
					.processIndicator(style.indicator, level, indicators)
				}
		.defineIndicators(indicators.Sort!(By(#level)))
		.defineMarkers(markers.Sort!(By(#level)))
		}

	// MARKERS
	MAXMARKERS: 24 // Limited based on Scinitlla Documentation
	defineMarkers(markers)
		{
		markers.Sort!({ |x, y|  x.level < y.level })
		if .MAXMARKERS < size = markers.Size()
			{
			remove = (.MAXMARKERS - size).Abs()
			SuneidoLog('ERROR: (CAUGHT) Too many markers defined: ' $ markers.Size(),
				params: markers[.. remove])
			markers = markers[remove ..]
			}
		for marker in markers
			.marker(@marker)
		}

	processMarker(marker, level, sortedMarkers)
		{
		// Ensures that any missing values are set
		marker = Object(fore: false, back: false, :level, type: #default).Merge(marker)
		if sortedMarkers.Any?({ it.level is marker.level and it.type is marker.type })
			.logDuplicateLevels(#marker, [level: marker.level, type: marker.type])
		else
			sortedMarkers.Add(marker)
		}

	logDuplicateLevels(style, params)
		{
		SuneidoLog('ERROR: (CAUGHT) Two ' $ style $ 's share the same type and level. ' $
			'Skipping ' $ style, :params)
		}

	processIndicator(indicator, level, sortedIndicators)
		{
		indicator = Object(fore: false, :level).Merge(indicator)
		if sortedIndicators.Any?({ it.level is indicator.level })
			.logDuplicateLevels(#indicator, [:level, type: indicator.level])
		else
			sortedIndicators.Add(indicator)
		}

	idx: 0
	marker(marker, fore, back, level, type)
		{
		if Number?(marker)
			.scintilla.DefineMarker(.idx, marker, fore, back)
		else
			.scintilla.DefineXPMMarker(.idx, marker, fore, back)
		.markerColors[.idx] = back
		.addMarker(.idx, level, type)
		.idx++
		}

	addMarker(idx, level, type)
		{
		markerGroup = .markers[type]
		for(i = 0; i < markerGroup.Size(); i++)
			if markerGroup[i].level <= level
				break
		markerGroup.Add(Object(:idx, :level), at: i)
		}

	MarkerIdx(level, type = #default)
		{ return .findIdx(.markers[type], level) }

	findIdx(stylesOb, level)
		{
		style = stylesOb.FindOne({ it.level is level })
		return style is false ? false : style.idx
		}

	ForEachMarkerByLevel(type, block)
		{
		for item in .markers[type]
			if block(item.idx) is true
				break
		}

	Getter_MarkerTypes()
		{ return .markers.Members() }

	MarkerColor(i)
		{ return .markerColors[i] }

	// INDICATORS
	nextIndicator: 8
	defineIndicators(indicators)
		{
		for indicator in indicators
			.indicator(@indicator)
		}

	indicator(style, fore, level)
		{
		.scintilla.DefineIndicator(.nextIndicator, style, fore)
		.indicators.Add([idx: .nextIndicator++, level: level])
		}

	IndicatorIdx(level)
		{ return .findIdx(.indicators, level) }

	IndicatorAtPos?(pos)
		{
		for indicator in .indicators
			if .scintilla.HasIndicator?(pos, indicator.idx)
				return true
		return false
		}
	}