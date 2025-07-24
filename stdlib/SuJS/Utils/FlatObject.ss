// Copyright (C) 2021 Axon Development Corporation All rights reserved worldwide.
// returns an object whose maximum nesting level is less or equal to maxlevel
class
	{
	CallClass(ob, maxLevel)
		{
		_info = Object(extra: Object(), :maxLevel, curObs: Object(), id: 0)
		newOb = Object()
		.flat(ob, newOb, 1) 	// (maxLevel - 1) levels
										// +
										// 1 (extra level from result.ob)
		return Object(ob: newOb, extra: _info.extra)
		}

	flat(ob, newOb, level)
		{
		for m in ob.Members()
			{
			if not Object?(ob[m])
				{
				newOb[m] = ob[m]
				continue
				}
			.checkCirculr?(ob[m])
			_info.curObs.Add(ob[m])
			if level < _info.maxLevel
				{
				newOb[m] = Object()
				.flat(ob[m], newOb[m], level + 1)
				}
			else
				{
				id = .getUniqueId()
				newOb[m] = id
				item = _info.extra[id] = Object()
				.flat(ob[m], item, 2) // (maxLevel - 2) levels
												// +
												// 2 (extra level from result.extra[ts])
				}
			_info.curObs.PopLast()
			}
		}

	checkCirculr?(ob)
		{
		for c in _info.curObs
			if Same?(c, ob)
				throw 'found circular object'
		}

	getUniqueId()
		{
		return 'FO_#' $ _info.id++ $ '#'
		}

	Build(ob)
		{
		if not Object?(ob) or not ob.Member?(#extra) or not ob.Member?(#ob)
			return ob
		extra = ob.extra
		ob = ob.ob
		if extra.Empty?()
			return ob
		.build(ob, extra)
		return ob
		}

	build(ob, extra)
		{
		for m in ob.Members()
			{
			if String?(ob[m]) and extra.Member?(ob[m])
				ob[m] = extra[ob[m]]
			if Object?(ob[m])
				.build(ob[m], extra)
			}
		}
	}
