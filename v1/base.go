// Copyright 2016 Tim Foerster <github@mailserver.1n3t.de>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"github.com/muhproductions/muh/helper"
	"github.com/muhproductions/muh/v1/resources"
	"gopkg.in/redis.v3"
	"os"
)

// Routes - Register all routes for API version 1
func Routes(api *gin.Engine) {
	version := api.Group("/v1")
	version.Use(Ratelimit())
	version.GET("/ping", Ping)

	resources.UserResource{
		Engine: version,
	}.Routes()

	resources.GistResource{
		Engine: version,
	}.Routes()

	go EventHandler(helper.RedisClient())

}

// EventHandler will subscribe to keyspace notifications
// and run callbacks.
func EventHandler(r *redis.Client) {
	db := "muh.db"
	if os.Getenv("DB") != "" {
		db = os.Getenv("DB")
	}
	b, err := bolt.Open(db, 0600, nil)
	if err != nil {
		panic(err)
	}
	helper.Bolt = b
	helper.BoltInit()
	pubsub, _ := r.Subscribe("__keyevent@0__:expired")
	for {
		msg, _ := pubsub.ReceiveMessage()
		for _, callback := range helper.Callbacks {
			callback(msg.Payload)
		}
	}
}

// Ping - a generic ping / pong route
func Ping(c *gin.Context) {
	c.JSON(418, gin.H{
		"message": "pong",
	})
}
