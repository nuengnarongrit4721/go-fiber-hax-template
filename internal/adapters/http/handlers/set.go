package handlers

type Set struct {
	User *UserHandler
}

type VersionedSet struct {
	V1 Set
	V2 Set
}
