package msgq

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestQueue(t *testing.T) {
	msqid, err := OpenQueue(0)
	require.NoError(t, err)
	require.NotZero(t, msqid)

	m := NewMsg(1, []byte("1"))
	require.NotNil(t, m)

	err2 := SndMsg(msqid, m)
	require.NoError(t, err2)

	m1 := NewMsg(1, []byte("2"))
	err3 := RcvMsg(msqid, m1)
	require.NoError(t, err3)
	require.Equal(t, m1, m)

	err4 := DestroyQueue(msqid)
	require.NoError(t, err4)
}

func TestMsgChan(t *testing.T) {
	msqid, err := OpenQueue(0)
	require.NoError(t, err)
	require.NotZero(t, msqid)

	m := NewMsg(1, []byte("1"))
	require.NotNil(t, m)

	go func() {
    time.Sleep(time.Second)
		err2 := SndMsg(msqid, m)
		require.NoError(t, err2)
	}()

	m1 := NewMsg(1, []byte("2"))
	<-MsgChan(msqid, m1)
	require.Equal(t, m1, m)

	err4 := DestroyQueue(msqid)
	require.NoError(t, err4)
}
