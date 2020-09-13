package cache

func getGuildUnion(ID string) cacheGuild {
	Cache.GuildsRWMutex.RLock()
	defer Cache.GuildsRWMutex.RUnlock()
	var cg = Cache.Guilds[ID]
	// defer is run here
	return cg
}
