// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
/*
TabCalcs acts as an interface class for its two children classes:
	- TabHorzCalcs: Used for horizontal tabs
	- TabVertCalcs: Used for vertical tabs

The children classes specify which methods are required in order to properly
position and manage the tabs. The children classes define the following public
members for use by the parent. Each member references the appropriate method given the
"orientation" argument.
	- RenderRect: Calculates the general rect required for the tab
	- DrawSpecs: Calculates the various rects / values required to draw the tabs
	- TextPos: Calculates the x/y coordinates for the tab text
	- ImagePos: Calculates the x/y coordinates for the tab image
	- LinePoints: Provides the x/y coordinates for drawing the tab border line
	- FontOrientation: Denotes the rotation for the tab's text
	- SelectSize: Adjusts the tabs render rect height/width based on the selected tab

The children classes also define various public methods for direct calls.
Barring special cases (IE: TabVertCalcs.ImageRect), the "orientation" argument does not
affect these methods.
*/
class
	{
	PaddingTop: 		5
	PaddingSide: 		7
	New(.controller, .selectedWeight = '', .trimChar = '.', .trimChars = 3)
		{
		.initFont()
		.imageSize = .controller.Ymin
		.TabHeight = .controller.Ymin += 2 * .PaddingTop
		.ButtonSize = .controller.Ymin - ScaleWithDpiFactor(6 /*= offset*/)
		}

	initFont()
		{
		.controller.WithSelectObject(.Font(selectedTab?:))
			{|hdc|
			GetTextExtentPoint32(hdc, 'M', 1, ex = Object())
			.controller.Ymin = ex.y
			}
		}

	Font(selectedTab? = false)
		{
		weight = selectedTab? ? .selectedWeight : ''
		return .fonts.GetInit(weight, { .createFont(weight, .FontOrientation) })
		}

	getter_fonts()
		{
		return .fonts = Object()
		}

	createFont(weight, orientation)
		{
		return CreateFontIndirect(.controller.LogFont(:weight, :orientation))
		}

	Getter_Ymin()
		{
		return .controller.Ymin
		}

	W: false
	H: false
	Resize(.W, .H)
		{
		}

	CalcRenderRect(i, tab, prevEnd, tabName = false)
		{
		if tabName isnt false
			.calcTabNameMetrics(tab, tabName)
		tab.width = tab.textWidth + .ImageWidth(tab.image) + .PaddingSide * 2
		tab.height = tab.textHeight + .PaddingTop * 2
		tab.renderRect = (.RenderRect)(i, tab, Max(prevEnd - 1, 0))
		tab.renderWidth = tab.width
		return tab.renderRect.end
		}

	calcTabNameMetrics(tab, tabName)
		{
		.controller.WithSelectObject(.Font())
			{|hdc|
			baseMetrics = .getTextMetric(hdc, tab.tabName = tabName)
			}
		.controller.WithSelectObject(.Font(selectedTab?:))
			{|hdc|
			boldMetrics = .getTextMetric(hdc, tabName)
			}
		tab.textWidth = boldMetrics.x
		tab.textHeight = boldMetrics.y
		tab.textBoldOffset = (boldMetrics.x - baseMetrics.x) / 2
		tab.boldCharSizes = boldMetrics.charSizes
		tab.baseCharSizes = baseMetrics.charSizes
		tab.trimCharSize = boldMetrics.trimCharSize
		}

	getTextMetric(hdc, text)
		{
		charSizes = Object()
		for i in .. text.Size()
			{
			GetTextExtentPoint32(hdc, text[i], 1, charMetric = Object())
			charSizes[i] = charMetric.x
			}
		GetTextExtentPoint32(hdc, text, text.Size(), metrics = Object())
		GetTextExtentPoint32(hdc, .trimChar, 1, trimMetric = Object())
		return [x: metrics.x, y: metrics.y, :charSizes, trimCharSize: trimMetric.x]
		}

	SelectOffset(i)
		{
		return i is .controller.GetSelected() ? 1 : 4 /*= select offset value*/
		}

	CalcDrawSpecs(wLarge, hLarge)
		{
		return (.DrawSpecs)(wLarge, hLarge)
		}

	CalcImageSpecs(tab)
		{
		imageSpecs = (.ImagePos)(tab)
		imageSpecs.size = .imageSize
		return imageSpecs
		}

	CalcLinePoints()
		{
		return (.LinePoints)()
		}

	CalcSelectChange(i, tab)
		{
		(.SelectSize)(i, tab.renderRect)
		}

	CalcTextSpecs(tab, selectedTab)
		{
		if '' is text = .fitText(tab, selectedTab)
			return false
		specs = (.TextPos)(tab, selectedTab)
		specs.text = text
		return specs
		}

	fitText(tab, selectedTab)
		{
		if tab.renderWidth is tab.width
			return tab.tabName
		availableSpace = tab.renderWidth - .ImageWidth(tab.image)
		charSizes = selectedTab ? tab.boldCharSizes : tab.baseCharSizes
		return .trimText(tab, charSizes, availableSpace)
		}

	ImageWidth(image)
		{
		return image is false ? 0 : .imageSize + 4 /*= image padding right*/
		}

	trimText(tab, charSizes, availableSpace)
		{
		allocatedSpace = .PaddingSide
		count = 0
		calc = {
			|charSize|
			if availableSpace <= allocatedSpace + charSize
				break
			allocatedSpace += charSize
			count++
			}

		// Calculate the first character to ensure it is displayed as long as possible
		charSizes[.. 1].Each(calc)
		tabNameChars = count

		// Calculate how many trim characters fit given the remaining space
		count = 0
		.trimCharOb(tab.trimCharSize).Each(calc)
		trimChars = count

		// Calculate how many tab name characters fit given the remaining space
		count = 0
		charSizes[1 ..].Each(calc)
		tabNameChars += count

		return tab.tabName[:: tabNameChars] $ .trimChar.Repeat(trimChars)
		}

	trimCharOb(charSize)
		{
		return Object().AddMany!(charSize, .trimChars)
		}

	ImageDimensions(tab)
		{
		// Image dimensions are not effected by the tab being selected and vice versa
		.controller.WithSelectObject(.Font())
			{|hdc|
			width = tab.image[0].Width(hdc)
			height = tab.image[0].Height(hdc)
			}
		return [:width, :height]
		}

	Destroy()
		{
		.fonts.Each(DeleteObject)
		.fonts.Delete(all:)
		}
	}