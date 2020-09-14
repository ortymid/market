package main

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/ortymid/t3-grpc/grpc"
	httpserver "github.com/ortymid/t3-grpc/http"
	"github.com/ortymid/t3-grpc/market"
	httpservice "github.com/ortymid/t3-grpc/service/http"
	"github.com/ortymid/t3-grpc/service/mem"
)

type Config struct {
	Port           int
	JWTAlg         string
	JWTSecret      interface{}
	UserServiceURL string
}

func main() {
	config := getConfig()

	userService := httpservice.NewUserService(config.UserServiceURL)
	productService := mem.NewProductService()
	m := &market.Market{
		UserService:    userService,
		ProductService: productService,
	}
	grpcServer := &grpc.Server{
		Market: m,
	}
	go grpcServer.Run(8081)

	grpcMarket, err := grpc.NewMarket(":8081")
	if err != nil {
		panic(err)
	}

	httpserver.Run(config.Port, config.JWTAlg, config.JWTSecret, grpcMarket)
}

func getConfig() *Config {
	portString := getEnvDefault("PORT", "8080")
	port, err := strconv.Atoi(portString)
	if err != nil {
		panic("cannot read PORT: " + err.Error())
	}

	jwtAlg := getEnvDefault("JWT_ALG", "HS256")
	jwtSecret, err := getKey(os.Getenv("KEY_SERVICE_URL"))
	if err != nil {
		panic(fmt.Errorf("cannot get JWT secret: %w", err))
	}

	usURL := os.Getenv("USER_SERVICE_URL")

	return &Config{
		Port:           port,
		JWTAlg:         jwtAlg,
		JWTSecret:      jwtSecret,
		UserServiceURL: usURL,
	}
}

func getKey(url string) (*rsa.PublicKey, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("something went wrong")
	}

	key := &rsa.PublicKey{}
	err = json.NewDecoder(resp.Body).Decode(key)
	defer resp.Body.Close()

	return key, err
}

func getEnvDefault(key string, d string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		val = d
	}
	return val
}
