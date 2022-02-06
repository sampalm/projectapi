package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/sampalm/projectapi/database"
	"gorm.io/gorm"
)

type API struct {
	ID          uint      `json:"-" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null" validate:"min=3,max=40,regexp=^[a-zA-Z-0-9]*$"`
	Description string    `json:"description"`
	Version     string    `json:"version" gorm:"not null" validate:"nonzero"`
	OpenApiFile string    `json:"openapi_file" validate:"nonzero"`
	ProjectName string    `json:"project_name" gorm:"not null" validate:"min=3,max=40,regexp=^[a-zA-Z-0-9]*$"`
	CreatedAt   time.Time `json:"create_time"`
	UpdatedAt   time.Time `json:"update_time"`
}

type APIModel struct{}

func (a APIModel) Find(apiName, projectName string) (api API, err error) {
	db := database.GetDatabase()

	err = db.First(&api, "name = ? AND project_name = ?", apiName, projectName).Error
	return api, err
}

func (a APIModel) All() (apis []API, err error) {
	db := database.GetDatabase()

	err = db.Find(&apis).Error
	return apis, err
}

func (a APIModel) Save(projectName string, api API) error {
	db := database.GetDatabase()

	result, err := a.Find(api.Name, projectName)
	if err != nil && err.Error() != gorm.ErrRecordNotFound.Error() {
		return err
	}
	if result.ProjectName == projectName {
		return fmt.Errorf("project already has this API registered")
	}

	return db.Create(&api).Error
}

func (a APIModel) Update(apiName, projectName string, api API) error {
	db := database.GetDatabase()

	result := db.Model(&API{}).Where("name = ? AND project_name = ? ", apiName, projectName).Updates(api)

	if result.RowsAffected == 0 {
		return errors.New("invalid api or project")
	}

	return result.Error
}

func (a APIModel) Destroy(apiName, projectName string) error {
	db := database.GetDatabase()

	api, err := a.Find(apiName, projectName)
	if err != nil {
		return err
	}
	if api.ProjectName != projectName {
		return fmt.Errorf("project does not has this API registered")
	}

	return db.Delete(&api).Error
}

func (ap *API) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("ProjectName") {
		return errors.New("project not allowed to change")
	}
	if tx.Statement.Changed("Name") {
		return errors.New("api name not allowed to change")
	}
	if tx.Statement.Changed("CreatedAt") {
		return errors.New("created time not allowed to change")
	}
	return nil
}
