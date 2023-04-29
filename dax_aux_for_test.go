package sabi

func ClearGlobalDaxSrcs() {
	isGlobalDaxSrcsFixed = false
	globalDaxSrcMap = make(map[string]DaxSrc)
}

func IsGlobalDaxSrcsFixed() bool {
	return isGlobalDaxSrcsFixed
}

func GlobalDaxSrcMap() map[string]DaxSrc {
	return globalDaxSrcMap
}

func IsLocalDaxSrcsFixed(base DaxBase) bool {
	return base.(*daxBaseImpl).isLocalDaxSrcsFixed
}

func LocalDaxSrcMap(base DaxBase) map[string]DaxSrc {
	return base.(*daxBaseImpl).localDaxSrcMap
}

func DaxConnMap(base DaxBase) map[string]DaxConn {
	return base.(*daxBaseImpl).daxConnMap
}

func Begin(base DaxBase) {
	base.begin()
}

func Commit(base DaxBase) Err {
	return base.commit()
}

func Rollback(base DaxBase) {
	base.rollback()
}

func End(base DaxBase) {
	base.end()
}
