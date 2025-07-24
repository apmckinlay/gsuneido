// Copyright (C) 2023 Suneido Software Corp. All rights reserved worldwide.
// NOTE: Use this class where ChooseColor needs to be called
class
	{
	CallClass(initVal, hwnd, custColors = false)
		{
		ob = Object(size: CHOOSECOLOR.Size(), flags: CC.FULLOPEN | CC.RGBINIT,
			owner: hwnd, custColors: custColors is false ? Object() : custColors,
				rgbResult: initVal)

		if false isnt Dialog.DoWithWindowsDisabled({ ChooseColor(ob) })
			return ob.rgbResult

		return false
		}
	}
