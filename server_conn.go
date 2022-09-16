//
// Copyright (c) 2018- yutopp (yutopp@gmail.com)
//
// Distributed under the Boost Software License, Version 1.0. (See accompanying
// file LICENSE_1_0.txt or copy at  https://www.boost.org/LICENSE_1_0.txt)
//

package rtmp

import (
	"log"

	"github.com/pkg/errors"

	"github.com/yutopp/go-rtmp/handshake"
)

// serverConn A wrapper of a connection. It prorives server-side specific features.
type serverConn struct {
	conn *Conn
}

func newServerConn(conn *Conn) *serverConn {
	return &serverConn{
		conn: conn,
	}
}

func (sc *serverConn) Serve() error {
	hs := handshake.HandshakeServer{}
	if err := hs.ReadC0C1(sc.conn.rwc); err != nil {
		return err
	}
	log.Printf("< R Handshake C0+C1.")

	log.Printf("> W Handshake S0+S1+S2.")
	if err := hs.WriteS0S1S2(sc.conn.rwc); err != nil {
		return err
	}

	if err := hs.ReadC2(sc.conn.rwc); err != nil {
		return err
	}
	log.Printf("< R Handshake C2.")

	// if err := handshake.PlainHandshakeWithClient(sc.conn.rwc, sc.conn.rwc, &handshake.Config{
	// 	SkipHandshakeVerification: sc.conn.config.SkipHandshakeVerification,
	// }); err != nil {
	// 	return errors.Wrap(err, "Failed to handshake")
	// }

	ctrlStream, err := sc.conn.streams.Create(ControlStreamID)
	if err != nil {
		return errors.Wrap(err, "Failed to create control stream")
	}
	ctrlStream.handler.ChangeState(streamStateServerNotConnected)

	sc.conn.streamer.controlStreamWriter = ctrlStream.Write

	if sc.conn.handler != nil {
		sc.conn.handler.OnServe(sc.conn)
	}

	return sc.conn.handleMessageLoop()
}

func (sc *serverConn) Close() error {
	return sc.conn.Close()
}
