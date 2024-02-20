package tests

import (
	"context"
	"fmt"
	"time"

	"github.com/EwvwGeN/authService/internal/config"
	"github.com/EwvwGeN/authService/internal/domain/models"
	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	CorrectApp = models.App{
		Id:	primitive.NewObjectIDFromTimestamp(time.Now()).Hex(),
		Name: "correct application",
		Secret: "_______8______16",
	}
	AdminUser models.User
)

func prepareContainers(cfg *config.Config) (UserСonfirmator func(email string) error, cancelFunc func(), err error) {
	mongoEnv, serverEnv := parseConfig(cfg)
    mongoCtx := context.Background()
	mongoC, err := testcontainers.GenericContainer(mongoCtx, testcontainers.GenericContainerRequest{
        ContainerRequest: testcontainers.ContainerRequest{
            Image: "mongo:7.0.4",
            WaitingFor: wait.ForListeningPort("27017/tcp"),
            Files: []testcontainers.ContainerFile{
                {
                    HostFilePath:      "./build/storage/mongo-init.sh",
                    ContainerFilePath: "/docker-entrypoint-initdb.d/mongo-init.sh",
                    FileMode:          0o777,
                },
            },
			Env: mongoEnv,
        },
        Started: false,
    })
    if err != nil {
		return
	}
	err = mongoC.Start(mongoCtx)
	if err != nil {
		return
	}

	_, _, err = mongoC.Exec(mongoCtx, []string{
		"mongosh",
		"-eval", fmt.Sprintf("use %s", cfg.MongoConfig.Database),
		"-eval", fmt.Sprintf(
			"db.%s.insertOne({_id: ObjectId(\"%s\"),name: \"%s\",secret: \"%s\"})",
			cfg.MongoConfig.AppCollection,
			CorrectApp.Id,
			CorrectApp.Name,
			CorrectApp.Secret,
		),
	})
	if err != nil {
		return
	}
	UserСonfirmator = func(email string) error {
		_, _, err := mongoC.Exec(mongoCtx, []string{
			fmt.Sprintf("mongosh"),
			"-eval", fmt.Sprintf("use %s", cfg.MongoConfig.Database),
			"-eval", fmt.Sprintf(
				"db.%s.updateOne({email: \"%s\"}, {\"$set\": {confirmed: true}})",
				cfg.MongoConfig.UserCollection,
				email,
			),
		})
		return err
	}
	mongoPorts, err := mongoC.Ports(mongoCtx)
	if err != nil {
		return
	}
	if len(mongoPorts) == 0 {
		err = fmt.Errorf("no mongo ports")
		return
	}
	mPort := mongoPorts["27017/tcp"][0].HostPort
	serverEnv["MONGO.DB_PORT"] = mPort
	serverEnv["MONGO.DB_HOST"] = "host.docker.internal"

	serverCtx := context.Background()
	serverC, err := testcontainers.GenericContainer(serverCtx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context: "../",
				Dockerfile: "./tests/build/Dockerfile",
			},
			ExposedPorts: []string{
				fmt.Sprintf("%s/tcp", cfg.GRPCConfig.Port),
			},
			WaitingFor: wait.ForExposedPort(),
			Env: serverEnv,
			HostConfigModifier: func(hc *container.HostConfig) {
				hc.ExtraHosts = append(hc.ExtraHosts, "host.docker.internal:host-gateway")
			},
		},
		Started: false,
	})
	if err != nil {
		return
	}

	err = serverC.Start(serverCtx)
	if err != nil {
		return
	}
	cancelFunc = func() {
		serverC.Terminate(serverCtx)
		mongoC.Terminate(mongoCtx)
    }

	serverPorts, err := serverC.Ports(serverCtx)
	if err != nil {
		return
	}
	port := serverPorts["44044/tcp"][0].HostPort
	cfg.GRPCConfig.Port = port
	
	return
}

func parseConfig(cfg *config.Config) (mongoEnv, serverEnv map[string]string) {
	mongoEnv = map[string]string{
		"MONGO_NEWUSER_NAME": cfg.MongoConfig.User,
		"MONGO_NEWUSER_PASSWORD": cfg.MongoConfig.Password,
		"MONGO_INITDB_NAME": cfg.MongoConfig.Database,
		"MONGO_INITDB_COL_USER": cfg.MongoConfig.UserCollection,
		"MONGO_INITDB_COL_APP": cfg.MongoConfig.AppCollection,
	}
	serverEnv = map[string]string{
		"LOG_LEVEL": cfg.LogLevel,
		"GRPC.PORT": cfg.GRPCConfig.Port,
		"HTTP.PORT": cfg.HttpConfig.Port,
		"VALIDATOR.EMAIL": cfg.Validator.EmailValidate,
		"VALIDATOR.PASSWORD": cfg.Validator.PasswordValidate,
		"VALIDATOR.USER_ID": cfg.Validator.UserIDValidate,
		"VALIDATOR.APP_ID": cfg.Validator.AppIDValidate,
		"MONGO.DB_CON_FORMAT": cfg.MongoConfig.ConectionFormat,
		"MONGO.DB_USER": cfg.MongoConfig.User,
		"MONGO.DB_PASS": cfg.MongoConfig.Password,
		"MONGO.DB_AUTH_SOURCE": cfg.MongoConfig.Database,
		"MONGO.DB_NAME": cfg.MongoConfig.Database,
		"MONGO.DB_COL_USER": cfg.MongoConfig.UserCollection,
		"MONGO.DB_COL_APP": cfg.MongoConfig.AppCollection,
		"RABBITMQ.PRODUCER.COUNT": fmt.Sprintf("%d",cfg.RabbitMQCfg.ProducerConfig.Count),
		"TOKEN_TTL": cfg.TokenTTL.String(),
	}
	return
}