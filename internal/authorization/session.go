package authorization

import (
	"fmt"
	"strconv"

	"github.com/qwuiemme/ellipsespace-server/pkg/client"
	"golang.org/x/crypto/bcrypt"
)

type Session struct {
	Id          int
	Login       string `json:"login"`
	Password    string `json:"password"`
	AccessLevel int8   `json:"access-level"`
}

func GetSession(name string) (s Session, err error) {
	conn := client.Connect()
	defer conn.Close()

	res, err := conn.Query(fmt.Sprintf("GET * FROM `sessions` WHERE Name = '%s'", name))

	if err != nil {
		return Session{}, err
	}

	defer res.Close()

	for res.Next() {
		err = res.Scan(&s.Id, &s.Login, &s.Password, &s.AccessLevel)

		if err != nil {
			return Session{}, err
		}
	}

	return
}

func (s *Session) Update() error {
	conn := client.Connect()
	defer conn.Close()

	res, err := conn.Query(fmt.Sprintf("UPDATE `session` SET Login = '%s', Password = '%s', AccessLevel = '%s' WHERE Id = '%s'", s.Login, s.Password, strconv.Itoa(int(s.AccessLevel)), strconv.Itoa(s.Id)))

	if err != nil {
		return err
	}

	defer res.Close()

	return nil
}

func (s *Session) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)

	if err != nil {
		return err
	}

	s.Password = string(hash)
	s.Update()
	return nil
}

func (s *Session) ComparePassword(input string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(input))

	if err != nil {
		return false, err
	}

	return true, nil
}
