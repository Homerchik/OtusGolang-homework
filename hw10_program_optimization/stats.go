package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

type User struct {
	ID int
	//nolint:staticcheck,tagliatelle
	Email string `json:"Email,nocopy"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

func getUsers(r io.Reader) ([]string, error) {
	br := bufio.NewReader(r)
	delimiter := byte(0x0A)
	var (
		user User
		emails []string
	)
	for {
		line, err := br.ReadBytes(delimiter)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if jsonErr := easyjson.Unmarshal(line, &user); jsonErr != nil {
			return nil, err
		}
		emails = append(emails, user.Email)
		if err == io.EOF {
			return emails, nil
		}
	}
}

func countDomains(emails []string, domain string) (DomainStat, error) {
	result := make(DomainStat)
	re, err := regexp.Compile(`\.` + domain)
	if err != nil {
		return nil, err
	}
	for _, email := range emails {
		if re.Match([]byte(email)) {
			fullDomain := strings.SplitN(email, "@", 2)
			result[strings.ToLower(fullDomain[1])]++
		}
	}
	return result, nil
}
