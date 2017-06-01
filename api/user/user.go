package user

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"io/ioutil"

	"github.com/gernest/cereal/internal"
	"github.com/gernest/cereal/messages"
	"github.com/gernest/cereal/models"
	"github.com/ngorm/ngorm"
)

// CreateUserRequest is the JSON object expected when creating anew user.
type CreateUserRequest struct {
	Name     string `json:"username"`
	Password string `json:"password"`
}

// Valid performs field validation
func (c *CreateUserRequest) Valid() *messages.Message {
	var m *messages.Message
	if c.Name == "" {
		m = &messages.Message{
			Message: messages.FailedValidation,
			Errors: []messages.Error{
				{
					Field:    "username",
					Resource: "user",
					Code:     messages.CodeMissing,
				},
			},
		}
	} else if len(c.Name) < 5 {
		if m == nil {
			m = &messages.Message{
				Message: messages.FailedValidation,
				Errors: []messages.Error{
					{
						Field:    "username",
						Resource: "user",
						Code:     messages.CodeAlreadyExists,
					},
				},
			}
		} else {
			m.Errors = append(m.Errors, messages.Error{
				Field:    "username",
				Resource: "user",
				Code:     messages.CodeAlreadyExists,
			})
		}
	}
	if c.Password == "" {
		if m == nil {
			m = &messages.Message{
				Message: messages.FailedValidation,
				Errors: []messages.Error{
					{
						Field:    "password",
						Resource: "user",
						Code:     messages.CodeMissing,
					},
				},
			}
		} else {
			m.Errors = append(m.Errors, messages.Error{
				Field:    "password",
				Resource: "user",
				Code:     messages.CodeMissing,
			})
		}
	}
	return m
}

// Create creates a new user based on username and password. This accepts json
// request and returns only StatusCreated 201 when ucess and appropriate error
// when something is not right
func Create(db *ngorm.DB, w http.ResponseWriter, r *http.Request) {
	req := &CreateUserRequest{}
	m := &messages.Message{}
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		m.Message = messages.InvalidJSON
		internal.WriteJSON(w, m, http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(b, req); err != nil {
		m.Message = messages.BadJSON
		internal.WriteJSON(w, m, http.StatusBadRequest)
		return
	}
	if v := req.Valid(); v != nil {
		internal.WriteJSON(w, v, http.StatusUnprocessableEntity)
		return
	}
	if userExists(db, req.Name) {
		m.Message = messages.FailedValidation
		m.Errors = append(m.Errors, messages.Error{
			Field:    "username",
			Resource: "user",
			Code:     messages.CodeAlreadyExists,
		})
		internal.WriteJSON(w, m, http.StatusBadRequest)
		return
	}
	h, err := HashString(req.Password)
	if err != nil {
		m.Message = messages.FailedValidation
		m.Errors = append(m.Errors, messages.Error{
			Field:    "password",
			Resource: "user",
			Code:     messages.CodeInvalid,
		})
		internal.WriteJSON(w, m, http.StatusBadRequest)
		return
	}
	u := models.User{
		Name:     req.Name,
		Password: h,
		Profie: models.Profile{
			Name: req.Name,
		},
	}
	if err := db.Create(&u); err != nil {
		//TODO:(gernest) Properly hadnle this error?
		//
		//It is considered bad practice to send this error back to the client,
		//we need first to identify what information is relevant to clients.
		//
		// The status code will be 500
	}
	internal.WriteJSON(w, messages.OK(), http.StatusCreated)
}

func userExists(db *ngorm.DB, username string) bool {
	if err := db.First(&models.User{}, &models.User{Name: username}); err != nil {
		return false
	}
	return true
}

// HashString uses bcrypt to encrypt the string
func HashString(secret string) (string, error) {
	s, err := bcrypt.GenerateFromPassword([]byte(secret), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(s), nil
}

// CompareHashedString compares the two passwords, returns true if they match
func CompareHashedString(hashed, str string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(str))
}
