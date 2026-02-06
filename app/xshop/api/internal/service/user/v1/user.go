package v1

import (
	"context"
	"fmt"
	"mxshop/app/pkg/code"
	"mxshop/pkg/errors"
	"mxshop/pkg/log"
	"mxshop/pkg/storage"
	"time"

	"mxshop/app/mxshop/api/internal/data"
	"mxshop/app/pkg/options"
	"mxshop/gmicro/server/restserver/middlewares"

	"github.com/dgrijalva/jwt-go"
)

type UserDTO struct {
	data.User

	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

type UserSrv interface {
	MobileLogin(ctx context.Context, mobile, password string) (*UserDTO, error)
	Register(ctx context.Context, mobile, password, code string) (*UserDTO, error)
	Update(ctx context.Context, userDTO *UserDTO) error
	Get(ctx context.Context, userID uint64) (*UserDTO, error)
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
	CheckPassWord(ctx context.Context, password, EncryptedPassword string) (bool, error)
}

type userService struct {
	//ud data.UserData
	data data.DataFactory

	jwtOpts *options.JwtOptions
}

func NewUserService(data data.DataFactory, jwtOpts *options.JwtOptions) UserSrv {
	return &userService{data: data, jwtOpts: jwtOpts}
}

func (us *userService) MobileLogin(ctx context.Context, mobile, password string) (*UserDTO, error) {
	user, err := us.data.Users().GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}

	//检查密码是否正确
	err = us.data.Users().CheckPassWord(ctx, password, user.PassWord)
	if err != nil {
		return nil, err
	}

	//生成token
	j := middlewares.NewJWT(us.jwtOpts.Key)
	claims := middlewares.CustomClaims{
		ID:          uint(user.ID),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                                   //签名的生效时间
			ExpiresAt: (time.Now().Local().Add(us.jwtOpts.Timeout)).Unix(), //30天过期
			Issuer:    us.jwtOpts.Realm,
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		User:      user,
		Token:     token,
		ExpiresAt: (time.Now().Local().Add(us.jwtOpts.Timeout)).Unix(),
	}, nil
}

func (us *userService) Register(ctx context.Context, mobile, password, codes string) (*UserDTO, error) {
	rstore := storage.RedisCluster{}

	value, err := rstore.GetKey(ctx, fmt.Sprintf("%s_%d", mobile, 1))
	if err != nil {
		return nil, errors.WithCode(code.ErrCodeNotExist, "验证码不存在")
	}

	if value != codes {
		return nil, errors.WithCode(code.ErrCodeInCorrect, "验证码错误")
	}

	var user = &data.User{
		Mobile:   mobile,
		PassWord: password,
	}
	err = us.data.Users().Create(ctx, user)
	if err != nil {
		log.Errorf("user register failed: %v", err)
		return nil, err
	}

	//生成token
	j := middlewares.NewJWT(us.jwtOpts.Key)
	claims := middlewares.CustomClaims{
		ID:          uint(user.ID),
		NickName:    user.NickName,
		AuthorityId: uint(user.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                                   //签名的生效时间
			ExpiresAt: (time.Now().Local().Add(us.jwtOpts.Timeout)).Unix(), //30天过期
			Issuer:    us.jwtOpts.Realm,
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		return nil, err
	}

	return &UserDTO{
		User:      *user,
		Token:     token,
		ExpiresAt: (time.Now().Local().Add(us.jwtOpts.Timeout)).Unix(),
	}, nil
}

func (u *userService) Update(ctx context.Context, userDTO *UserDTO) error {
	//TODO implement me
	panic("implement me")
}

func (us *userService) Get(ctx context.Context, userID uint64) (*UserDTO, error) {
	userDO, err := us.data.Users().Get(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &UserDTO{User: userDO}, nil
}

func (u *userService) GetByMobile(ctx context.Context, mobile string) (*UserDTO, error) {
	//TODO implement me
	panic("implement me")
}

func (u *userService) CheckPassWord(ctx context.Context, password, EncryptedPassword string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

var _ UserSrv = &userService{}
