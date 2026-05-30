package repository

import (
	"github.com/miaomiaopu/ipmp-server/internal/model"
	"gorm.io/gorm"
)

type CustomerRepository struct {
	db *gorm.DB
}

func NewCustomerRepository(db *gorm.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

func (r *CustomerRepository) FindByID(id string) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.Where("id = ?", id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) FindByCode(code string) (*model.Customer, error) {
	var customer model.Customer
	err := r.db.Where("customer_code = ?", code).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *CustomerRepository) List(page, pageSize int, keyword, status string) ([]model.Customer, int64, error) {
	var customers []model.Customer
	var total int64
	query := r.db.Model(&model.Customer{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR customer_code LIKE ? OR contact_person LIKE ?", like, like, like)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&customers).Error
	return customers, total, err
}

func (r *CustomerRepository) Create(customer *model.Customer) error {
	return r.db.Create(customer).Error
}

func (r *CustomerRepository) Update(customer *model.Customer) error {
	return r.db.Save(customer).Error
}

func (r *CustomerRepository) SoftDelete(id string) error {
	return r.db.Where("id = ?", id).Delete(&model.Customer{}).Error
}

// CountProjects 统计客户关联的项目数
func (r *CustomerRepository) CountProjects(customerID string) (int64, error) {
	var count int64
	err := r.db.Model(&model.Project{}).Where("customer_id = ?", customerID).Count(&count).Error
	return count, err
}

func (r *CustomerRepository) ForceDelete(id string) error { return ForceDelete(r.db, &model.Customer{}, id) }
func (r *CustomerRepository) Restore(id string) error { return Restore(r.db, &model.Customer{}, id) }
