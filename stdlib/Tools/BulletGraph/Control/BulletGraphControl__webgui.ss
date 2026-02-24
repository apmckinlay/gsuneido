// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
Control
	{
	ComponentName: 'BulletGraphComponent'
	New(data, satisfactory = 0, good = 0, target = false, range = #(0, 100),
		color = 0x506363, width = 128, height = 32, rectangle = true,
		outside = 5, vertical = false, axis = false, axisDensity = 5,
		axisFormat = false, selectedColor = false)
		{
		if vertical and width is 128 and height is 32 /*= default vertical sizing*/
			{ // swap w and h
			temp = width
			width = height
			height = temp
			}
		.validateData(data, range, satisfactory, good, target)
		data, satisfactory, good, target, range =
			.normalization(data, satisfactory, good, target, range.Copy())
		.ComponentArgs = Object(data, :satisfactory, :good, :target,
			:range, :color, :width, :height, :rectangle, :outside, :vertical,
			:axis, :axisDensity, :axisFormat, :selectedColor)
		}

	validateData(data, range, satisfactory, good, target)
		{
		Assert(data < range[1])
		Assert(satisfactory >= range[0] and satisfactory <= range[1])
		Assert(good >= range[0] and good <= range[1])
		Assert(range[0] < range[1])
		if target isnt false
			Assert(target >= range[0] and target <= range[1])
		}

	normalization(data, satisfactory, good, target, range)
		{
		data -= range[0]
		range[1] -= range[0]
		satisfactory -= range[0]
		good -= range[0]
		if target isnt false
			target -= range[0]
		return data, satisfactory, good, target, range
		}

	MOUSEMOVE()
		{
		.Send('BulletGraph_Hover')
		return false
		}

	LBUTTONUP()
		{
		.Send('BulletGraph_Click')
		return false
		}

	Selected(selected)
		{
		.Act('Selected', selected)
		}
	}
