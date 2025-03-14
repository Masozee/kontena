# Project Management System

A comprehensive project management system built with Go and GORM.

## Models

### Project
- Core entity that contains all project information
- Has relationships with all other entities
- Tracks budget, timeline, and various project components

### Person
- Represents team members or stakeholders
- Can be assigned to projects and tasks
- Tracks roles and contact information

### KPI (Key Performance Indicator)
- Measures project performance
- Includes target and current values
- Automatically calculates progress and achievement status

### Task
- Represents work items within a project
- Can be assigned to team members
- Tracks status and due dates

### Report
- Documents project status and updates
- Linked to specific projects

### Milestone
- Represents significant project checkpoints
- Tracks status and due dates

### Risk
- Identifies potential project risks
- Includes impact, probability, and mitigation strategies

### Issue
- Tracks problems that arise during the project
- Can be assigned to team members for resolution

### Document
- Stores project-related files and attachments
- Tracks who uploaded each document

### TimeTracking
- Records time spent on tasks
- Links tasks to people who worked on them

## Database Schema

The system uses a relational database with the following schema:

```
Project
  ├── People (many-to-many)
  ├── KPIs (one-to-many)
  ├── Tasks (one-to-many)
  ├── Reports (one-to-many)
  ├── Milestones (one-to-many)
  ├── Risks (one-to-many)
  ├── Issues (one-to-many)
  └── Documents (one-to-many)

Task
  ├── AssignedTo (many-to-one with Person)
  └── TimeEntries (one-to-many with TimeTracking)

TimeTracking
  ├── Task (many-to-one)
  └── Person (many-to-one)
```

## Getting Started

1. Clone the repository
2. Set up your database connection in `.env`
3. Run the migrations to create the database schema
4. Start building your project management application

## Environment Variables

Create a `.env` file with the following variables:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=project_management
```

## Database Migration

The system uses GORM's AutoMigrate to create and update the database schema:

```go
db.AutoMigrate(
    &models.Project{},
    &models.Person{},
    &models.KPI{},
    &models.Task{},
    &models.Report{},
    &models.Milestone{},
    &models.Risk{},
    &models.Issue{},
    &models.Document{},
    &models.TimeTracking{},
)
```

## License

MIT 