package user

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"gohub/internal/domain"
	"gohub/internal/service"
	"net/http"
)

type Handler struct {
	compiledEmailRegexp    *regexp.Regexp
	compiledPasswordRegexp *regexp.Regexp
	userService            *service.UserService
}

func InitHandler(userService *service.UserService) *Handler {
	const (
		emailRegexpPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		passwordRegexpPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	return &Handler{
		compiledEmailRegexp:    regexp.MustCompile(emailRegexpPattern, regexp.None),
		compiledPasswordRegexp: regexp.MustCompile(passwordRegexpPattern, regexp.None),
		userService:            userService,
	}
}

func (uh *Handler) RegisterRoutes(server *gin.RouterGroup) {
	server.POST("/sign-up", uh.SignUp)
	server.POST("/sign-in", uh.SignIn)
	server.PATCH("/:id", uh.UpdateUserInfo)
	server.GET("/:id", uh.GetUserInfo)
}

func (uh *Handler) SignUp(ctx *gin.Context) {
	type Req struct {
		Email             string `json:"email"`
		Password          string `json:"password"`
		ConfirmedPassword string `json:"confirmedPassword"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	if req.Password != req.ConfirmedPassword {
		ctx.String(http.StatusOK, "passwords do not match")
		return
	}

	ok, err := uh.compiledEmailRegexp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "invalid email")
		return
	}
	ok, err = uh.compiledPasswordRegexp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "invalid password")
		return
	}

	err = uh.userService.SignUp(
		ctx.Request.Context(),
		domain.User{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}

	ctx.String(http.StatusOK, "ok")
	return
}

func (uh *Handler) SignIn(ctx *gin.Context) {}

func (uh *Handler) UpdateUserInfo(ctx *gin.Context) {}

func (uh *Handler) GetUserInfo(ctx *gin.Context) {}
