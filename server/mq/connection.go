package mq

import "github.com/streadway/amqp"

// getNewConn 获取一个新的Connection
func (mq *RabbitMQ) getNewConn() (*amqp.Connection, error) {
	var err error
	mq.conn, err = amqp.Dial(mq.uri)
	return mq.conn, err
}

// GetConn 获取Connection，优先获取缓存Connection，缓存过期则新建Connection，并更新缓存
func (mq *RabbitMQ) getConn() error {
	if mq.conn != nil && !mq.conn.IsClosed() {
		return nil
	}

	var e error
	mq.conn, e = mq.getNewConn()
	return e
}

// Close 关闭Connection
func (mq *RabbitMQ) Close() {
	if mq.conn != nil && !mq.conn.IsClosed() {
		mq.conn.Close()
	}
	mq.publishChan = nil
	mq.consumeChan = nil
}
