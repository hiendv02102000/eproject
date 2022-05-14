package mutation

import (
	"backend-food/graph/input"
	"backend-food/graph/output"
	"backend-food/internal/pkg/domain/domain_model/dto"
	"backend-food/internal/pkg/domain/domain_model/entity"
	"backend-food/internal/pkg/domain/service"
	"backend-food/pkg/share/utils"
	"errors"

	"github.com/graphql-go/graphql"
)

func CreateUserMutation(containerRepo map[string]interface{}) *graphql.Field {
	return &graphql.Field{
		Type:        output.CreateUserOutput(),
		Description: "User Register",

		Args: graphql.FieldConfigArgument{
			"user": &graphql.ArgumentConfig{
				Type: input.UserInput(),
			},
		},
		Resolve: func(p graphql.ResolveParams) (result interface{}, err error) {
			req := p.Args["user"].(map[string]interface{})
			createUserReq := dto.CreateUserRequest{
				Username:  req["username"].(string),
				Password:  req["password"].(string),
				FirstName: req["first_name"].(string),
				LastName:  req["last_name"].(string),
				Email:     req["email"].(string),
			}
			err = utils.CheckValidate(createUserReq)
			if err != nil {
				return
			}

			userRepo := containerRepo["user_repository"].(service.UserRepositoryInterface)
			user, err := userRepo.FirstUserListWithAnyCondition("SELECT * FROM users where username = ? OR email = ? ", createUserReq.Username, createUserReq.Email)

			if err != nil {
				return
			}

			if user.ID != 0 {
				if user.Username == createUserReq.Username {
					err = errors.New("username already exists")
					return
				}
				err = errors.New("email already exists")
				return
			}
			createUserReq.Password = utils.EncryptPassword(createUserReq.Password)
			user, err = userRepo.CreateUser(entity.Users{
				FirstName:    createUserReq.FirstName,
				LastName:     createUserReq.LastName,
				Username:     createUserReq.Username,
				Email:        createUserReq.Email,
				Password:     createUserReq.Password,
				Role:         entity.ClientRole,
				ActiveStatus: entity.ACTIVE_STATUS,
			})

			if err != nil {
				return
			}
			createUserRes := dto.CreateUserResponse{
				Username:  user.Username,
				FirstName: user.FirstName,
				LastName:  user.LastName,
				Email:     user.Email,
				Role:      string(user.Role),
				CreatedAt: utils.FormatTime(user.CreatedAt),
			}
			result = createUserRes
			return
		},
	}
}
