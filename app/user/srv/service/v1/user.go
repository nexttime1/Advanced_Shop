package v1

import (
	"Advanced_Shop/app/pkg/code"
	dv1 "Advanced_Shop/app/user/srv/data/v1"
	metav1 "Advanced_Shop/pkg/common/meta/v1"
	"Advanced_Shop/pkg/errors"
	"context"
)

type UserDTO struct {
	dv1.UserDO
}

type UserSrv interface {
	List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDTOList, error)
	Create(ctx context.Context, user *UserDTO) error
	Update(ctx context.Context, user *UserDTO) error
	GetByID(ctx context.Context, ID uint64) (*UserDTO, error)
	GetByMobile(ctx context.Context, mobile string) (*UserDTO, error)
}

type userService struct {
	userStore dv1.UserStore
}

func (u *userService) Create(ctx context.Context, user *UserDTO) error {
	//先判断用户是否存在
	_, err := u.userStore.GetByMobile(ctx, user.Mobile)
	if err != nil && errors.IsCode(err, code.ErrUserNotFound) {
		return u.userStore.Create(ctx, &user.UserDO)
	}

	//这里应该区别到底是什么错误，用户已经存在？ 数据访问错误？
	return errors.WithCode(code.ErrUserAlreadyExists, "用户已经存在")
}

func (u *userService) Update(ctx context.Context, user *UserDTO) error {
	//先查询用户是否存在
	_, err := u.userStore.GetByID(ctx, uint64(user.ID))
	if err != nil {
		return err
	}

	return u.userStore.Update(ctx, &user.UserDO)
}

func (u *userService) GetByID(ctx context.Context, ID uint64) (*UserDTO, error) {
	userDO, err := u.userStore.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return &UserDTO{*userDO}, nil
}

func (u *userService) GetByMobile(ctx context.Context, mobile string) (*UserDTO, error) {
	userDO, err := u.userStore.GetByMobile(ctx, mobile)
	if err != nil {
		return nil, err
	}

	return &UserDTO{*userDO}, nil
}

func NewUserService(us dv1.UserStore) UserSrv {
	return &userService{
		userStore: us,
	}
}

var _ UserSrv = &userService{}

type UserDTOList struct {
	TotalCount int64      `json:"totalCount,omitempty"` //总数
	Items      []*UserDTO `json:"data"`                 //数据
}

func (u *userService) List(ctx context.Context, orderby []string, opts metav1.ListMeta) (*UserDTOList, error) {

	doList, err := u.userStore.List(ctx, orderby, opts)
	if err != nil {
		return nil, err
	}

	var userDTOList UserDTOList
	for _, value := range doList.Items {
		projectDTO := UserDTO{*value}
		userDTOList.Items = append(userDTOList.Items, &projectDTO)
	}

	return &userDTOList, nil
}
