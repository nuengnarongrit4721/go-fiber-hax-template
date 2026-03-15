package handlers

type Set struct {
	User *UserHandler
	Auth *AuthHandler
}

type VersionedSet struct {
	V1 Set
	V2 Set
}
