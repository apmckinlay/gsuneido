// Copyright (C) 2018 Suneido Software Corp. All rights reserved worldwide.
class
	{
	New(svc, .lib, .name)
		{
		.svc = svc
		.data = Object()
		.lastCommit = Date.Begin()
		.GetData()
		}

	SetName(name)
		{
		.name = name
		.GetData()
		}

	GetData()
		{
		.data = Object()
		when = Timestamp()
		while ((list = .svc.Get10Before(.lib, .name, when)) isnt #()) {
			for item in list
				.data.Add(item.Project(#(lib_committed, id, comment)))
			when = .data.Last().lib_committed
			}
		if .data.NotEmpty?()
			.lastCommit = .data[0].lib_committed
		.data.Sort!({ |x,y| x.lib_committed < y.lib_committed })
		}

	GetContributors()
		{
		contributers = Object()
		for change in .data
			for id in change.id.Split(",")
				contributers.AddUnique(id.Trim())
		return contributers
		}

	GetNumOfContributors()
		{
		return .GetContributors().Size()
		}

	GetContribPercentages() // outputs % of all changes each user contributed to
		{
		totalChanges = .data.Size()
		userPercentages = .GetContributors().Map({ |user|
			Object(:user, percentage: .getPercentage(user, totalChanges).Round(2)) })
		return userPercentages
		}

	getPercentage(user, totalChanges)
		{
		n = .data.CountIf({ it.id.Split(",").Map(#Trim).Has?(user) })
		return n * 100 / totalChanges /*= percent*/
		}

	// Attributes a score to each contributor that represents how knowledgeable they are
	// about the current code in the chosen record based on their contributions
	WeighContributions()
		{
		points = Object().Set_default(0).ListToNamed(@.GetContributors())
		previousId = ""
		for item in .data
			{
			comment = item.comment.Lower()
			if comment.Has?("cosmetic")
				continue // don't weigh cosmetic changes

			basePoint = .getBasePoint(comment, item.id)
			weightByDate = .getWeightByDate(item.lib_committed)

			// weight is added if changes by a user are successive
			for id in item.id.Split(",").Map(#Trim)
				{
				adjustedBasePoint = item.id.Prefix?(id)
					? basePoint
					: basePoint - 0.5 /*=points taken off for being secondary partner*/
				consecutiveBonus = previousId.Has?(id) ? 1 : 0
				points[id] += (adjustedBasePoint + consecutiveBonus) * weightByDate
				}

			previousId = item.id
			}
		for contributor in .GetContributors()
			points[contributor] *= .getPercentMultiplier(contributor)
		points = points.Assocs()
		return points.Sort!({ |x,y| x[1] > y[1] })
		}

	getBasePoint(comment, id)
		{
		basePoint = 4 // comment not labeled as issue or minor refactor
		if comment.Prefix?("issue") or comment =~ `\(\d+\)\s*\Z`
			basePoint = 6 // changes for a suggestion are weighted the most
		else if comment.Prefix?("minor")
			basePoint = 2 // minor refactors are weighted the least

		return id.Count(",") >= 1
			? --basePoint // working with others slightly decreases weight
			: basePoint
		}

	getWeightByDate(lib_committed)
		{
		interval = 0.1
		monthDiff = (.lastCommit.MinusDays(lib_committed) / 30 /*=days in a month*/).
			RoundDown(0)

		return monthDiff >= 7 /*=month weighting cutoff*/
			? 0.4 /*=month weight cutoff*/
			: 1 - interval * monthDiff
		}

	getPercentMultiplier(id)
		{
		contribPercentage = .GetContribPercentages().FindOne({ it.user is id }).percentage
		return (1 + contribPercentage / 100 /*=percent to decimal*/)
		}
	}
