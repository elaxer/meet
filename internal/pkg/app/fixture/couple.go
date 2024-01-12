package fixture

import "meet/internal/pkg/app/model"

var Couples = map[*model.User][]*model.User{
	Elaxer:  {Mariya},
	Dmitriy: {Kristina},
	Vasiliy: {Mariya},
}
