package authorization

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/qwuiemme/ellipsespace-server/pkg/client"
	"golang.org/x/crypto/bcrypt"
)

type SessionBase struct {
	Id          int    `json:"id"`
	SessionName string `json:"sname"`
	AccessLevel int8   `json:"access-level"`
}

type Session struct {
	SessionBase
	Password string `json:"password"`
}

func Unmarshal(r io.Reader) (*Session, error) {
	jsonByte, err := io.ReadAll(r)

	if err != nil {
		return &Session{}, err
	}

	var obj Session
	err = json.Unmarshal(jsonByte, &obj)

	if err != nil {
		return &Session{}, err
	}

	return &obj, nil
}

func GetSession(sessionName string) (s Session, err error) {
	conn := client.Connect()
	defer conn.Close()

	res, err := conn.Query(fmt.Sprintf("SELECT * FROM `sessions` WHERE SessionName = '%s'", sessionName))

	if err != nil {
		return Session{}, err
	}

	defer res.Close()

	for res.Next() {
		err = res.Scan(&s.Id, &s.SessionName, &s.Password, &s.AccessLevel)

		if err != nil {
			return Session{}, err
		}
	}

	return
}

func (s *Session) AddToDatabase() error {
	conn := client.Connect()
	defer conn.Close()

	res, err := conn.Query(fmt.Sprintf("INSERT INTO `sessions` VALUES (0, '%s', '%s', 0)", s.SessionName, s.Password))

	if err != nil {
		return err
	} else {
		defer res.Close()
		return nil
	}
}

func (s *Session) AddFullToDatabase() error {
	conn := client.Connect()
	defer conn.Close()

	res, err := conn.Query(fmt.Sprintf("INSERT INTO `sessions` VALUES (0, '%s', '%s', '%s')", s.SessionName, s.Password, strconv.Itoa(int(s.AccessLevel))))

	if err != nil {
		return err
	} else {
		defer res.Close()
		return nil
	}
}

func (s *Session) Update() error {
	conn := client.Connect()
	defer conn.Close()

	res, err := conn.Query(fmt.Sprintf("UPDATE `sessions` SET SessionName = '%s', Password = '%s' WHERE Id = '%s'", s.SessionName, s.Password, strconv.Itoa(int(s.AccessLevel))))

	if err != nil {
		return err
	}

	defer res.Close()

	return nil
}

func (s *Session) UpdateAll() error {
	conn := client.Connect()
	defer conn.Close()

	res, err := conn.Query(fmt.Sprintf("UPDATE `sessions` SET SessionName = '%s', Password = '%s', AccessLevel = '%s' WHERE Id = '%s'", s.SessionName, s.Password, strconv.Itoa(int(s.AccessLevel)), strconv.Itoa(s.Id)))

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

func (s *Session) ComparePassword(input string) bool {
	return bcrypt.CompareHashAndPassword([]byte(s.Password), []byte(input)) == nil
}
