// services/user_services.go

package services

import (
    "errors"
    "auth-user-api/models"
    "auth-user-api/repository"
    "auth-user-api/utils"

    "golang.org/x/crypto/bcrypt"
)

type UserService interface {
    Register(username, email, password1, password2 string, role int) error
    Update(id, username, email, password1, password2 string) error
    Delete(id string) error
    Authenticate(username, password string) (*models.User, error)
    GetAllUsers() ([]*models.User, error)
    GetUserByID(id string) (*models.User, error)
    GetUserByUsername(username string) (*models.User, error)
}

type userService struct {
    repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
    return &userService{repo}
}

// Register - Untuk mendaftarkan user baru
func (s *userService) Register(username, email, password1, password2 string, role int) error {
    if password1 != password2 {
        return errors.New("password didn't match")
    }

    if err := utils.ValidatePassword(password1); err != nil {
        return err
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := &models.User{
        Username: username,
        Email:    email,
        Password: string(hashedPassword),
        Role:     role, // Set the role here
    }

    return s.repo.CreateUser(user)
}

// GetAllUsers - Mendapatkan semua user
func (s *userService) GetAllUsers() ([]*models.User, error) {
    users, err := s.repo.GetAllUsers()
    if err != nil {
        return nil, err
    }
    return users, nil
}

// Update - Mengupdate data user
func (s *userService) Update(id, username, email, password1, password2 string) error {
    user, err := s.repo.GetUserByID(id)
    if err != nil {
        return err
    }

    // Update username jika diberikan
    if username != "" {
        user.Username = username
    }

    // Update email jika diberikan
    if email != "" {
        user.Email = email
    }

    // Update password jika diberikan dan valid
    if password1 != "" || password2 != "" {
        if password1 != password2 {
            return errors.New("password didn't match")
        }

        if err := utils.ValidatePassword(password1); err != nil {
            return err
        }

        hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
        if err != nil {
            return err
        }

        user.Password = string(hashedPassword) // Simpan password yang sudah di-hash
    }

    // Update user di database
    return s.repo.UpdateUser(user)
}

// Delete - Menghapus user
func (s *userService) Delete(id string) error {
    return s.repo.DeleteUser(id)
}

// Authenticate - Autentikasi user berdasarkan username dan password
func (s *userService) Authenticate(username, password string) (*models.User, error) {
    user, err := s.repo.GetUserByUsername(username)
    if err != nil {
        if err.Error() == "record not found" {
            return nil, errors.New("user not found")
        }
        return nil, err
    }

    if user.DeletedAt.Valid {
        return nil, errors.New("user not found")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
        return nil, errors.New("invalid username or password")
    }

    return user, nil
}

// GetUserByID - Mengambil user berdasarkan ID
func (s *userService) GetUserByID(id string) (*models.User, error) {
    return s.repo.GetUserByID(id)
}

// GetUserByUsername - Mengambil user berdasarkan username
func (s *userService) GetUserByUsername(username string) (*models.User, error) {
    return s.repo.GetUserByUsername(username)
}
