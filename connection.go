//
// Copyright © 2011-2012 Guy M. Allard
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
//

package stompngo

import (
	"log"
)

// Exported Connection methods

/*
	Connected returns the current connection status.
*/
func (c *Connection) Connected() bool {
	return c.connected
}

/*
	Session returns the broker assigned session id.
*/
func (c *Connection) Session() string {
	return c.session
}

/*
	Protocol returns the current connection protocol level.
*/
func (c *Connection) Protocol() string {
	return c.protocol
}

/*
	SetLogger enables a client defined logger for this connection.

	Set to "nil" to disable logging.

	Example:
		// Start logging
		l := log.New(os.Stdout, "", log.Ldate|log.Lmicroseconds)
		c.SetLogger(l)
*/
func (c *Connection) SetLogger(l *log.Logger) {
	c.logger = l
}

/*
	SendTickerInterval returns any heartbeat send ticker interval in ms.  A return 
	value of zero means	no heartbeats are being sent.
*/
func (c *Connection) SendTickerInterval() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.sti / 1000000
}

/*
	ReceiveTickerInterval returns any heartbeat receive ticker interval in ms.  
	A return value of zero means no heartbeats are being received.
*/
func (c *Connection) ReceiveTickerInterval() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.rti / 1000000
}

/*
	SendTickerCount returns any heartbeat send ticker count.  A return value of
	zero usually indicates no send heartbeats are enabled.
*/
func (c *Connection) SendTickerCount() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.sc
}

/*
	ReceiveTickerCount returns any heartbeat receive ticker count. A return
	value of zero usually indicates no read heartbeats are enabled.
*/
func (c *Connection) ReceiveTickerCount() int64 {
	if c.hbd == nil {
		return 0
	}
	return c.hbd.rc
}

// Package exported functions

/*
	Supported checks if a particular STOMP version is supported in the current 
	implementation.
*/
func Supported(v string) bool {
	return supported.Supported(v)
}

// Unexported Connection methods

/*
	Log data if possible.
*/
func (c *Connection) log(v ...interface{}) {
	if c.logger == nil {
		return
	}
	c.logger.Print(c.session, v)
	return
}

/*
	Shutdown logic.
*/
func (c *Connection) shutdown() {
	// Shutdown heartbeats if necessary
	if c.hbd != nil {
		if c.hbd.hbs {
			c.hbd.ssd <- true
		}
		if c.hbd.hbr {
			c.hbd.rsd <- true
		}
	}
	// Stop writer go routine
	c.wsd <- true
	// We are not connected
	c.connected = false
	return
}

/*
	Read error handler.
*/
func (c *Connection) handleReadError(md MessageData) {
	// Notify any general subscriber of error
	c.input <- md
	// Notify all individual subscribers of error
	c.subsLock.Lock()
	for key := range c.subs {
		c.subs[key] <- md
	}
	c.subsLock.Unlock()
	// Let further shutdown logic proceed normally.
	return
}
