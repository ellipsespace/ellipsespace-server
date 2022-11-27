package catalogueobject

import (
	"fmt"
	"strings"

	"github.com/qwuiemme/ellipsespace-server/pkg/client"
)

type CatalogueObject struct {
	// Название спутника
	Name string `json:"name"`
	// Подробное описание
	Description string `json:"description"`
	// Дата обнаружения EllipseSpace
	OpeningDateTime string `json:"o-date-time"`
	// Сидерический период обращения
	SidericConversionPeriod float32 `json:"s-conversion-period"`
	// Орбитальная скорость
	BodyOrbitalVelocity float32 `json:"orbital-vel"`
	// Наклонение
	Inclination float32 `json:"inclination"`
	// Спутники
	Satelites []string `json:"satelites"`
	// Чей спутник
	WhoseSatelite string `json:"whose-satelite"`
	// Экваториальный радиус
	EquatorialRadius float32 `json:"equator-radius"`
	// Полярный радиус
	PolarRadius float32 `json:"polar-radius"`
	// Средний радиус
	AverageRadius float32 `json:"avg-radius"`
	// Площадь
	Square float64 `json:"s"`
	// Объем
	Volume float64 `json:"v"`
	// Масса
	Weight float64 `json:"m"`
	// Средняя плотность
	AverageDensity float32 `json:"p"`
	// Ускорение свободного падения
	GravityAcceleration float32 `json:"g"`
	// Первая космическая скорость
	FirstSpaceVelocity float32 `json:"v1"`
	// Вторая космическая скорость
	SecondSpaceVelocity float32 `json:"v2"`
	//Фотографии
	Photos []string `json:"photos"`
}

func (c *CatalogueObject) Save() {

}

func (c *CatalogueObject) AddToDatabase() error {
	conn := client.Connect()
	defer conn.Close()

	res, err := conn.Query(fmt.Sprintf("INSERT INTO `catalogue` VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s')",
		c.Name,
		c.Description,
		c.OpeningDateTime,
		fmt.Sprintf("%g", c.SidericConversionPeriod),
		fmt.Sprintf("%g", c.BodyOrbitalVelocity),
		fmt.Sprintf("%g", c.Inclination),
		strings.Join(c.Satelites, "\n"),
		c.WhoseSatelite,
		fmt.Sprintf("%g", c.EquatorialRadius),
		fmt.Sprintf("%g", c.PolarRadius),
		fmt.Sprintf("%g", c.AverageRadius),
		fmt.Sprintf("%e", c.Square),
		fmt.Sprintf("%e", c.Volume),
		fmt.Sprintf("%e", c.Weight),
		fmt.Sprintf("%g", c.AverageDensity),
		fmt.Sprintf("%g", c.GravityAcceleration),
		fmt.Sprintf("%g", c.FirstSpaceVelocity),
		fmt.Sprintf("%g", c.SecondSpaceVelocity),
		strings.Join(c.Photos, "\n"),
	))

	if err != nil {
		return err
	} else {
		defer res.Close()
		return nil
	}
}

func GetFromDatabase(name string) (c CatalogueObject, err error) {
	conn := client.Connect()
	defer conn.Close()

	res, err := conn.Query(fmt.Sprintf("SELECT * FROM `catalogue` WHERE Name = '%s'", name))

	if err != nil {
		return CatalogueObject{}, err
	}

	for res.Next() {
		err = res.Scan(&c.Name, &c.Description, &c.OpeningDateTime, &c.SidericConversionPeriod, &c.BodyOrbitalVelocity,
			&c.Inclination, &c.Satelites, &c.WhoseSatelite, &c.EquatorialRadius, &c.PolarRadius, &c.AverageRadius, &c.Square,
			&c.Volume, &c.Weight, &c.AverageDensity, &c.GravityAcceleration, &c.FirstSpaceVelocity, &c.SecondSpaceVelocity, &c.Photos)

		if err != nil {
			return CatalogueObject{}, err
		}
	}

	return
}

func (c *CatalogueObject) AddPhoto(link string) {
	c.Photos = append(c.Photos, link)
	c.Save()
}

func (c *CatalogueObject) DeletePhoto(link string) {
	for i, v := range c.Photos {
		if v == link {
			c.Photos = append(c.Photos[:i], c.Photos[i+1:]...)
			break
		}
	}
}
