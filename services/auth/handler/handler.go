package handler

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"gorm.io/gorm"

	log "hotwave/logger"
	"hotwave/services/auth/cache"
	"hotwave/services/auth/config"
	"hotwave/services/auth/database"
	"hotwave/services/auth/database/models"
	"hotwave/services/auth/proto"
)

func New(c *config.Config) (*Handler, error) {
	udb, err := database.CreateMysqlClient(c.UserDBDSN)
	if err != nil {
		return nil, err
	}

	return &Handler{
		userCache: cache.NewMemory(),
		userDB:    udb,
	}, nil
}

type Handler struct {
	userDB    *gorm.DB
	userCache cache.AuthCache
	conf      *config.Config
}

func generateToken(u *models.Users) *jwt.Token {
	// Id: strconv.FormatInt(u.Id, 10)

	//TODO:
	claims := &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
		ID:        uuid.NewString(),
		Audience:  []string{u.Uname},
		NotBefore: jwt.NewNumericDate(time.Now()),
		Issuer:    "hotwave-auth",
		Subject:   "assess",
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

// func isLoginMethod(method string) bool {
// 	methodList := []string{
// 		"hotwave.Login",
// 	}
// 	method = strings.ToLower(method)
// 	for _, v := range methodList {
// 		if method == strings.ToLower(v) {
// 			return true
// 		}
// 	}
// 	return false
// }

// func LogWrapper(fn server.HandlerFunc) server.HandlerFunc {
// 	return func(ctx context.Context, req server.Request, resp interface{}) error {
// 		//b := time.Now()
// 		//log.Debugf("[Log Wrapper] Before serving request method: %v", req.Method())
// 		err := fn(ctx, req, resp)
// 		//log.Debugf("[Log Wrapper] After serving request method: %v,cost:%d", req.Method(), time.Since(b).Milliseconds())
// 		return err
// 	}
// }

// func LoginWrapper(fn server.HandlerFunc) server.HandlerFunc {
// 	return func(ctx context.Context, req server.Request, resp interface{}) error {
// 		tokenRaw, present := metadata.Get(ctx, "Authorization")
// 		tokenRaw = strings.TrimPrefix(tokenRaw, "Bearer ")
// 		if !present || len(tokenRaw) == 0 {
// 			if isLoginMethod(req.Method()) {
// 				return fn(ctx, req, resp)
// 			}
// 			return fmt.Errorf("token is invalid")
// 			//TODO:
// 			//return errors.Unauthorized()
// 		}

// 		claims := &jwt.StandardClaims{}
// 		token, err := jwt.ParseWithClaims(tokenRaw, claims, func(t *jwt.Token) (interface{}, error) {
// 			return "", nil
// 		})
// 		if err != nil || !token.Valid {
// 			return fmt.Errorf("StatusUnauthorized")
// 		}
// 		// user, _ := cache.FetchUser(claims.Id)
// 		// if user == nil {
// 		// 	user, err = TokenLogin(token)
// 		// 	if err != nil || user == nil {
// 		// 		return fmt.Errorf("can'nt found user info with token")
// 		// 	}
// 		// }
// 		// if user.Token != tokenRaw {
// 		// 	return fmt.Errorf("token is expired, please login again")
// 		// }
// 		// newCtx := context.WithValue(ctx, cache.UserCTXKey, user)
// 		// return fn(newCtx, req, resp)
// 		return nil
// 	}
// }

func TokenLogin(t *jwt.Token) (*cache.Memory, error) {
	if !t.Valid {
		return nil, fmt.Errorf("token is not valid")
	}

	claim := t.Claims.(*jwt.RegisteredClaims)
	if claim == nil {
		return nil, fmt.Errorf("token is not valid")
	}
	uid, err := strconv.ParseInt(claim.ID, 10, 64)
	if err != nil {
		return nil, err
	}

	fmt.Println(uid)

	// dbuser := gDBConn.GetUser(claim.Id)
	// if dbuser == nil {
	// 	return nil, fmt.Errorf("token is not valid")
	// }
	// rawToken, _ := t.SignedString(JwtKey)
	// u := &cache.UserCache{
	// 	Dbuser: dbuser,
	// 	Token:  rawToken,
	// }
	// err := cache.StoreUser(u)
	// if err != nil {
	// 	log.Error(err)
	// }
	return nil, nil
}

func (*Handler) Captcha(ctx context.Context, in *proto.CaptchaRequest) (*proto.CaptchaResponse, error) {
	//TODO:
	return &proto.CaptchaResponse{}, nil
}

func (h *Handler) LoginWithPasswd(ctx context.Context, in *proto.LoginWithPasswdRequest) (*proto.LoginWithPasswdResponse, error) {
	if len(in.Account) < 4 {
		return nil, fmt.Errorf("please input right account")
	}
	if len(in.Password) < 6 {
		return nil, fmt.Errorf("passwd is required")
	}

	// cacheUser, _ := cache.FetchUser(in.Username)
	// if cacheUser != nil {
	// out.AssessToken = cacheUser.Token
	// return nil
	// }

	user := &models.Users{
		Uname: in.Account,
	}

	res := h.userDB.Debug().Limit(1).Find(user, user)

	if err := res.Error; err != nil {
		log.Error(err)
		return nil, fmt.Errorf("server internal error")
	}

	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("account is not exist")
	}

	if user.Passwd != in.Password {
		return nil, fmt.Errorf("passwd is not right")
	}

	if user.Stat != 0 {
		return nil, fmt.Errorf("stat error %d", user.Stat)
	}

	rawToekn := generateToken(user)

	assess, err := rawToekn.SignedString(h.conf.JwtKey)
	if err != nil {
		return nil, err
	}

	out := &proto.LoginWithPasswdResponse{}
	out.Token = &proto.Token{
		AssessToken: assess,
	}

	cacheInfo := &cache.AuthCacheInfo{
		User:        user,
		AssessToken: assess,
	}

	if err = h.userCache.StoreUser(ctx, cacheInfo, time.Hour); err != nil {
		log.Error(err)
	}

	out.UserInfo = &proto.UserInfo{
		Userid:   user.UID,
		Username: user.Uname,
		Stat:     int32(user.Stat),
		Created:  user.CreateAt.Unix(),
	}

	//TODO: event and log
	// loginEvent := &proto.UserLoginEvent{}
	// log.Info(loginEvent)
	return out, nil
}

func (*Handler) Logout(ctx context.Context, in *proto.LogoutRequest) (*proto.LogoutResponse, error) {
	return nil, nil
}

func (*Handler) RefreshToken(ctx context.Context, in *proto.RefreshTokenRequest) (*proto.RefreshTokenResponse, error) {
	//TODO
	return nil, nil
}

func (h *Handler) UserInfo(ctx context.Context, in *proto.UserInfoRequest) (*proto.UserInfoResponse, error) {
	claims, err := h.VerifyAdminTokenWithClaims(in.AccessToken)
	if err != nil {
		return nil, err
	}

	user := h.userCache.FetchUserByName(ctx, claims.Audience[0])

	if user == nil {
		return nil, fmt.Errorf("user no found")
	}

	out := &proto.UserInfoResponse{}

	out.Info = &proto.UserInfo{
		Userid:   user.User.UID,
		Username: user.User.Uname,
		Stat:     int32(user.User.Stat),
		Created:  user.User.CreateAt.Unix(),
	}
	return out, nil
}

func (h *Handler) VerifyAdminTokenWithClaims(tokenRaw string) (*jwt.RegisteredClaims, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenRaw, claims, func(t *jwt.Token) (interface{}, error) {
		return h.conf.JwtKey, nil
	})
	return claims, err
}
