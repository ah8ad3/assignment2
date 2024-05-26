package main

func bootstrapUserQuotas() map[UserID]Quota {
	userQuotas := make(map[UserID]Quota, 2)
	userQuotas[1] = NewQuota(1, 1000, 1)
	userQuotas[2] = NewQuota(2, 2000, 20)

	return userQuotas
}
