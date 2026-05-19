package auth

import userfeature "backend-travelku/internal/user"

// Repository adalah data access layer untuk fitur auth.
// Saat ini hanya thin wrapper yang delegate ke user repository.
type Repository interface {
	Create(user *userfeature.User) error
	GetByEmail(email string) (*userfeature.User, error)
	ExistsByEmail(email string) bool
	Update(user *userfeature.User) error
}

type repository struct {
	userRepo userfeature.Repository
}

func NewRepository(userRepo userfeature.Repository) Repository {
	return &repository{userRepo: userRepo}
}

func (r *repository) Create(user *userfeature.User) error {
	return r.userRepo.Create(user)
}

func (r *repository) GetByEmail(email string) (*userfeature.User, error) {
	return r.userRepo.GetByEmail(email)
}

func (r *repository) ExistsByEmail(email string) bool {
	return r.userRepo.ExistsByEmail(email)
}

func (r *repository) Update(user *userfeature.User) error {
	return r.userRepo.Update(user)
}
