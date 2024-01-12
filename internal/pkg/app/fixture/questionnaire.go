package fixture

import (
	"meet/internal/pkg/app/model"
	"time"

	"github.com/guregu/null"
)

var (
	astana = null.IntFrom(65737)
	moscow = null.IntFrom(99972)
)

var queBaseModel = baseModelSeq()

var (
	ElaxerQue = &model.Questionnaire{
		BaseModel:      queBaseModel(),
		UserID:         Elaxer.ID,
		Name:           "Elaxer",
		BirthDate:      model.NewBirthDate(2000, time.March, 31),
		Gender:         model.GenderMale,
		Orientation:    model.OrientationHetero,
		MeetingPurpose: model.MeetingPurposeRelationship,
		AgeRange:       model.NewAgeRange(20, 30),
		CityID:         astana,
		About:          "Elaxer, 23 года, познакомлюсь с девушкой старше 20 и младше 30 лет",
		Photos:         []*model.Photo{},
		IsActive:       true,
	}
	MariyaQue = &model.Questionnaire{
		BaseModel:      queBaseModel(),
		UserID:         Mariya.ID,
		Name:           "Мария",
		BirthDate:      model.NewBirthDate(1990, time.June, 18),
		Gender:         model.GenderFemale,
		Orientation:    model.OrientationBi,
		MeetingPurpose: model.MeetingPurposeRelationship,
		AgeRange:       model.NewAgeRange(30, 40),
		CityID:         astana,
		About:          "Мария, 33 года, познакомлюсь с паренем или девушкой старше 30 и младше 40 лет",
		Photos:         []*model.Photo{},
		IsActive:       true,
	}
	ElenaQue = &model.Questionnaire{
		BaseModel:      queBaseModel(),
		UserID:         Elena.ID,
		Name:           "Елена Петрова",
		BirthDate:      model.NewBirthDate(2003, time.December, 1),
		Gender:         model.GenderFemale,
		Orientation:    model.OrientationHomo,
		MeetingPurpose: model.MeetingPurposeSex,
		AgeRange:       model.NewAgeRange(18, 25),
		CityID:         astana,
		About:          "Елена, 20 лет, познакомлюсь с девушкой старше 18 и младше 25 лет",
		Photos:         []*model.Photo{},
		IsActive:       true,
	}
	KristinaQue = &model.Questionnaire{
		BaseModel:      queBaseModel(),
		UserID:         Kristina.ID,
		Name:           "Кристина",
		BirthDate:      model.NewBirthDate(1998, time.October, 23),
		Gender:         model.GenderFemale,
		Orientation:    model.OrientationHetero,
		MeetingPurpose: model.MeetingPurposeSex,
		AgeRange:       model.NewAgeRange(18, 25),
		CityID:         moscow,
		About:          "Кристина, 26 лет, познакомлюсь с парнем старше 18 и младше 25 лет",
		Photos:         []*model.Photo{},
		IsActive:       true,
	}
	DmitriyQue = &model.Questionnaire{
		BaseModel:      queBaseModel(),
		UserID:         Dmitriy.ID,
		Name:           "Дима",
		BirthDate:      model.NewBirthDate(2000, time.September, 5),
		Gender:         model.GenderMale,
		Orientation:    model.OrientationHetero,
		MeetingPurpose: model.MeetingPurposeRelationship,
		AgeRange:       model.NewAgeRange(18, 23),
		CityID:         astana,
		About:          "Дмитрий, 23 года, познакомлюсь с девушкой старше 18 и младше 23 лет",
		Photos:         []*model.Photo{},
		IsActive:       true,
	}
	VasiliyQue = &model.Questionnaire{
		BaseModel:      queBaseModel(),
		UserID:         Vasiliy.ID,
		Name:           "Василий",
		BirthDate:      model.NewBirthDate(1993, time.April, 1),
		Gender:         model.GenderMale,
		Orientation:    model.OrientationBi,
		MeetingPurpose: model.MeetingPurposeFriendship,
		AgeRange:       model.NewAgeRange(20, 60),
		CityID:         astana,
		About:          "Василий, 30 лет, познакомлюсь с парнем или девушкой старше 20 и младше 60 лет",
		Photos:         []*model.Photo{},
		IsActive:       true,
	}
)

var Questionnaires = map[*model.User]*model.Questionnaire{
	Elaxer:   ElaxerQue,
	Mariya:   MariyaQue,
	Elena:    ElenaQue,
	Kristina: KristinaQue,
	Dmitriy:  DmitriyQue,
	Vasiliy:  VasiliyQue,
}
