package user

import (
	v1 "Advanced_Shop/api/user/v1"
	"Advanced_Shop/app/pkg/code"
	"Advanced_Shop/app/pkg/gorm"
	DOv1 "Advanced_Shop/app/user/srv/data/v1"
	DTOv1 "Advanced_Shop/app/user/srv/service/v1"
	srv1 "Advanced_Shop/app/user/srv/service/v1"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/errors"
	"Advanced_Shop/pkg/log"
	"context"
	"github.com/google/wire"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"
)

var ProviderSet = wire.NewSet(NewUserServer)

func DTOToResponse(userDTO srv1.UserDTO) *v1.UserInfoResponse {
	//在grpc的message中字段有默认值，你不能随便赋值nil进去，容易出错
	//这里要搞清， 哪些字段是有默认值
	userInfoRsp := v1.UserInfoResponse{
		Id:       userDTO.ID,
		PassWord: userDTO.Password,
		NickName: userDTO.NickName,
		Gender:   userDTO.Gender,
		Role:     int32(userDTO.Role),
		Mobile:   userDTO.Mobile,
	}
	if userDTO.Birthday != nil {
		userInfoRsp.BirthDay = uint64(userDTO.Birthday.Unix())
	}
	return &userInfoRsp
}

type userServer struct {
	v1.UnimplementedUserServer
	srv srv1.UserSrv
}

// NewUserServer java中的ioc，控制翻转 ioc = injection of control
// 代码分层，第三方服务， rpc， redis， 等等， 带来一定的复杂度
func NewUserServer(srv srv1.UserSrv) v1.UserServer {
	return &userServer{srv: srv}
}

var _ v1.UserServer = &userServer{}

func (u *userServer) GetUserList(ctx context.Context, info *v1.PageInfo) (*v1.UserListResponse, error) {
	log.Info("GetUserList is called")
	srvOpts := metav1.ListMeta{
		Page:     int(info.Pn),
		PageSize: int(info.PSize),
	}
	dtoList, err := u.srv.List(ctx, []string{}, srvOpts)
	if err != nil {
		return nil, err
	}

	var rsp v1.UserListResponse
	for _, value := range dtoList.Items {
		userRsp := DTOToResponse(*value)
		rsp.Data = append(rsp.Data, userRsp)
	}
	return &rsp, nil
}

func (u *userServer) GetUserByMobile(ctx context.Context, request *v1.MobileRequest) (*v1.UserInfoResponse, error) {
	log.Infof("get user by mobile function called.")
	user, err := u.srv.GetByMobile(ctx, request.Mobile)
	if err != nil {
		log.Errorf("get user by mobile: %s, error: %v", request.Mobile, err)
		return nil, err
	}

	userInfoRsp := DTOToResponse(*user)
	return userInfoRsp, nil
}

func (u *userServer) GetUserById(ctx context.Context, request *v1.IdRequest) (*v1.UserInfoResponse, error) {
	log.Infof("get user by id function called.")
	user, err := u.srv.GetByID(ctx, uint64(request.Id))
	if err != nil {
		log.Errorf("get user by id: %s, error: %v", request.Id, err)
		return nil, err
	}

	userInfoRsp := DTOToResponse(*user)
	return userInfoRsp, nil
}

func (u *userServer) CreateUser(ctx context.Context, info *v1.CreateUserInfo) (*v1.UserInfoResponse, error) {
	log.Infof("create user function called.")

	userDO := DOv1.UserDO{
		Mobile:   info.Mobile,
		NickName: info.NickName,
		Password: info.PassWord,
	}
	userDTO := DTOv1.UserDTO{userDO}

	err := u.srv.Create(ctx, &userDTO)
	if err != nil {
		log.Errorf("create user: %v, error: %v", userDTO, err)
		return nil, err
	}

	userInfoRsp := DTOToResponse(userDTO)
	return userInfoRsp, nil
}

func (u *userServer) UpdateUser(ctx context.Context, info *v1.UpdateUserInfo) (*emptypb.Empty, error) {
	log.Infof("update user function called.")

	birthDay := time.Unix(int64(info.BirthDay), 0)
	userDO := DOv1.UserDO{
		Model: gorm.Model{
			ID: info.Id,
		},
		NickName: info.NickName,
		Gender:   info.Gender,
		Birthday: &birthDay,
	}
	userDTO := DTOv1.UserDTO{UserDO: userDO}

	err := u.srv.Update(ctx, &userDTO)
	if err != nil {
		log.Errorf("update user: %v, error: %v", userDTO, err)
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (u *userServer) CheckPassWord(ctx context.Context, info *v1.PasswordCheckInfo) (*v1.CheckResponse, error) {
	if info.EncryptedPassword == info.Password {
		return &v1.CheckResponse{Success: true}, nil
	}
	return &v1.CheckResponse{Success: false}, errors.WithCode(code.ErrUserPasswordIncorrect, "password err", nil)
}

// TODO
//func (u *userServer) CheckPassWord(ctx context.Context, info *v1.PasswordCheckInfo) (*v1.CheckResponse, error) {
//	//校验密码
//	options := &password.Options{16, 100, 32, sha512.New}
//	passwordInfo := strings.Split(info.EncryptedPassword, "$")
//	check := password.Verify(info.Password, passwordInfo[2], passwordInfo[3], options)
//	return &v1.CheckResponse{Success: check}, nil
//}

func (u *userServer) mustEmbedUnimplementedUserServer() {
	//TODO implement me
	panic("implement me")
}
