package controller

import (
	"github.com/jinzhu/gorm"
	"github.com/simon-ding/lklk/models"
	"github.com/simon-ding/lklk/yyets"
)

type ImportYYetsAPI struct {
	Params struct{
		Username string `json:"username"`
		Password string `json:"password"`
	} `http:"json"`
	DB *gorm.DB
}

func (i *ImportYYetsAPI) Handle() (*Response, error) {

	c := yyets.Client{}
	c.SetLogin(i.Params.Username, i.Params.Password)
	favs, err := c.UserFavs()
	if err != nil {
		return nil, err
	}
	for _, f := range favs {
		tv := models.TVSeries{
			YyetsID:  f,
		}
		i.DB.Create(&tv)
	}

	return SuccessReturn(), nil
}