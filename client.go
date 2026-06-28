package books

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) GetBook(ID string) (Book, error) {
	res, err := http.Get("http://" + c.addr + "/v1/find/" + ID)

	if err != nil {
		fmt.Printf("Request: %v\n", err)
		return Book{}, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Reading response: %v\n", err)
		return Book{}, err
	}

	book := Book{}
	err = json.Unmarshal(body, &book)

	if err != nil {
		fmt.Printf("Unmarshalling response: %v\n", err)
		return Book{}, err
	}

	return book, nil
}

func (c *Client) GetAllBooks() ([]Book, error) {
	res, err := http.Get("http://" + c.addr + "/v1/list")

	if err != nil {
		fmt.Printf("Request: %v\n", err)
		return []Book{}, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Reading response: %v\n", err)
		return []Book{}, err
	}

	list := []Book{}
	err = json.Unmarshal(body, &list)

	if err != nil {
		fmt.Printf("Unmarshalling response: %v\n", err)
		return []Book{}, err
	}

	return list, nil
}
