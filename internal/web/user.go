package web

import (
	"errors"
	"net/http"
	"time"

	"github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/tinywaves/gohub/internal"
	"github.com/tinywaves/gohub/internal/domain"
	"github.com/tinywaves/gohub/internal/service"
)

type userHandler struct {
	uSvc                          *service.UserService
	emailRegexpPatternCompiled    *regexp2.Regexp
	passwordRegexpPatternCompiled *regexp2.Regexp
}

func InitUserHandler(uSvc *service.UserService) *userHandler {
	const (
		emailRegexpPattern    = `^\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		passwordRegexpPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,72}$`
	)
	emailRegexpPatternCompiled := regexp2.MustCompile(emailRegexpPattern, regexp2.Compiled)
	passwordRegexpPatternCompiled := regexp2.MustCompile(passwordRegexpPattern, regexp2.Compiled)
	return &userHandler{
		uSvc:                          uSvc,
		emailRegexpPatternCompiled:    emailRegexpPatternCompiled,
		passwordRegexpPatternCompiled: passwordRegexpPatternCompiled,
	}
}

func (uh *userHandler) RegisterRoutes(userGroup *gin.RouterGroup) {
	userGroup.POST("/sign-up", uh.signUp)
	userGroup.POST("/sign-in", uh.signIn)
	userGroup.PATCH("", uh.edit)
	userGroup.GET("", uh.profile)
}

func (uh *userHandler) signUp(ctx *gin.Context) {
	type req struct {
		Email             string `json:"email"`
		Password          string `json:"password"`
		ConfirmedPassword string `json:"confirmedPassword"`
	}
	var reqBody req
	// Bind parses the request body into reqBody based on the Content-Type.
	// If parsing fails, return a 4xx error.
	if err := ctx.Bind(&reqBody); err != nil {
		return
	}

	if reqBody.Password != reqBody.ConfirmedPassword {
		ctx.String(http.StatusOK, "passwords do not match")
		return
	}

	ok, err := uh.emailRegexpPatternCompiled.MatchString(reqBody.Email)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "invalid email")
		return
	}
	ok, err = uh.passwordRegexpPatternCompiled.MatchString(reqBody.Password)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "invalid password")
		return
	}

	err = uh.uSvc.SignUp(ctx, domain.User{
		Email:    reqBody.Email,
		Password: reqBody.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrServiceEmailDuplicate) {
			ctx.String(http.StatusOK, err.Error())
			return
		}
		ctx.String(http.StatusOK, "system error")
		return
	}

	ctx.String(http.StatusOK, "signUp success")
}

func (uh *userHandler) signIn(ctx *gin.Context) {
	type req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var reqBody req
	if err := ctx.Bind(&reqBody); err != nil {
		return
	}

	if reqBody.Email == "" || reqBody.Password == "" {
		ctx.String(http.StatusOK, "invalid email or password")
		return
	}

	user, err := uh.uSvc.SignIn(ctx, reqBody.Email, reqBody.Password)
	if err != nil {
		if errors.Is(err, service.ErrServiceEmailPasswordNotMatch) {
			ctx.String(http.StatusOK, err.Error())
			return
		}
		ctx.String(http.StatusOK, "system error")
		return
	}

	jwtToken := jwt.NewWithClaims(
		jwt.SigningMethodEdDSA,
		internal.UserClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(internal.JwtTokenDuration)),
			UserId:    user.UserId,
			UserAgent: ctx.Request.UserAgent(),
		},
	)
	jwtTokenString, err := jwtToken.SignedString(internal.Ed25519PrivateKey)
	if err != nil {
		ctx.String(http.StatusOK, "system error")
		return
	}
	ctx.Header(internal.JwtTokenHeaderKey, jwtTokenString)
	ctx.String(http.StatusOK, "signIn success")
}

func (uh *userHandler) edit(ctx *gin.Context) {
	type req struct {
		Nickname string `json:"nickname"`
		Bio      string `json:"bio"`
		Gender   int    `json:"gender"`
		Birthday int    `json:"birthday"`
	}
	var reqBody req
	if err := ctx.Bind(&reqBody); err != nil {
		return
	}

	err := uh.uSvc.Edit(ctx, domain.User{
		UserId:   internal.HandleUserId(ctx),
		Nickname: reqBody.Nickname,
		Bio:      reqBody.Bio,
		Gender:   reqBody.Gender,
		Birthday: reqBody.Birthday,
	})
	if err != nil {
		if errors.Is(err, service.ErrServiceUserNotFound) {
			ctx.String(http.StatusOK, err.Error())
			return
		}
		ctx.String(http.StatusOK, "system error")
		return
	}
	ctx.String(http.StatusOK, "edit success")
}

func (uh *userHandler) profile(ctx *gin.Context) {
	user, err := uh.uSvc.Profile(ctx, internal.HandleUserId(ctx))
	if err != nil {
		if errors.Is(err, service.ErrServiceUserNotFound) {
			ctx.String(http.StatusOK, err.Error())
			return
		}
		ctx.String(http.StatusOK, "system error")
		return
	}
	ctx.JSON(http.StatusOK, user)
}
