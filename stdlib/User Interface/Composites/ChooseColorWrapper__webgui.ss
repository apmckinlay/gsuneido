// Copyright (C) 2023 Axon Development Corporation All rights reserved worldwide.
class
	{
	CallClass(initVal, hwnd /*unused*/, custColors /*unused*/)
		{
		ob = Object(rgbResult: initVal)
		return ChooseColor(ob)
			? ob.rgbResult
			: false
		}
	}
