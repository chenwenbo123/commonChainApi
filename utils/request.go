package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Get(url, apikey string) ([]byte, error) {
	defer func() {
		if error := recover(); error != nil {

		}
	}()
	req, error := http.NewRequest("GET", url, nil)
	if error != nil {
		return []byte("error"), error
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("TRON_PRO_API_KEY", apikey)
	res, error1 := http.DefaultClient.Do(req)
	if error1 != nil {
		return []byte("error"), error1
	}
	defer res.Body.Close()
	body, error2 := ioutil.ReadAll(res.Body)
	if error2 != nil {
		return []byte("error"), error2
	}
	return body, nil
}

func Post(url string, data map[string]interface{}) ([]byte, error) {
	defer func() {
		if error := recover(); error != nil {

		}
	}()
	payload, error := json.Marshal(data)
	if error != nil {
		return nil, errors.New("error")
	}
	req, error := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if error != nil {
		return nil, errors.New("error")
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body, nil
}

func TronPost(url, apikey string, data map[string]interface{}) ([]byte, error) {
	defer func() {
		if error := recover(); error != nil {

		}
	}()
	payload, error := json.Marshal(data)
	if error != nil {
		return nil, errors.New("error")
	}
	req, error := http.NewRequest("POST", url, strings.NewReader(string(payload)))
	if error != nil {
		return nil, errors.New("error")
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("TRON_PRO_API_KEY", apikey)
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	return body, nil
}

func IdomaticGet(url string) ([]byte, error) {
	defer func() {
		if error := recover(); error != nil {

		}
	}()
	req, error := http.NewRequest("GET", url, nil)
	if error != nil {
		return []byte("error"), error
	}
	req.Header.Add("accept", "application/json")
	res, error1 := http.DefaultClient.Do(req)
	if error1 != nil {
		return []byte("error"), error1
	}
	defer res.Body.Close()
	body, error2 := ioutil.ReadAll(res.Body)
	if error2 != nil {
		return []byte("error"), error2
	}
	return body, nil
}

func IdomaticTronGet(url, apikey string) ([]byte, error) {
	defer func() {
		if error := recover(); error != nil {

		}
	}()
	req, error := http.NewRequest("GET", url, nil)
	if error != nil {
		return []byte("error"), errors.New("Http Request Error")
	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("TRON_PRO_API_KEY", apikey)
	res, error1 := http.DefaultClient.Do(req)
	if error1 != nil {
		return []byte("error"), errors.New("Http Request Error")
	}
	defer res.Body.Close()
	body, error2 := ioutil.ReadAll(res.Body)
	if error2 != nil {
		return []byte("error"), errors.New("Http Request Error")
	}
	var result map[string]interface{}
	json.Unmarshal(body, &result)
	return body, nil
}

// g *sync.WaitGroup,
func SendInform(urls string, chain, blockNum int64, coinName, contractAddress, typeInfo, fromAddress, toAddress string, num float64, txid string) {
	params := url.Values{}
	Url, err := url.Parse(urls)
	if err != nil {
		return
	}
	params.Set("chain", fmt.Sprint(chain))
	params.Set("block_num", fmt.Sprint(blockNum))
	params.Set("coin_name", fmt.Sprint(coinName))
	params.Set("contract_address", contractAddress)
	params.Set("type", typeInfo)
	params.Set("from_address", fromAddress)
	params.Set("to_address", toAddress)
	params.Set("num", fmt.Sprint(num))
	params.Set("txid", txid)
	//如果参数中有中文参数,这个方法会进行URLEncode
	Url.RawQuery = params.Encode()
	urlPath := Url.String()
	resp, err := http.Get(urlPath)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	//g.Done()
}
