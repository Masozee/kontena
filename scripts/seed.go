package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Masozee/kontena/api/database"
	"github.com/Masozee/kontena/api/models"
)

func main() {
	// Initialize database
	database.InitDB()

	// Create tenants
	tenants := createTenants()

	// For each tenant, create test data
	for _, tenant := range tenants {
		createPeople(tenant)
		createProjects(tenant)
	}

	log.Println("Seed data created successfully!")
}

func createTenants() []models.Tenant {
	tenants := []models.Tenant{
		{
			Name:        "Acme Corporation",
			Description: "A global conglomerate with diverse business interests",
			Plan:        "enterprise",
			Status:      "active",
			Domain:      "acme.example.com",
			LogoURL:     "https://example.com/logos/acme.png",
			Settings:    `{"theme":"dark","language":"en","notifications":true}`,
		},
		{
			Name:        "Stark Industries",
			Description: "A technology company specializing in advanced weaponry and clean energy",
			Plan:        "enterprise",
			Status:      "active",
			Domain:      "stark.example.com",
			LogoURL:     "https://example.com/logos/stark.png",
			Settings:    `{"theme":"light","language":"en","notifications":true}`,
		},
		{
			Name:        "Wayne Enterprises",
			Description: "A multinational conglomerate with interests in various sectors",
			Plan:        "enterprise",
			Status:      "active",
			Domain:      "wayne.example.com",
			LogoURL:     "https://example.com/logos/wayne.png",
			Settings:    `{"theme":"dark","language":"en","notifications":false}`,
		},
	}

	for i := range tenants {
		result := database.DB.Create(&tenants[i])
		if result.Error != nil {
			log.Fatalf("Failed to create tenant: %v", result.Error)
		}
		fmt.Printf("Created tenant: %s (ID: %d)\n", tenants[i].Name, tenants[i].ID)
	}

	return tenants
}

func createPeople(tenant models.Tenant) []models.Person {
	roles := []string{"Project Manager", "Developer", "Designer", "QA Engineer", "Product Owner"}
	people := []models.Person{
		{
			TenantID: tenant.ID,
			Name:     "John Doe",
			Email:    fmt.Sprintf("john.doe@%s", tenant.Domain),
			Role:     roles[0],
			Position: "Senior Project Manager",
			Phone:    "+1-555-123-4567",
			Avatar:   "https://example.com/avatars/john.png",
		},
		{
			TenantID: tenant.ID,
			Name:     "Jane Smith",
			Email:    fmt.Sprintf("jane.smith@%s", tenant.Domain),
			Role:     roles[1],
			Position: "Senior Developer",
			Phone:    "+1-555-234-5678",
			Avatar:   "https://example.com/avatars/jane.png",
		},
		{
			TenantID: tenant.ID,
			Name:     "Bob Johnson",
			Email:    fmt.Sprintf("bob.johnson@%s", tenant.Domain),
			Role:     roles[2],
			Position: "UI/UX Designer",
			Phone:    "+1-555-345-6789",
			Avatar:   "https://example.com/avatars/bob.png",
		},
		{
			TenantID: tenant.ID,
			Name:     "Alice Williams",
			Email:    fmt.Sprintf("alice.williams@%s", tenant.Domain),
			Role:     roles[3],
			Position: "QA Lead",
			Phone:    "+1-555-456-7890",
			Avatar:   "https://example.com/avatars/alice.png",
		},
		{
			TenantID: tenant.ID,
			Name:     "Charlie Brown",
			Email:    fmt.Sprintf("charlie.brown@%s", tenant.Domain),
			Role:     roles[4],
			Position: "Product Manager",
			Phone:    "+1-555-567-8901",
			Avatar:   "https://example.com/avatars/charlie.png",
		},
	}

	for i := range people {
		result := database.DB.Create(&people[i])
		if result.Error != nil {
			log.Fatalf("Failed to create person: %v", result.Error)
		}
		fmt.Printf("Created person: %s (ID: %d) for tenant %s\n", people[i].Name, people[i].ID, tenant.Name)
	}

	return people
}

func createProjects(tenant models.Tenant) []models.Project {
	// Get people for this tenant
	var people []models.Person
	database.DB.Where("tenant_id = ?", tenant.ID).Find(&people)

	// Create projects
	projects := []models.Project{
		{
			TenantID:    tenant.ID,
			Name:        "Website Redesign",
			Description: "Redesign the company website with a modern look and feel",
			Budget:      50000.00,
			StartDate:   time.Now(),
			EndDate:     func() *time.Time { t := time.Now().AddDate(0, 3, 0); return &t }(),
			Status:      "in_progress",
		},
		{
			TenantID:    tenant.ID,
			Name:        "Mobile App Development",
			Description: "Develop a mobile app for iOS and Android",
			Budget:      75000.00,
			StartDate:   time.Now(),
			EndDate:     func() *time.Time { t := time.Now().AddDate(0, 6, 0); return &t }(),
			Status:      "planning",
		},
		{
			TenantID:    tenant.ID,
			Name:        "CRM Integration",
			Description: "Integrate our systems with a new CRM solution",
			Budget:      30000.00,
			StartDate:   time.Now().AddDate(0, -1, 0),
			EndDate:     func() *time.Time { t := time.Now().AddDate(0, 2, 0); return &t }(),
			Status:      "in_progress",
		},
	}

	for i := range projects {
		// Create the project
		result := database.DB.Create(&projects[i])
		if result.Error != nil {
			log.Fatalf("Failed to create project: %v", result.Error)
		}
		fmt.Printf("Created project: %s (ID: %d) for tenant %s\n", projects[i].Name, projects[i].ID, tenant.Name)

		// Associate people with the project
		for j, person := range people {
			if j < 3 { // Assign first 3 people to each project
				database.DB.Model(&projects[i]).Association("People").Append(&person)
			}
		}

		// Create KPIs for the project
		createKPIs(projects[i])

		// Create tasks for the project
		createTasks(projects[i], people)

		// Create milestones for the project
		createMilestones(projects[i])

		// Create risks for the project
		createRisks(projects[i])
	}

	return projects
}

func createKPIs(project models.Project) []models.KPI {
	kpis := []models.KPI{
		{
			ProjectID:    project.ID,
			Description:  "Complete project on time",
			TargetValue:  100,
			CurrentValue: 25,
			Unit:         "%",
		},
		{
			ProjectID:    project.ID,
			Description:  "Stay within budget",
			TargetValue:  project.Budget,
			CurrentValue: project.Budget * 0.3,
			Unit:         "$",
		},
		{
			ProjectID:    project.ID,
			Description:  "Customer satisfaction",
			TargetValue:  5,
			CurrentValue: 4.2,
			Unit:         "rating",
		},
	}

	for i := range kpis {
		// Update achievement status
		kpis[i].UpdateAchievement()

		// Save to database
		result := database.DB.Create(&kpis[i])
		if result.Error != nil {
			log.Fatalf("Failed to create KPI: %v", result.Error)
		}
		fmt.Printf("Created KPI: %s (ID: %d) for project %s\n", kpis[i].Description, kpis[i].ID, project.Name)
	}

	return kpis
}

func createTasks(project models.Project, people []models.Person) []models.Task {
	tasks := []models.Task{
		{
			ProjectID:   project.ID,
			Title:       "Requirements gathering",
			Description: "Gather and document project requirements",
			Status:      models.TaskStatusCompleted,
			DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, -7); return &t }(),
		},
		{
			ProjectID:   project.ID,
			Title:       "Design phase",
			Description: "Create wireframes and design mockups",
			Status:      models.TaskStatusInProgress,
			DueDate:     func() *time.Time { t := time.Now().AddDate(0, 0, 7); return &t }(),
		},
		{
			ProjectID:   project.ID,
			Title:       "Development",
			Description: "Implement the design and functionality",
			Status:      models.TaskStatusTodo,
			DueDate:     func() *time.Time { t := time.Now().AddDate(0, 1, 0); return &t }(),
		},
		{
			ProjectID:   project.ID,
			Title:       "Testing",
			Description: "Perform QA testing and fix bugs",
			Status:      models.TaskStatusTodo,
			DueDate:     func() *time.Time { t := time.Now().AddDate(0, 2, 0); return &t }(),
		},
		{
			ProjectID:   project.ID,
			Title:       "Deployment",
			Description: "Deploy the project to production",
			Status:      models.TaskStatusTodo,
			DueDate:     func() *time.Time { t := time.Now().AddDate(0, 2, 15); return &t }(),
		},
	}

	for i := range tasks {
		// Assign a person to the task
		personIndex := i % len(people)
		tasks[i].AssignedToID = &people[personIndex].ID

		// Save to database
		result := database.DB.Create(&tasks[i])
		if result.Error != nil {
			log.Fatalf("Failed to create task: %v", result.Error)
		}
		fmt.Printf("Created task: %s (ID: %d) for project %s\n", tasks[i].Title, tasks[i].ID, project.Name)

		// Create time entries for completed and in-progress tasks
		if tasks[i].Status == models.TaskStatusCompleted || tasks[i].Status == models.TaskStatusInProgress {
			createTimeEntries(tasks[i], people[personIndex])
		}
	}

	return tasks
}

func createTimeEntries(task models.Task, person models.Person) []models.TimeTracking {
	// Create 3 time entries for the task
	timeEntries := []models.TimeTracking{
		{
			TaskID:   task.ID,
			PersonID: person.ID,
			Hours:    2.5,
			Date:     time.Now().AddDate(0, 0, -5),
		},
		{
			TaskID:   task.ID,
			PersonID: person.ID,
			Hours:    3.0,
			Date:     time.Now().AddDate(0, 0, -3),
		},
		{
			TaskID:   task.ID,
			PersonID: person.ID,
			Hours:    4.0,
			Date:     time.Now().AddDate(0, 0, -1),
		},
	}

	for i := range timeEntries {
		result := database.DB.Create(&timeEntries[i])
		if result.Error != nil {
			log.Fatalf("Failed to create time entry: %v", result.Error)
		}
		fmt.Printf("Created time entry: %.1f hours (ID: %d) for task %s\n", timeEntries[i].Hours, timeEntries[i].ID, task.Title)
	}

	return timeEntries
}

func createMilestones(project models.Project) []models.Milestone {
	milestones := []models.Milestone{
		{
			ProjectID:   project.ID,
			Title:       "Project Kickoff",
			Description: "Initial project kickoff meeting and requirements gathering",
			DueDate:     time.Now().AddDate(0, 0, -14),
			Status:      models.MilestoneStatusCompleted,
		},
		{
			ProjectID:   project.ID,
			Title:       "Design Approval",
			Description: "Get client approval on design mockups",
			DueDate:     time.Now().AddDate(0, 0, 7),
			Status:      models.MilestoneStatusInProgress,
		},
		{
			ProjectID:   project.ID,
			Title:       "Alpha Release",
			Description: "Release alpha version for internal testing",
			DueDate:     time.Now().AddDate(0, 1, 0),
			Status:      models.MilestoneStatusPlanned,
		},
		{
			ProjectID:   project.ID,
			Title:       "Beta Release",
			Description: "Release beta version for client testing",
			DueDate:     time.Now().AddDate(0, 2, 0),
			Status:      models.MilestoneStatusPlanned,
		},
		{
			ProjectID:   project.ID,
			Title:       "Final Release",
			Description: "Release final version to production",
			DueDate:     time.Now().AddDate(0, 3, 0),
			Status:      models.MilestoneStatusPlanned,
		},
	}

	for i := range milestones {
		result := database.DB.Create(&milestones[i])
		if result.Error != nil {
			log.Fatalf("Failed to create milestone: %v", result.Error)
		}
		fmt.Printf("Created milestone: %s (ID: %d) for project %s\n", milestones[i].Title, milestones[i].ID, project.Name)
	}

	return milestones
}

func createRisks(project models.Project) []models.Risk {
	risks := []models.Risk{
		{
			ProjectID:   project.ID,
			Description: "Client may change requirements",
			Impact:      "High",
			Probability: "Medium",
			Mitigation:  "Regular client meetings and sign-off on requirements",
			Status:      models.RiskStatusMonitoring,
		},
		{
			ProjectID:   project.ID,
			Description: "Team members may be unavailable",
			Impact:      "Medium",
			Probability: "Low",
			Mitigation:  "Cross-train team members and have backup resources",
			Status:      models.RiskStatusIdentified,
		},
		{
			ProjectID:   project.ID,
			Description: "Technology limitations",
			Impact:      "High",
			Probability: "Low",
			Mitigation:  "Conduct technical feasibility study early in the project",
			Status:      models.RiskStatusMitigated,
		},
	}

	for i := range risks {
		result := database.DB.Create(&risks[i])
		if result.Error != nil {
			log.Fatalf("Failed to create risk: %v", result.Error)
		}
		fmt.Printf("Created risk: %s (ID: %d) for project %s\n", risks[i].Description, risks[i].ID, project.Name)
	}

	return risks
}
