package user

import (
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/app/pkg/common"
	"Advanced_Shop/gnova/server/rpcserver"
	"Advanced_Shop/gnova/server/rpcserver/clientinterceptors"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
	"time"

	upbv1 "Advanced_Shop/api/user/v1"
	"Advanced_Shop/app/xshop/api/internal/data"
	"Advanced_Shop/gnova/registry"
	itime "Advanced_Shop/pkg/common/time"
)

const serviceName = "discovery:///xshop-user-srv"

type users struct {
	uc upbv1.UserClient
}

func NewUsers(uc upbv1.UserClient) *users {
	return &users{uc}
}

func NewUserServiceClient(r registry.Discovery) upbv1.UserClient {
	conn, err := rpcserver.DialInsecure(
		context.Background(),
		rpcserver.WithEndpoint(serviceName),
		rpcserver.WithDiscovery(r),
		rpcserver.WithClientUnaryInterceptor(clientinterceptors.UnaryTracingInterceptor),
	)
	if err != nil {
		panic(err)
	}
	c := upbv1.NewUserClient(conn)
	return c
}

func (u *users) CheckPassWord(ctx context.Context, password, encryptedPwd string) error {
	cres, err := u.uc.CheckPassWord(ctx, &upbv1.PasswordCheckInfo{
		Password:          password,
		EncryptedPassword: encryptedPwd,
	})
	if err != nil {
		log.Errorf("CheckPassWord err:%v", err)
		return err
	}
	if cres.Success {
		return nil
	}
	return errors.WithCode(code.ErrUserPasswordIncorrect, "密码错误")
}

func (u *users) Create(ctx context.Context, user *data.User) error {
	protoUser := &upbv1.CreateUserInfo{
		Mobile:   user.Mobile,
		NickName: user.NickName,
		PassWord: user.PassWord,
	}
	userRsp, err := u.uc.CreateUser(ctx, protoUser)
	if err != nil {
		log.Errorf("CreateUser err:%v", err)
		return err
	}
	user.ID = uint64(userRsp.Id)
	return err
}

func (u *users) Update(ctx context.Context, user *data.User) error {
	protoUser := &upbv1.UpdateUserInfo{
		Id:       int32(user.ID),
		NickName: user.NickName,
		Gender:   user.Gender,
		BirthDay: uint64(user.Birthday.Unix()),
	}
	_, err := u.uc.UpdateUser(ctx, protoUser)
	if err != nil {
		log.Errorf("UpdateUser err:%v", err)
		return err
	}
	return nil
}

func (u *users) Get(ctx context.Context, userID uint64) (data.User, error) {
	user, err := u.uc.GetUserById(ctx, &upbv1.IdRequest{
		Id: int32(userID),
	})
	if err != nil {
		log.Errorf("GetUser err:%v", err)
		return data.User{}, err
	}

	return data.User{
		ID:       uint64(user.Id),
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Birthday: itime.Time{time.Unix(int64(user.BirthDay), 0)},
		Gender:   user.Gender,
		Role:     user.Role,
		PassWord: user.PassWord,
	}, nil
}

func (u *users) List(ctx context.Context, pageInfo common.PageInfo) (data.UserList, error) {
	var response data.UserList

	list, err := u.uc.GetUserList(ctx, &upbv1.PageInfo{
		Pn:    uint32(pageInfo.Page),
		PSize: uint32(pageInfo.Limit),
	})
	if err != nil {
		log.Errorf("List err:%v", err)
		return response, err
	}
	response.TotalCount = int64(list.Total)

	var resp []*data.User
	for _, user := range list.Data {
		resp = append(resp, &data.User{
			ID:       uint64(user.Id),
			Mobile:   user.Mobile,
			NickName: user.NickName,
			Birthday: itime.Time{time.Unix(int64(user.BirthDay), 0)},
			Gender:   user.Gender,
			Role:     user.Role,
			PassWord: user.PassWord,
		})
	}
	response.Items = resp
	return response, nil
}

func (u *users) GetByMobile(ctx context.Context, mobile string) (data.User, error) {
	user, err := u.uc.GetUserByMobile(ctx, &upbv1.MobileRequest{
		Mobile: mobile,
	})
	if err != nil {
		log.Errorf("get user by mobile %s error: %v", mobile, err)
		return data.User{}, err
	}

	return data.User{
		ID:       uint64(user.Id),
		Mobile:   user.Mobile,
		NickName: user.NickName,
		Birthday: itime.Time{time.Unix(int64(user.BirthDay), 0)},
		Gender:   user.Gender,
		Role:     user.Role,
		PassWord: user.PassWord,
	}, nil
}

var _ data.UserData = &users{}
