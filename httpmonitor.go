package httpmonitor

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type Monitor struct {
	ToEmail   string
	FromEmail string
	Password  string
	Services  []Service
}

type Service struct {
	Address   string
	Type      string
	Body      string
	BodyType  string
	Frequency time.Duration
	DownSince time.Time
}

func (m *Monitor) Run() error {
	errCh := make(chan error)
	for _, s := range m.Services {
		go func(s Service) {
			for {
				_, err := s.Check()
				if err != nil {
					fmt.Println(s.Address, "is down, sending e-mail")
					err := m.SendNotification(&s, err.Error())
					if err != nil {
						errCh <- err
					}
				} else {
					fmt.Println(s.Address, "is alive")
				}
				time.Sleep(s.Frequency)
			}
		}(s)
	}
	err := <-errCh
	return err
}

func (m *Monitor) SendNotification(s *Service, message string) error {
	from := m.FromEmail
	pass := m.Password
	to := m.ToEmail

	if s.DownSince.IsZero() {
		s.DownSince = time.Now()
	}

	timestring := s.DownSince.Format("2.1.2006 15:04:05")

	msg := []byte("To: " + m.ToEmail + "\r\n" +
		"Subject: Service " + s.Address + " is down!\r\n" +
		"\r\n" +
		"It looks like the service " + s.Address + " is down,\r\n" +
		"It has been down since " + timestring + ".\r\n" +
		"Error message:\r\n" +
		message)

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	return errors.Wrap(err, "could not send e-mail notification")

}

func (s *Service) Check() (string, error) {
	url := "http://" + s.Address
	var resp *http.Response
	var err error
	if s.Type == "GET" {
		resp, err = http.Get(url)
		if err != nil {
			return "", errors.Wrap(err, "get failed")
		}

	} else {
		postArgs := strings.NewReader(string(s.Body))
		resp, err = http.Post(url, s.BodyType, postArgs)
		if err != nil {
			return "", errors.Wrap(err, "post failed")
		}
	}

	if resp.StatusCode != 200 {
		return "", errors.New("HTTP status code != 200")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "reading body failed")
	}

	if len(body) <= 0 {
		return "", errors.Wrap(err, "empty response body")
	}

	s.DownSince = time.Time{}

	return string(body), nil
}
