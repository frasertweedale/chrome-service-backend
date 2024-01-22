package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"gorm.io/datatypes"
)

type AvailableTemplates string

const (
	LandingPage AvailableTemplates = "landingPage"
)

func (at *AvailableTemplates) Scan(value interface{}) error {
	*at = AvailableTemplates(value.(string))
	return nil
}

func (at AvailableTemplates) Value() (driver.Value, error) {
	return string(at), nil
}

func (at AvailableTemplates) String() string {
	return string(at)
}

func (at AvailableTemplates) IsValid() error {
	switch at {
	case LandingPage:
		return nil
	}

	return fmt.Errorf("invalid dashboard template. Expected one of %s, got %s", LandingPage, at)
}

type TemplateConfig struct {
	Sx datatypes.JSON `gorm:"not null;default null" json:"sx"`
	Md datatypes.JSON `gorm:"not null;default null" json:"md"`
	Lg datatypes.JSON `gorm:"not null;default null" json:"lg"`
	Xl datatypes.JSON `gorm:"not null;default null" json:"xl"`
}

type GridItem struct {
	Title     string `json:"title"`
	ID        string `json:"i"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	Width     int    `json:"w"`
	Height    int    `json:"h"`
	MaxHeight int    `json:"maxH"`
	MinHeight int    `json:"minH"`
}

type GridSizes string

const (
	Sx GridSizes = "sx"
	Md GridSizes = "md"
	Lg GridSizes = "lg"
	Xl GridSizes = "xl"
)

func (gs GridSizes) IsValid() error {
	switch gs {
	case Sx, Md, Lg, Xl:
		return nil
	default:
		return errors.New(fmt.Errorf("invalid grid size, expected one of %s, %s, %s, %s", Sx, Md, Lg, Xl).Error())
	}
}

func (gs GridSizes) GetMaxWidth() (int, error) {
	err := gs.IsValid()
	if err != nil {
		return 0, err
	}
	switch gs {
	case Sx:
		return 1, nil
	case Md:
		return 2, nil
	case Lg:
		return 3, nil
	case Xl:
		return 4, nil
	default:
		return 0, errors.New("invalid grid size")
	}
}

func (gi GridItem) IsValid(variant GridSizes) error {
	if err := variant.IsValid(); err != nil {
		return err
	}

	if gi.ID == "" {
		return errors.New(`invalid grid item, field id "i" is required`)
	}

	if gi.Width < 1 || gi.Height < 1 || gi.MaxHeight < 1 || gi.MinHeight < 1 {
		return errors.New(`invalid grid item, height "h", width "w", maxHeight "maxH", mixHeight "minH" must be greater than 0`)
	}

	if gi.Height > gi.MaxHeight {
		return errors.New(fmt.Errorf(`invalid grid item, height "h" %d must be less than or equal to max height "maxH" %d`, gi.Height, gi.MaxHeight).Error())
	}

	if gi.Height < gi.MinHeight {
		return errors.New(fmt.Errorf(`invalid grid item, height "h" %d must be greater than or equal to min height "minH" %d`, gi.Height, gi.MinHeight).Error())
	}

	maxGridSize, err := variant.GetMaxWidth()
	if err != nil {
		return err
	}

	if gi.Width > maxGridSize {
		return errors.New(fmt.Errorf("invalid grid item, layout variant %s, width must be less than or equal to %d", variant, maxGridSize).Error())
	}

	if gi.X > maxGridSize {
		return errors.New(fmt.Errorf("invalid grid item, layout variant %s, coordinate X must be less than or equal to %d", variant, maxGridSize).Error())
	}

	return nil
}

func (tc *TemplateConfig) SetLayoutSizeItems(layoutSize string, items []GridItem) *TemplateConfig {
	bytes, err := json.Marshal(items)
	if err != nil {
		panic(err)
	}
	reflect.ValueOf(tc).Elem().FieldByName(layoutSize).Set(reflect.ValueOf(bytes))
	return tc
}

type DashboardTemplateBase struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
}

type DashboardTemplate struct {
	BaseModel
	UserIdentityID uint                  `json:"userIdentityID"`
	Default        bool                  `gorm:"not null;default:false" json:"default"`
	TemplateBase   DashboardTemplateBase `gorm:"not null;default null; embedded" 'json:"templateBase"`
	TemplateConfig TemplateConfig        `gorm:"not null;default null; embedded" json:"templateConfig"`
}

type BaseDashboardTemplate struct {
	Name           string         `json:"name"`
	DisplayName    string         `json:"displayName"`
	TemplateConfig TemplateConfig `json:"templateConfig"`
}

type BaseTemplates map[AvailableTemplates]BaseDashboardTemplate
