package service

import (
	"errors"
	"time"

	"github.com/miaomiaopu/ipmp-server/internal/model"
	"github.com/miaomiaopu/ipmp-server/internal/repository"
	"gorm.io/gorm"
)

var ErrCustomerNotFound = errors.New("customer not found")

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) List(page, pageSize int, keyword, status string) ([]model.Customer, int64, error) {
	return s.repo.List(page, pageSize, keyword, status)
}

func (s *CustomerService) GetByID(id string) (*model.Customer, int64, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, 0, ErrCustomerNotFound
		}
		return nil, 0, err
	}
	count, err := s.repo.CountProjects(id)
	if err != nil {
		return nil, 0, err
	}
	return customer, count, nil
}

func (s *CustomerService) Create(req *model.Customer) error {
	existing, err := s.repo.FindByCode(req.CustomerCode)
	if err == nil && existing != nil {
		return errors.New("customer code already exists")
	}
	return s.repo.Create(req)
}

func (s *CustomerService) Update(id string, updates map[string]interface{}) error {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCustomerNotFound
		}
		return err
	}
	// 应用更新
	for key, value := range updates {
		switch key {
		case "name":
			customer.Name = value.(string)
		case "customer_code":
			customer.CustomerCode = value.(string)
		case "contact_person":
			customer.ContactPerson = value.(string)
		case "contact_phone":
			customer.ContactPhone = model.EncryptedField(value.(string))
		case "contact_email":
			customer.ContactEmail = model.EncryptedField(value.(string))
		case "address":
			customer.Address = model.EncryptedField(value.(string))
		case "notes":
			customer.Notes = value.(string)
		case "status":
			customer.Status = value.(string)
		}
	}
	customer.UpdatedAt = time.Now()
	return s.repo.Update(customer)
}

func (s *CustomerService) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrCustomerNotFound
		}
		return err
	}
	return s.repo.SoftDelete(id)
}
