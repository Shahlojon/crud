package security

import (
	"time"
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
	"errors"
	"context"
	"log"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/jackc/pgx/v4"
)

var ErrNoSuchUser = errors.New("no such user")
var ErrInvalidPassword = errors.New("invalid password")
var ErrInternal = errors.New("invalid error")
var ErrExpireToken=errors.New("token expired")
type Service struct {
	db *pgxpool.Pool
}

func NewService(db *pgxpool.Pool) *Service{
	return &Service{db: db}
}

func (s *Service) Auth(login, password string) bool {
	sqlStatement :=`select login, password from managers where login=$1 and password=$2`

	err:=s.db.QueryRow(context.Background(),sqlStatement, login, password).Scan(&login, &password)
	if err!=nil {
		log.Print(err)
		return false
	}
	return true
}

// func (s *Service) TokenForCustomer(
// 	ctx context.Context,
// 	phone string,
// 	password string,
// )(token string, err error){
// 	var hash string
// 	var id int64
// 	err = s.db.QueryRow(ctx, `SELECT id, password FROM customers WHERE phone=$1`, phone).Scan(&password)
// 	if err == pgx.ErrNoRows{
// 		return "", ErrNoSuchUser
// 	}
// 	if err!=nil{
// 		return "", ErrInternal
// 	}

// 	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
// 	if err!=nil{
// 		return "",ErrInvalidPassword
// 	}

// 	buffer:=make([]byte, 256)
// 	n,err:=rand.Read(buffer)
// 	if n!=len(buffer)|| err!=nil{
// 		return "",ErrInternal
// 	}

// 	token=hex.EncodeToString(buffer)
// 	_,err=s.db.Exec(ctx, `INSERT INTO customer_tokens(token, customer_id) VALUES($1,$2)`, token, id)
// 	if err!=nil{
// 		return "",ErrInternal
// 	}

// 	return token, nil
// }

// func (s *Service) AuthenticateCustomer(
// 	ctx context.Context,
// 	token string,
// ) (id int64, err error) {
// 	err = s.db.QueryRow(ctx, `SELECT customer_id FROM customers_tokens WHERE token = $1`, token).Scan(&id)

// 	if err == pgx.ErrNoRows {
// 		return 0, ErrNoSuchUser
// 	}

// 	if err!=nil {
// 		return 0, ErrInternal
// 	}
// 	return id,nil
// }

//TokenForCustomer
func (s *Service) TokenForCustomer(ctx context.Context, phone, password string)(string, error){
	var hash string
	var id int64

	err := s.db.QueryRow(ctx, "select id, password from customers where phone = $1", phone).Scan(&id, &hash)

	if err == pgx.ErrNoRows{
		return "", ErrNoSuchUser
	}
	if err != nil{ 
		return "", ErrInternal
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil{
		return "", ErrInvalidPassword
	}

	buffer := make([]byte, 256)
	n, err := rand.Read(buffer)
	if n != len(buffer) || err != nil{
		return "", ErrInternal
	}

	token := hex.EncodeToString(buffer)
	_, err = s.db.Exec(ctx, "insert into customers_tokens(token, customer_id) values($1, $2)", token, id)
	if err != nil{
		return "", ErrInternal
	}

	return token, nil

}

func (s *Service) AuthenticateCustomer(ctx context.Context, tkn string)(int64, error){
	var id int64
	var expire time.Time
	err := s.db.QueryRow(ctx, "select customer_id, expire from customers_tokens where token=$1", tkn).Scan(&id, &expire)
	if err == pgx.ErrNoRows{
		return 0, ErrNoSuchUser
	}
	if err != nil{
		return 0, ErrInternal
	}

	tNow := time.Now().Format("2006-01-02 15:04:05")
	tEnd := expire.Format("2006-01-02 15:04:05")

	if tNow > tEnd {
		return 0, ErrExpireToken
	}

	return id, nil
}