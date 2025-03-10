package handler

import (
	"elearning_api/dto"
	"elearning_api/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler interface {
	GetUserProfile(c *gin.Context)
	Confirm(c *gin.Context)
}

type UserHandlerImplement struct {
	UserService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return &UserHandlerImplement{
		UserService: userService,
	}
}

// @Summary      Get User Profile
// @Description  Get User Profile
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dto.UserProfileResponse
// @Router       /user/profile [get]
func (uh UserHandlerImplement) GetUserProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := uh.UserService.GetUserProfile(userID)

	if err != nil {
		if apiErr, ok := err.(*dto.ApiError); ok {
			c.JSON(apiErr.Code, dto.UserProfileResponse{
				Data:  dto.UserBasicInfo{},
				Error: apiErr.Message,
				Meta:  dto.MetaResponse{},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.UserProfileResponse{
			Data:  dto.UserBasicInfo{},
			Error: "unexpected error",
			Meta:  dto.MetaResponse{},
		})
		return
	}
	c.JSON(http.StatusOK, dto.UserProfileResponse{
		Data:  user,
		Error: "",
		Meta:  dto.MetaResponse{},
	})
}

// @Summary Logout
// @Summary Confirm
// @Description Confirm
// @Tags User
// @Accept json
// @Produce json
// @Param token body dto.UserConfirmRequest true "User Confirm Request"
// @Success 200 {object} dto.UserConfirmResponse
// @Router /user/confirm [post]
func (uh UserHandlerImplement) Confirm(c *gin.Context) {
	request := dto.UserConfirmRequest{}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, dto.UserConfirmResponse{
			Error: "Error when bind data : " + err.Error(),
		})
		return
	}

	err := uh.UserService.Confirm(request)

	if err != nil {
		if apiErr, ok := err.(*dto.ApiError); ok {
			c.JSON(apiErr.Code, dto.UserConfirmResponse{
				Data:  dto.Message{},
				Error: apiErr.Message,
				Meta:  dto.MetaResponse{},
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.UserConfirmResponse{
			Data:  dto.Message{},
			Error: "unexpected error",
			Meta:  dto.MetaResponse{},
		})
		return
	}
	c.JSON(http.StatusOK, dto.UserConfirmResponse{
		Data: dto.Message{
			Message: "Account has been confirmed",
		},
		Error: "",
		Meta:  dto.MetaResponse{},
	})
}
