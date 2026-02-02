// Copyright (C) 2022 Axon Development Corporation All rights reserved worldwide.
_CanvasItem
	{
	id: false
	Getter_Id()
		{
		if .id isnt false
			return .id
		return .id = SuRenderBackend().NextId()
		}
	}