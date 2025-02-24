package taskService

type TaskService struct {
    repo TaskRepository
}

func NewService(repo TaskRepository) *TaskService {
    return &TaskService{repo: repo}
}

func (s *TaskService) CreateTask(task Tasks) (Tasks, error) {
    return s.repo.CreateTask(task)
}

func (s *TaskService) GetAllTasks() ([]Tasks, error) {
    return s.repo.GetAllTasks()
}

func (s *TaskService) UpdateTaskByID(id uint, task Tasks) (Tasks, error) {
    return s.repo.UpdateTaskByID(id, task)
}

func (s *TaskService) DeleteTaskByID(id uint) error {
    return s.repo.DeleteTaskByID(id)
}

func (s *TaskService) PatchTaskByID(id uint, updates map[string]interface{}) (Tasks, error) {
 return s.repo.PatchTaskByID(id, updates)
}