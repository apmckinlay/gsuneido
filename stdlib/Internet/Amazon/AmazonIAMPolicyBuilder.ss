// Copyright (C) 2026 Suneido Software Corp. All rights reserved worldwide.
class
	{
	BuildPolicy(statements)
		{
		return Object(Version: '2012-10-17', Statement: statements)
		}

	BuildStatement(sid, effect, principal, actions, resource, conditions = #())
		{
		return Object(
			Sid: sid,
			Effect: effect,
			Principal: principal,
			Action: actions,
			Resource: resource,
			Condition: conditions
			)
		}

	FederatedUserCondition(accountId)
		{
		return Object(ArnLike: Object('aws:PrincipalArn':
			"arn:aws:sts::" $ accountId $ ":federated-user/*"))
		}

	// Convenience: builds an Allow statement that always includes
	// the federated-user ArnLike condition
	AllowStatement(sid, actions, resource, accountId, extraConditions = #())
		{
		conditions = .FederatedUserCondition(accountId).
			Merge(extraConditions)
		return .BuildStatement(sid, 'Allow', '*', actions,
			resource, conditions)
		}
	}
