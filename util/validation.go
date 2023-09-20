package util

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type (
	CreateAccountRequest struct {
		Owner    string `validate:"required,min=5,max=20"` // Required field, min 5 char long max 20
		Currency string `validate:"required,contains=USD"` // Required field, and client needs to implement our 'teener' tag format which we'll see later
	}

	ErrorResponse struct {
		Error       bool
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
)

// type User struct {
// 	FirstName      string     `json:"fname"`
// 	LastName       string     `json:"lname"`
// 	Age            uint8      `validate:"gte=0,lte=130"`
// 	Email          string     `json:"e-mail" validate:"required,email"`
// 	FavouriteColor string     `validate:"hexcolor|rgb|rgba"`
// 	Addresses      []*Address `validate:"required,dive,required"` // a person can have a home and cottage...
// 	Gender         Gender     `json:"gender" validate:"required,gender_custom_validation"`
// }

var validate = validator.New(validator.WithRequiredStructEnabled())

var MyValidator = &XValidator{
	validator: validate,
}

func (v *XValidator) Validate(data interface{}) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := validate.Struct(data)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			// In this case data object is actually holding the User struct
			var elem ErrorResponse

			elem.FailedField = err.Field() // Export struct field name
			elem.Tag = err.Tag()           // Export struct tag
			elem.Value = err.Value()       // Export field value
			elem.Error = true

			validationErrors = append(validationErrors, elem)
		}
	}

	return validationErrors
}

func CheckValidationErrors(req interface{}) *fiber.Error {

	if errs := MyValidator.Validate(req); len(errs) > 0 && errs[0].Error {
		errMsgs := make([]string, 0)

		for _, err := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf(
				"[%s]: '%v' | Needs to implement '%s'",
				err.FailedField,
				err.Value,
				err.Tag,
			))
		}

		return &fiber.Error{
			Code:    fiber.ErrBadRequest.Code,
			Message: strings.Join(errMsgs, " and "),
		}
	}
	return nil

}

// func main() {

// 	app := fiber.New(fiber.Config{
// 		// Global custom error handler
// 		ErrorHandler: func(c *fiber.Ctx, err error) error {
// 			return c.Status(fiber.StatusBadRequest).JSON(GlobalErrorHandlerResp{
// 				Success: false,
// 				Message: err.Error(),
// 			})
// 		},
// 	})

// 	app.Get("/", func(c *fiber.Ctx) error {
// 		user := &CreateAccountRequest{
// 			Owner:    c.Query("owner"),
// 			Currency: c.Query("currency"),
// 		}

// 		// Validation

// 		// Logic, validated with success
// 		return c.SendString("Hello, World!")
// 	})

// 	log.Fatal(app.Listen(":3000"))
// }

/**
OUTPUT

[1]
Request:

GET http://127.0.0.1:3000/

Response:

{"success":false,"message":"[Name]: '' | Needs to implement 'required' and [Age]: '0' | Needs to implement 'required'"}

[2]
Request:

GET http://127.0.0.1:3000/?name=efdal&age=9

Response:
{"success":false,"message":"[Age]: '9' | Needs to implement 'teener'"}

[3]
Request:

GET http://127.0.0.1:3000/?name=efdal&age=

Response:
{"success":false,"message":"[Age]: '0' | Needs to implement 'required'"}

[4]
Request:

GET http://127.0.0.1:3000/?name=efdal&age=18

Response:
Hello, World!

**/
