package handlers

type Set struct {
	User *UserHandler
	Auth *AuthHandler
	JWKS *JWKSHandler
}

type VersionedSet struct {
	V1 Set
	V2 Set
}
