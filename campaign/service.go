package campaign

type Service interface {
	FindCampaigns(userId int) ([]Campaign, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// fetch campaign from database
// if client sent user_id,
// then fetch all campaigns that belong to that user_id
// if not, then fetch all campaigns from database
func (s *service) FindCampaigns(userId int) ([]Campaign, error) {
	if userId != 0 {
		campaigns, err := s.repository.FindByUserId(userId)
		if err != nil {
			return campaigns, err
		}
		return campaigns, nil
	}

	campaigns, err := s.repository.FindAll()
	if err != nil {
		return campaigns, err
	}
	return campaigns, nil
}
