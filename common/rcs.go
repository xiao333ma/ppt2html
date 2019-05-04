package common

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/samuel/go-zookeeper/zk"
	"strings"
	"time"
)

/**
* RCS client
**/
type RcsClient struct {
	zkUrl        string
	clientMap    map[string]interface{}
	tokenTypeMap map[string]string
	zkConnect    *zk.Conn
}

/**
Get a standard redis client
**/
func (self RcsClient) GetRedisClient(token string) *redis.Client {
	client := self.clientMap[token]

	if client != nil {
		result, _ := client.(*redis.Client)
		return result
	}
	return self.createOrRefreshRedisClient(token)
}

/**
Get a sentinel client
**/
func (self RcsClient) GetRedisSentinel(token string) *redis.Client {
	client := self.clientMap[token]

	if client != nil {
		result, _ := client.(*redis.Client)
		return result
	}
	return self.createOrReFreshSentinelClient(token)
}

/**
Get a cluster client
**/
func (self RcsClient) GetRedisCluster(token string) *redis.ClusterClient {
	client := self.clientMap[token]
	if client != nil {
		result, _ := client.(*redis.ClusterClient)
		return result
	}
	return self.createOrRefreshClusterClient(token)
}

func (self RcsClient) watchNodeChange(ech <-chan zk.Event) {
	event := <-ech
	path := event.Path

	if event.Type.String() != "EventNodeDataChanged" {
		return
	}

	fmt.Println("Path: ", path, " changed")

	token := strings.Replace(path, "/wy-redis/app/", "", -1)
	clientType := self.tokenTypeMap[token]
	if "client" == clientType {
		self.createOrRefreshRedisClient(token)
	} else if "cluster" == clientType {
		self.createOrRefreshClusterClient(token)
	} else {
		self.createOrReFreshSentinelClient(token)
	}
}

func (self RcsClient) getAddressAndPorts(token string) []string {
	var path string = "/wy-redis/app/" + token
	body, _, ech, _ := self.zkConnect.GetW(path)
	content := string(body)

	var result []string = make([]string, 0, 5)
	for _, value := range strings.Split(content, "\n") {
		if len(strings.TrimSpace(value)) > 0 {
			result = append(result, value)
		}
	}

	go self.watchNodeChange(ech)
	return result
}

func (self RcsClient) createOrRefreshRedisClient(token string) *redis.Client {
	addressAndPorts := self.getAddressAndPorts(token)
	client := redis.NewClient(&redis.Options{
		Addr: addressAndPorts[0],
	})

	self.clientMap[token] = client
	self.tokenTypeMap[token] = "client"

	fmt.Println("A new standard client created. ")

	return client
}

func (self RcsClient) createOrReFreshSentinelClient(token string) *redis.Client {
	addressAndPorts := self.getAddressAndPorts(token)
	serviceName := strings.Split(addressAndPorts[0], ":")[0]

	client := redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:    serviceName,
		SentinelAddrs: addressAndPorts[1:],
	})
	self.clientMap[token] = client
	self.tokenTypeMap[token] = "sentinel"

	fmt.Println("A new sentinel client created. ")
	return client
}

func (self RcsClient) createOrRefreshClusterClient(token string) *redis.ClusterClient {
	addressAndPorts := self.getAddressAndPorts(token)
	client := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addressAndPorts,
	})
	self.clientMap[token] = client
	self.tokenTypeMap[token] = "cluster"

	fmt.Println("A new cluster client created. ")
	return client
}

/**
Create a RCS client
zkUrl : zkURL
*/
func CreateRcsClient(zkUrl string) *RcsClient {
	client := new(RcsClient)
	client.zkUrl = zkUrl
	client.clientMap = make(map[string]interface{}, 3)
	client.tokenTypeMap = make(map[string]string, 3)

	connect, _, error := zk.Connect(strings.Split(zkUrl, ","), time.Second*5)
	if error != nil {
		panic(error)
	}
	client.zkConnect = connect

	fmt.Println("RCS client created")

	return client
}
