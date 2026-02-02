// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
function (key, pressed = false)
	{
	if pressed is false
		return false

	if key is VK.SHIFT
		return pressed.GetDefault('shift', false)

	if key is VK.CONTROL
		return pressed.GetDefault('control', false)

	if key is VK.MENU
		return pressed.GetDefault('alt', false)

	throw 'KeyPressed? not handling key ' $ String(key)
	}