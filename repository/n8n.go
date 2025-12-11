package repository

import (
	"main/models"
)

// CreateN8NWorkflow creates a new workflow record in the database
func (r *Repository) CreateN8NWorkflow(workflow *models.N8NWorkflow) error {
	return r.db.Create(workflow).Error
}

// GetN8NWorkflowByID retrieves a specific workflow by ID
func (r *Repository) GetN8NWorkflowByID(workflowID uint) (*models.N8NWorkflow, error) {
	var workflow models.N8NWorkflow
	err := r.db.Where("id = ?", workflowID).First(&workflow).Error
	if err != nil {
		return nil, err
	}
	return &workflow, nil
}

// GetN8NWorkflowsByUserID retrieves all workflows for a specific user
func (r *Repository) GetN8NWorkflowsByUserID(userID uint) ([]models.N8NWorkflow, error) {
	var workflows []models.N8NWorkflow
	err := r.db.Where("user_id = ?", userID).Order("updated_at DESC").Find(&workflows).Error
	if err != nil {
		return nil, err
	}
	return workflows, nil
}

// UpdateN8NWorkflow updates an existing workflow record
func (r *Repository) UpdateN8NWorkflow(workflow *models.N8NWorkflow) error {
	return r.db.Save(workflow).Error
}

// DeleteN8NWorkflow deletes a specific workflow by ID
func (r *Repository) DeleteN8NWorkflow(workflowID uint) error {
	return r.db.Delete(&models.N8NWorkflow{}, workflowID).Error
}