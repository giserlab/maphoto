package jwt

import "maphoto/internal/env"

var (
	Env env.Options
)

func InitJWT(env *env.Options) {
	Env = *env
}
