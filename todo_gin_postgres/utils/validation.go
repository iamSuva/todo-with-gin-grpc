package utils

import (
	"fmt"
	"regexp"
	"todowithgin/models"

	"github.com/go-playground/validator/v10"
)

type ErrResponse struct {
	Message string `json:"message"`
}

func GetMessage(tag string) string {
	switch tag {
	case "required":
		return "required"
	case "min":
		return "too short atleast 5 character needed"
	case "password":
		return "must include at least 8 characters, one uppercase, one lowercase, one number, and one special character"
		
	default:
		return "invalid input"
	}
}

func ValidateTaskHandler(task models.Task) ([]ErrResponse, error) {
	myValidator := validator.New()

	err := myValidator.Struct(task)
	if err != nil {
		fmt.Println(err)
		var errResponses []ErrResponse
		for _, verr := range err.(validator.ValidationErrors) {
			errResponses = append(errResponses, ErrResponse{
				Message: fmt.Sprintf("%s is %s", verr.Field(), GetMessage(verr.Tag())),
			})
		}
		return errResponses, err
	}
	return nil, nil

}
func ValidateUserHandler(user models.User) ([]ErrResponse, error) {
	userValidator := validator.New()
	userValidator.RegisterValidation("password", validatePassword)
	err := userValidator.Struct(user)
	if err != nil {
		fmt.Println("user validation: ", err)
		var errResponse []ErrResponse
		for _, verr := range err.(validator.ValidationErrors) {
			errResponse = append(errResponse, ErrResponse{
				Message: fmt.Sprintf("%s is %s", verr.Field(), GetMessage(verr.Tag())),
			})
		}
		return errResponse, err
	}
	return nil, nil
}

func validatePassword(f1 validator.FieldLevel) bool {
	password := f1.Field().String()
	fmt.Println(password)

	isvalid := len(password) >= 8 &&
		regexp.MustCompile(`[A-Z]`).MatchString(password) &&
		regexp.MustCompile(`[a-z]`).MatchString(password) &&
		regexp.MustCompile(`[0-9]`).MatchString(password) &&
		regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString(password)

	return isvalid

}
