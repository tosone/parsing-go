package interact

import (
	"context"
	"testing"

	"smartconn.cc/tosone/ra-plus/common/util"
	"smartconn.cc/tosone/ra-plus/drivers/audio"
	"smartconn.cc/tosone/ra-plus/store"
)

func TestNew(t *testing.T) {
	var err error
	var interact *State

	if err = store.Initialize(); err != nil {
		t.Fatal(err)
	}

	if err = audio.Startup(); err != nil {
		t.Fatal(err)
	}

	if interact, err = New(); err != nil {
		t.Fatal(err)
	}

	if err = interact.Setting(util.UUID(), "c249b2cc-e19c-4c10-bdc7-780aa1d70371", "59a7836abbfd1856ece6b2f7"); err != nil {
		t.Fatal(err)
	}

	if err = interact.Run(context.Background(), "c249b2cc-e19c-4c10-bdc7-780aa1d70371"); err != nil {
		t.Fatal(err)
	}

	interact.Destroy()
}
