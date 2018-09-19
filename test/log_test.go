package test

import (
	_ "github.com/chrisho/mosquito/test/config"
	_ "github.com/chrisho/mosquito/alilog"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestAlilog(t *testing.T) {
	logrus.Println("Test 1")

	forever := make(chan int)
	<-forever

	t.Log("Test test done")
}