package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"reflect"

	"github.com/kandithws/sharespace-api/auth-service/src/common/validator"
	"github.com/kandithws/sharespace-api/auth-service/src/genproto"
	"github.com/kandithws/sharespace-api/auth-service/src/model"
	"github.com/kandithws/sharespace-api/auth-service/src/store"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func NewAuthGrpcHandler(gserver *grpc.Server) {
	s := &handler{userStore: store.NewUserStore(), val: validator.NewValidator()}

	genproto.RegisterAuthServiceServer(gserver, s)
}

// AuthServiceHandler implements genproto.AuthServiceServer
type handler struct {
	userStore *store.UserStore
	val       validator.Validator
}

var (
	ErrRequestBindingError = errors.New("Fail to bind the request")
)

func bindJSON(in interface{}, out interface{}) error {
	if reflect.TypeOf(out).Kind() != reflect.Ptr {
		return ErrRequestBindingError
	}

	if reflect.ValueOf(out).IsNil() {
		return ErrRequestBindingError
	}

	byt, _ := json.Marshal(in)
	if err := json.Unmarshal(byt, out); err != nil {
		return err
	}

	return nil
}

func hashAndSalt(pwdStr string) string {

	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)

	pwd := []byte(pwdStr)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	}
	// GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func (h *handler) Register(c context.Context, req *genproto.RegisterRequest) (*genproto.RegisterResponse, error) {
	m := &model.User{}
	// Bind
	if err := bindJSON(req, m); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	m.Password = hashAndSalt(req.Password)

	if err := h.val.Validate(m); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := h.userStore.CreateUser(m); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &genproto.RegisterResponse{Message: "OK"}, nil
}

func comparePasswords(hashedPwd string, plainPwdString string) bool {
	// Since we'll be getting the hashed password from the DB it
	// will be a string so we'll need to convert it to a byte slice
	byteHash := []byte(hashedPwd)

	plainPwd := []byte(plainPwdString)

	err := bcrypt.CompareHashAndPassword(byteHash, plainPwd)
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}

func (h *handler) Login(c context.Context, req *genproto.LoginRequest) (*genproto.LoginResponse, error) {

	// Find by username
	user, err := h.userStore.FindUserBy(&model.User{Username: req.Username})

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "Username not found")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if !comparePasswords(user.Password, req.Password) {
		return nil, status.Error(codes.Unauthenticated, "Unauthorized")
	}

	// Generate JWT token at API Gateway

	return &genproto.LoginResponse{Authorized: true}, nil
}

func (h *handler) GetUser(c context.Context, req *genproto.GetUserRequest) (*genproto.User, error) {
	m := &model.User{}

	if err := bindJSON(req, m); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	user, err := h.userStore.FindUserBy(m)
	if user == nil {
		return nil, status.Errorf(codes.NotFound, "Username not found")
	}

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	// Bind response
	res := &genproto.User{}
	if err := bindJSON(user, res); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return res, nil
}
