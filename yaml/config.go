package Config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

const path = "./config.yaml"

type Conf struct {
	Tron struct {
		ApiKey string `yaml:"apiKey"`
	} `yaml:"Tron"`

	Db struct {
		Host     string `yaml:"host"`
		DbName   string `yaml:"dbName"`
		UserName string `yaml:"userName"`
		Password string `yaml:"password"`
	} `yaml:"Db"`

	Node struct {
		Ethereum string `yaml:"ethereum"`
		Bsc      string `yaml:"bsc"`
		Tron     string `yaml:"tron"`
	} `yaml:"Node"`

	Inform struct {
		Url string `yaml:"url"`
	} `yaml:"Inform"`

	System struct {
		Port   string `yaml:"port"`
		Router string `yaml:"router"`
	} `yaml:"System"`
}

func LoadConfig() *Conf {
	config := new(Conf)
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		fmt.Println(err)
	}
	return config
}

func SaveConfig(c *Conf) {
	data, err := yaml.Marshal(c)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(path, data, 0777)
	if err != nil {
		panic("错误")
	}
	//fmt.Println("保存成功")
}
