package handlers

import (
	"strconv"

	"github.com/Masozee/kontena/api/database"
	"github.com/Masozee/kontena/api/models"
	"github.com/gofiber/fiber/v2"
)

// GetTasks retrieves all tasks for a specific project
// @Summary Get all tasks for a project
// @Description Get all tasks for a specific project
// @Tags tasks
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param project_id path int true "Project ID"
// @Success 200 {array} models.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{project_id}/tasks [get]
func GetTasks(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id")
	if tenantIDStr == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID is required",
		})
	}

	tenantID, err := strconv.Atoi(tenantIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID format",
		})
	}

	projectID, err := c.ParamsInt("project_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID format",
		})
	}

	// Verify project belongs to tenant
	var project models.Project
	result := database.DB.Where("id = ? AND tenant_id = ?", projectID, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	var tasks []models.Task
	result = database.DB.Where("project_id = ?", projectID).Preload("AssignedTo").Find(&tasks)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve tasks: " + result.Error.Error(),
		})
	}

	return c.JSON(tasks)
}

// GetTask retrieves a specific task by ID
// @Summary Get a task by ID
// @Description Get a specific task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [get]
func GetTask(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id")
	if tenantIDStr == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID is required",
		})
	}

	tenantID, err := strconv.Atoi(tenantIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID format",
		})
	}

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID format",
		})
	}

	var task models.Task
	result := database.DB.Preload("AssignedTo").First(&task, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	// Verify task belongs to a project owned by the tenant
	var project models.Project
	result = database.DB.Where("id = ? AND tenant_id = ?", task.ProjectID, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied to this task",
		})
	}

	return c.JSON(task)
}

// CreateTask creates a new task for a project
// @Summary Create a task
// @Description Create a new task for a project
// @Tags tasks
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param project_id path int true "Project ID"
// @Param task body models.Task true "Task object"
// @Success 201 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /projects/{project_id}/tasks [post]
func CreateTask(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id")
	if tenantIDStr == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID is required",
		})
	}

	tenantID, err := strconv.Atoi(tenantIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID format",
		})
	}

	projectID, err := c.ParamsInt("project_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid project ID format",
		})
	}

	// Verify project belongs to tenant
	var project models.Project
	result := database.DB.Where("id = ? AND tenant_id = ?", projectID, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Project not found",
		})
	}

	task := new(models.Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Validate required fields
	if task.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task title is required",
		})
	}

	task.ProjectID = uint(projectID)

	// Verify assigned person belongs to the same tenant if provided
	if task.AssignedToID != nil && *task.AssignedToID > 0 {
		var person models.Person
		result = database.DB.Where("id = ? AND tenant_id = ?", *task.AssignedToID, tenantID).First(&person)
		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Assigned person not found or not in the same tenant",
			})
		}
	}

	result = database.DB.Create(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create task: " + result.Error.Error(),
		})
	}

	// Load the assigned person if exists
	if task.AssignedToID != nil {
		database.DB.Preload("AssignedTo").First(&task, task.ID)
	}

	return c.Status(fiber.StatusCreated).JSON(task)
}

// UpdateTask updates a task by ID
// @Summary Update a task
// @Description Update a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param id path int true "Task ID"
// @Param task body models.Task true "Task object"
// @Success 200 {object} models.Task
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [put]
func UpdateTask(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id")
	if tenantIDStr == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID is required",
		})
	}

	tenantID, err := strconv.Atoi(tenantIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID format",
		})
	}

	taskID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID format",
		})
	}

	// Find the task
	var task models.Task
	result := database.DB.Preload("Project").First(&task, taskID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	// Verify the task's project belongs to the tenant
	var project models.Project
	result = database.DB.Where("id = ? AND tenant_id = ?", task.ProjectID, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied: task does not belong to this tenant",
		})
	}

	// Parse the updated task data
	updatedTask := new(models.Task)
	if err := c.BodyParser(updatedTask); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body: " + err.Error(),
		})
	}

	// Validate required fields
	if updatedTask.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Task title is required",
		})
	}

	// Verify assigned person belongs to the same tenant if provided
	if updatedTask.AssignedToID != nil && *updatedTask.AssignedToID > 0 {
		var person models.Person
		result = database.DB.Where("id = ? AND tenant_id = ?", *updatedTask.AssignedToID, tenantID).First(&person)
		if result.Error != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Assigned person not found or not in the same tenant",
			})
		}
	}

	// Update the task
	task.Title = updatedTask.Title
	task.Description = updatedTask.Description
	task.Status = updatedTask.Status
	task.DueDate = updatedTask.DueDate
	task.AssignedToID = updatedTask.AssignedToID

	result = database.DB.Save(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update task: " + result.Error.Error(),
		})
	}

	// Load the assigned person if exists
	if task.AssignedToID != nil {
		database.DB.Preload("AssignedTo").First(&task, task.ID)
	}

	return c.Status(fiber.StatusOK).JSON(task)
}

// DeleteTask deletes a task by ID
// @Summary Delete a task
// @Description Delete a task by ID
// @Tags tasks
// @Accept json
// @Produce json
// @Param tenant_id header string true "Tenant ID"
// @Param id path int true "Task ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
func DeleteTask(c *fiber.Ctx) error {
	tenantIDStr := c.Locals("tenant_id")
	if tenantIDStr == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Tenant ID is required",
		})
	}

	tenantID, err := strconv.Atoi(tenantIDStr.(string))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid tenant ID format",
		})
	}

	taskID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid task ID format",
		})
	}

	// Find the task
	var task models.Task
	result := database.DB.Preload("Project").First(&task, taskID)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Task not found",
		})
	}

	// Verify the task's project belongs to the tenant
	var project models.Project
	result = database.DB.Where("id = ? AND tenant_id = ?", task.ProjectID, tenantID).First(&project)
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied: task does not belong to this tenant",
		})
	}

	// Delete the task
	result = database.DB.Delete(&task)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete task: " + result.Error.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Task deleted successfully",
	})
}
