package fakers

type user struct {
	Email        string `fakers:"email,unique"`
	PasswordHash string `fakers:"password"`
	Nickname     string `fakers:"username"`
	AvatarUrl    string `fakers:"url"`
	IsPrivate    bool
}

type anime struct {
	Title       string  `faker:"sentence"`
	ReleaseDate string  `faker:"date"`
	Rating      float32 `faker:"boundary_start=0.0, boundary_end=5.0"`
	Genre       string  `faker:"word"`
	Type        string  `faker:"oneof: SERIES, MOVIE, OVA"`
	Studio      string  `faker:"word"`
	Status      string  `faker:"oneof: ONGOING, PLANNED, RELEASED"`
	AgeRating   string  `faker:"oneof: G, PG, PG-13, R, R+, Rx"`
	CoverUrl    string  `faker:"url"`
}

type userAnimeStatus struct {
	AnimeId      int    `faker:"boundary_start=1, boundary_end=1000000"`
	UserId       int    `faker:"boundary_start=1, boundary_end=1000000"`
	List         string `faker:"oneof: WATCHING, WATCHED, WANT_TO_WATCH, DROPPED, DEFERRED, REVISING"`
	IsSubscribed bool
}

type AnimeSeries struct {
	AnimeId      int    `faker:"boundary_start=1, boundary_end=1000000"`
	SeriesUrl    string `faker:"url"`
	SecondsCount int
}

type Series struct {
	AnimeId int     `faker:"boundary_start=1, boundary_end=1000000"`
	UserId  int     `faker:"boundary_start=1, boundary_end=1000000"`
	Rating  float32 `faker:"boundary_start=0.0, boundary_end=5.0"`
	Text    string  `faker:"paragraph"`
}

type Friends struct {
	UserId1 int `faker:"boundary_start=1, boundary_end=1000000"`
	UserId2 int `faker:"boundary_start=1, boundary_end=1000000"`
}

type Messages struct {
	SenderId   int    `faker:"boundary_start=1, boundary_end=1000000"`
	ReceiverId int    `faker:"boundary_start=1, boundary_end=1000000"`
	Text       string `faker:"paragraph"`
	FileUrl    string `faker:"url"`
	Time       string `faker:"time"`
}

type Achievements struct {
	AnimeId     int    `faker:"boundary_start=1, boundary_end=1000000"`
	Name        string `faker:"first_name"`
	Description string `faker:"sentence"`
}

type UserAchievements struct {
	AchievementId int    `faker:"boundary_start=1, boundary_end=1000000"`
	UserId        int    `faker:"boundary_start=1, boundary_end=1000000"`
	Time          string `faker:"time"`
}

type Notifications struct {
	UserId int    `faker:"boundary_start=1, boundary_end=1000000"`
	Type   string `faker:"oneof: FRIEND_REQUEST, NEW_EPISODE, NEW_MESSAGE"`
	Body   string
	Time   string `faker:"time"`
}

type Reactions struct {
	ReviewId int `faker:"boundary_start=1, boundary_end=1000000"`
	UserId   int `faker:"boundary_start=1, boundary_end=1000000"`
	IsLike   bool
}
