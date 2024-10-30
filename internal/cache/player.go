package cache

var playerCache = &syncedMap[int64, int]{}

func InsertPlayer(dbId int, pdgaNum int64) {
    playerCache.Set(pdgaNum, dbId)
}

func PlayerDbId(pdgaNum int64) int {
   return playerCache.Get(pdgaNum)
}
