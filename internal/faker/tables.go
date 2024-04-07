package faker

type user struct {
	Email        string `faker:"email"`
	PasswordHash string `faker:"password"`
	Nickname     string `faker:"username"`
	AvatarUrl    string `faker:"url"`
	IsPrivate    bool
}
