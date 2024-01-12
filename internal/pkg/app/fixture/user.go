package fixture

import (
	"meet/internal/pkg/app/model"

	"github.com/guregu/null"
)

func passwordHash(password string) null.String {
	h, err := model.Password(password).Hash()
	if err != nil {
		panic(err)
	}

	return null.StringFrom(h)
}

var userBaseModel = baseModelSeq()

var (
	Elaxer = &model.User{
		BaseModel:    userBaseModel(),
		Login:        "elaxer",
		PasswordHash: passwordHash("123456"),
	}
	Mariya = &model.User{
		BaseModel:    userBaseModel(),
		Login:        "mariya669",
		PasswordHash: passwordHash("qwerty"),
	}
	Elena = &model.User{
		BaseModel:    userBaseModel(),
		Login:        "elena420",
		PasswordHash: passwordHash("йцукен"),
	}
	Kristina = &model.User{
		BaseModel:    userBaseModel(),
		Login:        "kristina",
		PasswordHash: passwordHash("kristina"),
	}
	Dmitriy = &model.User{
		BaseModel:    userBaseModel(),
		Login:        "dmitriy",
		PasswordHash: passwordHash("322322"),
	}
	Vasiliy = &model.User{
		BaseModel:    userBaseModel(),
		Login:        "vasiliy",
		PasswordHash: passwordHash("hell666"),
	}
)

var Users = []*model.User{
	Elaxer, Mariya, Elena, Kristina, Dmitriy, Vasiliy,
}
