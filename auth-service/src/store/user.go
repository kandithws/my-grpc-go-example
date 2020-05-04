package store

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/kandithws/sharespace-api/auth-service/src/common/db"
	"github.com/kandithws/sharespace-api/auth-service/src/model"
	"github.com/lib/pq"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore() *UserStore {
	client := db.Client()
	// Migrate Here
	client.AutoMigrate(&model.User{})
	return &UserStore{db: client}
}

type UniqueConstriantError struct {
	Message string `json:"error_message"`
}

func (e *UniqueConstriantError) Error() string {
	return fmt.Sprintf(e.Message)
}

func makeUniqueConstriantError(err error) error {
	e, ok := err.(*pq.Error)
	if ok && e.Code.Name() == "unique_violation" {
		// handle error
		return &UniqueConstriantError{Message: e.Detail}
	}
	return nil

}

func (s *UserStore) CreateUser(m *model.User) error {

	err := s.db.Create(m).Error

	if err != nil {
		// Check for specific error types
		e := makeUniqueConstriantError(err)
		if e != nil {
			return e
		}
	}

	return err
}

func (s *UserStore) FindUserById(id uint) (*model.User, error) {
	var m model.User
	if err := s.db.First(&m, id).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (s *UserStore) FindUserBy(q *model.User) (*model.User, error) {
	var m model.User
	if err := s.db.Where(q).First(&m).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (s *UserStore) Update(m *model.User) error {
	err := s.db.Model(m).Updates(m).Error // update"s" non-blank
	if err != nil {
		// Check for specific error types
		e := makeUniqueConstriantError(err)
		if e != nil {
			return e
		}
	}
	return err
}

func (s *UserStore) Delete(m *model.User) error {
	return s.db.Delete(m).Error
}
