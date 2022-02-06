package models

import (
	"errors"
	"time"

	"github.com/sampalm/projectapi/database"
	"gorm.io/gorm"
)

type Project struct {
	ID          uint      `json:"-" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"size:255, unique, not null" validate:"min=3,max=40,regexp=^[a-zA-Z-0-9]*$"`
	DisplayName string    `json:"display_name" validate:"min=3,max=40"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"create_time"`
	UpdatedAt   time.Time `json:"update_time"`
}

type ProjectModel struct{}

func (p ProjectModel) Find(projectname string) (project Project, err error) {

	db := database.GetDatabase()
	err = db.First(&project, "name = ?", projectname).Error

	return project, err
}

func (p ProjectModel) All() (projects []Project, err error) {

	db := database.GetDatabase()
	err = db.Find(&projects).Error

	return projects, err

}

func (p ProjectModel) Save(project Project) error {

	db := database.GetDatabase()

	result := db.Find(&project, "name = ?", project.Name)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		return errors.New("project already exists")
	}

	return db.Create(&project).Error
}

func (p ProjectModel) Update(projectName string, project Project) error {

	db := database.GetDatabase()

	result := db.Model(&Project{}).Where("name = ?", projectName).Updates(project)

	if result.RowsAffected == 0 {
		return errors.New("could not update project")
	}

	return result.Error
}

func (p ProjectModel) Destroy(projectName string) error {

	db := database.GetDatabase()

	project, err := p.Find(projectName)
	if err != nil {
		return err
	}

	return db.Delete(&project).Error
}

func (ap *Project) BeforeUpdate(tx *gorm.DB) (err error) {
	if tx.Statement.Changed("Name") {
		return errors.New("project name not allowed to change")
	}
	if tx.Statement.Changed("CreatedAt") {
		return errors.New("created time not allowed to change")
	}
	return nil
}
