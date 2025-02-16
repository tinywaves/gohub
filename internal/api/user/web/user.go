package web

import (
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"gohub/internal"
	"gohub/internal/api/user/domain"
	"gohub/internal/api/user/service"
	"net/http"
)

type UserHandler struct {
	compiledEmailRegexp    *regexp.Regexp
	compiledPasswordRegexp *regexp.Regexp
	userService            *service.UserService
}

func InitUserHandler(userService *service.UserService) *UserHandler {
	const (
		emailRegexpPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		passwordRegexpPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	)
	return &UserHandler{
		compiledEmailRegexp:    regexp.MustCompile(emailRegexpPattern, regexp.None),
		compiledPasswordRegexp: regexp.MustCompile(passwordRegexpPattern, regexp.None),
		userService:            userService,
	}
}

func (uh *UserHandler) RegisterRoutes(server *gin.RouterGroup) {
	server.POST("/sign-up", uh.SignUp)
	server.POST("/sign-in", uh.SignIn)
	server.PATCH("", uh.UpdateUserInfo)
	server.GET("", uh.GetUserInfo)
	server.POST("/sign-out", uh.SignOut)
}

func (uh *UserHandler) SignUp(ctx *gin.Context) {
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
		if errors.Is(err, service.ErrUserEmailDuplicated) {
			ctx.String(http.StatusOK, err.Error())
			return
		}
		ctx.String(http.StatusOK, "system error")
		return
	}

	ctx.String(http.StatusOK, "ok")
	return
}

func (uh *UserHandler) SignIn(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := uh.userService.SignIn(
		ctx.Request.Context(),
		domain.User{
			Email:    req.Email,
			Password: req.Password,
		},
	)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ctx.String(http.StatusOK, "your email is not registered")
			return
		}
		if errors.Is(err, service.ErrUserEmailPasswordNotMatch) {
			ctx.String(http.StatusOK, err.Error())
			return
		}
		ctx.String(http.StatusOK, "system error")
		return
	}

	session := sessions.Default(ctx)
	session.Set(internal.SessionDataKey, user.Id)
	err = session.Save()
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}

	ctx.String(http.StatusOK, "ok")
	return
}

func (uh *UserHandler) UpdateUserInfo(ctx *gin.Context) {
	type Req struct {
		Nickname string `json:"nickname"`
		Bio      string `json:"bio"`
		Birthday int64  `json:"birthday"`
		Gender   string `json:"gender"`
	}

	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}

	session := sessions.Default(ctx)
	userId := session.Get(internal.SessionDataKey).(string)

	err := uh.userService.UpdateUserInfo(
		ctx.Request.Context(),
		domain.User{
			Id:       userId,
			Nickname: req.Nickname,
			Bio:      req.Bio,
			Birthday: req.Birthday,
			Gender:   req.Gender,
		},
	)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ctx.String(http.StatusOK, "your account is not registered")
			return
		}
		ctx.String(http.StatusOK, "system error")
		return
	}

	ctx.String(http.StatusOK, "ok")
	return
}

func (uh *UserHandler) GetUserInfo(ctx *gin.Context) {
	session := sessions.Default(ctx)
	userId := session.Get(internal.SessionDataKey).(string)

	user, err := uh.userService.GetUserInfo(ctx.Request.Context(), userId)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			ctx.String(http.StatusOK, "your account is not registered")
			return
		}
		ctx.String(http.StatusOK, "system error")
		return
	}

	ctx.JSON(http.StatusOK, user)
	return
}

func (uh *UserHandler) SignOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	err := session.Save()
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}

	ctx.String(http.StatusOK, "ok")
	return
}
