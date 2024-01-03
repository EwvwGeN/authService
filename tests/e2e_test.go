package tests

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/EwvwGeN/authService/internal/config"
	authProto "github.com/EwvwGeN/authService/proto/gen/go"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

func TestSuiteRun(t *testing.T) {
	suite.Run(t, new(testSuite))
}

type testSuite struct {
	suite.Suite
	cfg        *config.Config
	grpcC      authProto.AuthClient
	clientConn *grpc.ClientConn
	cancelFunc func()
}

func (suite *testSuite) SetupSuite() {
	cfg, err := config.LoadConfig("./configs/test_config.yaml")
	suite.Require().NoError(err)
	cancel, err := prepareContainers(cfg)
	suite.cancelFunc = cancel
	suite.Require().NoError(err)
	clientConn, err := grpc.DialContext(context.Background(),
		fmt.Sprintf("localhost:%d", cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	suite.Require().NoError(err)
	suite.cancelFunc = cancel
	suite.cfg = cfg
	suite.clientConn = clientConn
	suite.grpcC = authProto.NewAuthClient(clientConn)
}

func (suite *testSuite) TearDownSuite() {
	suite.cancelFunc()
	suite.clientConn.Close()
}

func (suite *testSuite) TestRegisterHappyPass() {
	email := gofakeit.Email()
	pwd := gofakeit.Password(true, true, true, true, false, 10)
	resReg, err := suite.grpcC.Register(context.Background(), &authProto.RegisterRequest{
		Email: email,
		Password: pwd,
	})
	suite.Require().NoError(err)
	suite.NotEmpty(resReg.GetUserId())
}

func (suite *testSuite) TestRegisterIncorrectValues() {
	email := gofakeit.Email()
	pwd := "1"
	resReg, err := suite.grpcC.Register(context.Background(), &authProto.RegisterRequest{
		Email: email,
		Password: pwd,
	})
	suite.Equal(codes.InvalidArgument, status.Convert(err).Code())
	suite.Empty(resReg)
}

func (suite *testSuite) TestRegisterUserAlreadyExist() {
	email := gofakeit.Email()
	pwd := gofakeit.Password(true, true, true, true, false, 10)

	resReg, err := suite.grpcC.Register(context.Background(), &authProto.RegisterRequest{
		Email: email,
		Password: pwd,
	})
	suite.Require().NoError(err)
	suite.Require().NotEmpty(resReg.GetUserId())
	
	resRegTwo, err := suite.grpcC.Register(context.Background(), &authProto.RegisterRequest{
		Email: email,
		Password: pwd,
	})
	suite.Equal(codes.AlreadyExists, status.Convert(err).Code())
	suite.Empty(resRegTwo)
	
}

func (suite *testSuite) TestLoginHappyPass() {
	email := gofakeit.Email()
	pwd := gofakeit.Password(true, true, true, true, false, 10)
	resReg, err := suite.grpcC.Register(context.Background(), &authProto.RegisterRequest{
		Email: email,
		Password: pwd,
	})
	suite.Require().NoError(err)
	suite.NotEmpty(resReg.GetUserId())
	log.Println(CorrectApp)
	resLog, err := suite.grpcC.Login(context.Background(), &authProto.LoginRequest{
		Email:    email,
		Password: pwd,
		AppId:    CorrectApp.Id,
	})
	createTime := time.Now()
	suite.Require().NoError(err)
	token := resLog.GetToken()
	suite.NotEmpty(token)
	parsedToken, _ := jwt.Parse(token, nil)
	suite.Require().NotEmpty(parsedToken)
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	suite.Require().True(ok)
	suite.Equal(resReg.GetUserId(), claims["user_id"].(string))
	suite.Equal(email, claims["email"].(string))
	suite.Equal(CorrectApp.Id, claims["app_id"].(string))
	suite.InDelta(createTime.Add(suite.cfg.TokenTTL).Unix(), int64(claims["exp"].(float64)), 1)
}

func (suite *testSuite) TestLoginIncorrectUserEmail() {
	email := gofakeit.Email()
	pwd := gofakeit.Password(true, true, true, true, false, 10)
	resReg, err := suite.grpcC.Register(context.Background(), &authProto.RegisterRequest{
		Email: email,
		Password: pwd,
	})
	suite.Require().NoError(err)
	suite.NotEmpty(resReg.GetUserId())
	resLog, err := suite.grpcC.Login(context.Background(), &authProto.LoginRequest{
		Email:    "",
		Password: pwd,
		AppId:    "",
	})
	suite.Equal(codes.InvalidArgument, status.Convert(err).Code())
	suite.Empty(resLog)
}

func (suite *testSuite) TestLoginIncorrectUserPassword() {
	email := gofakeit.Email()
	pwd := gofakeit.Password(true, true, true, true, false, 10)
	resReg, err := suite.grpcC.Register(context.Background(), &authProto.RegisterRequest{
		Email: email,
		Password: pwd,
	})
	suite.Require().NoError(err)
	suite.NotEmpty(resReg.GetUserId())
	resLog, err := suite.grpcC.Login(context.Background(), &authProto.LoginRequest{
		Email:    email,
		Password: "",
		AppId:    "",
	})
	suite.Equal(codes.InvalidArgument, status.Convert(err).Code())
	suite.Empty(resLog)
}

func (suite *testSuite) TestLoginIncorrectAppId() {
	email := gofakeit.Email()
	pwd := gofakeit.Password(true, true, true, true, false, 10)
	resReg, err := suite.grpcC.Register(context.Background(), &authProto.RegisterRequest{
		Email: email,
		Password: pwd,
	})
	suite.Require().NoError(err)
	suite.NotEmpty(resReg.GetUserId())
	resLog, err := suite.grpcC.Login(context.Background(), &authProto.LoginRequest{
		Email:    email,
		Password: pwd,
		AppId:    "",
	})
	suite.Equal(codes.InvalidArgument, status.Convert(err).Code())
	suite.Empty(resLog)
}

func (suite *testSuite) TestIsAdminTrueHappyPass() {
}

func (suite *testSuite) TestIsAdminFalseHappyPass() {
	email := gofakeit.Email()
	pwd := gofakeit.Password(true, true, true, true, false, 10)
	resReg, err := suite.grpcC.Register(context.Background(), &authProto.RegisterRequest{
		Email: email,
		Password: pwd,
	})
	suite.Require().NoError(err)
	suite.NotEmpty(resReg.GetUserId())
	resCheck, err := suite.grpcC.IsAdmin(context.Background(), &authProto.IsAdminRequest{
		UserId: resReg.GetUserId(),
	})
	suite.Require().NoError(err)
	suite.False(resCheck.GetIsAdmin())
}

func (suite *testSuite) TestIsAdminIncorrestUserId() {
	email := gofakeit.Email()
	pwd := gofakeit.Password(true, true, true, true, false, 10)
	resReg, err := suite.grpcC.Register(context.Background(), &authProto.RegisterRequest{
		Email: email,
		Password: pwd,
	})
	suite.Require().NoError(err)
	suite.NotEmpty(resReg.GetUserId())
	resCheck, err := suite.grpcC.IsAdmin(context.Background(), &authProto.IsAdminRequest{
		UserId: "",
	})
	suite.Equal(codes.InvalidArgument, status.Convert(err).Code())
	suite.Empty(resCheck)
}

func (suite *testSuite) TestIsAdminUserNotFound() {
	resCheck, err := suite.grpcC.IsAdmin(context.Background(), &authProto.IsAdminRequest{
		UserId: "111111111111111111111111",
	})
	suite.Equal(codes.NotFound, status.Convert(err).Code())
	suite.Empty(resCheck)
}
