package env

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"product-service/src/common/log"
)

type currentApp struct {
	RestfulAddr   string
	GrpcAddr      string
	GrpcBasicAuth string
}

type userService struct {
	GrpcAddr string
	GrpcAuth string
}

type postgres struct {
	DSN string
}

type jwt struct {
	PublicKey *rsa.PublicKey
}

type imagekit struct {
	PublicKey   string
	PrivateKey  string
	UrlEndpoint string
}

type kafka struct {
	Addr1 string
	Addr2 string
	Addr3 string
}

type Config struct {
	CurrentApp  *currentApp
	Postgres    *postgres
	Jwt         *jwt
	ImageKit    *imagekit
	Kafka       *kafka
	UserService *userService
}

var Conf *Config

func Load() {
	if os.Getenv("MODE") == "NON-DEV" {
		LoadFromVault()
	}

	currentAppConf := new(currentApp)
	currentAppConf.RestfulAddr = os.Getenv("RESTFUL_ADDR")
	currentAppConf.GrpcAddr = os.Getenv("GRPC_ADDR")
	currentAppConf.GrpcBasicAuth = os.Getenv("GRPC_BASIC_AUTH")

	postgresConf := new(postgres)
	postgresConf.DSN = os.Getenv("POSTGRES_DSN")

	jwtConf := new(jwt)
	base64Str, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_PUBLIC_KEY"))
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	jwtConf.PublicKey = loadRsaPublicKey(string(base64Str))

	imagekitConf := new(imagekit)
	imagekitConf.PublicKey = os.Getenv("IMAGEKIT_PUBLIC_KEY")
	imagekitConf.PrivateKey = os.Getenv("IMAGEKIT_PRIVATE_KEY")
	imagekitConf.UrlEndpoint = os.Getenv("IMAGEKIT_URL_ENDPOINT")

	kafkaConf := new(kafka)
	kafkaConf.Addr1 = os.Getenv("KAFKA_ADDR_1")
	kafkaConf.Addr2 = os.Getenv("KAFKA_ADDR_2")
	kafkaConf.Addr3 = os.Getenv("KAFKA_ADDR_3")

	userServiceConf := new(userService)
	userServiceConf.GrpcAddr = os.Getenv("USER_SERVICE_GRPC_ADDR")
	userServiceConf.GrpcAuth = os.Getenv("USER_SERVICE_GRPC_AUTH")

	Conf = &Config{
		CurrentApp:  currentAppConf,
		Postgres:    postgresConf,
		Jwt:         jwtConf,
		ImageKit:    imagekitConf,
		Kafka:       kafkaConf,
		UserService: userServiceConf,
	}
}

func loadRsaPublicKey(publicKeyStr string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(publicKeyStr))
	if block == nil {
		log.Logger.Fatal("failed to parse public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		log.Logger.Fatal(err.Error())
	}

	rsaPublicKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		log.Logger.Fatal(err.Error())
	}

	return rsaPublicKey
}
