// Copyright (C) 2001 Suneido Software Corp. All rights reserved worldwide.
class
	{
	CallClass(func)
		{
		if not .Member?("F_" $ func)
			throw "Accumulator: invalid function: " $ func
		return new this["F_" $ func]
		}
	F_total: class
		{
		total: 0
		Reset()
			{ .total = 0 }
		Value(value)
			{ .total += value }
		Result()
			{ return .total }
		}
	F_average: class
		{
		total: 0
		count: 0
		Reset()
			{ .total = .count = 0 }
		Value(value)
			{ .total += value; ++.count }
		Result()
			{ return .total / .count }
		}
	F_max: class
		{
		New()
			{ .max = this }
		Reset()
			{ .max = this }
		Value(value)
			{ if (.max is this or value > .max) .max = value }
		Result()
			{ return .max }
		}
	F_min: class
		{
		New()
			{ .min = this }
		Reset()
			{ .min = this }
		Value(value)
			{ if (.min is this or value < .min) .min = value }
		Result()
			{ return .min }
		}
	F_count: class
		{
		count: 0
		Reset()
			{ .count = 0 }
		Value(unused)
			{ ++.count }
		Result()
			{ return .count }
		}
	}
