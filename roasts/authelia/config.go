package authelia

import (
	"bytes"
	"database/sql"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"testing"

	"coldbrew.go/cb/cmd/db/create"
	"coldbrew.go/cb/common"
	"github.com/sirupsen/logrus"
)

func EncryptPasswords(password string) (string, error) {
	// TODO pipe a password into it (ideal but not sure how) or pull up the memo in this function
	command := fmt.Sprintf("run authelia/authelia:latest authelia crypto hash generate argon2 --password %v", password)
	args := strings.Split(command, " ")
	r, w, _ := os.Pipe()
	os.Stdout = w

	passwordCmd := exec.Command("docker", args...)
	passwordCmd.Stdout = os.Stdout
	passwordCmd.Stderr = os.Stderr
	err := passwordCmd.Run()

	outC := make(chan string)
	// TODO MAKE THIS MY OWN
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// back to normal state
	w.Close()
	out := <-outC

	return out, err
}

func Preinstall() {
	// TODO make this an easily callable flow, the use case of needing custom variable
	// insertion seems like it will be common
	password, err := EncryptPasswords(common.PromptInput("password"))

	if err != nil {
		logrus.Error(err)
	}

	// TODO still don't hardcode the db location
	db, err := sql.Open("sqlite3", "./InstallDir/data.db")

	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
	defer db.Close()

	err = create.CreatePortsTable(db)

	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	err = create.CreateTable(db, "authelia")

	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}

	errs := create.Seed(db, map[string]string{"password": password}, "authelia")

	for _, err := range errs {
		logrus.Error(err)
	}
}

func TestEncryptPasswords(t *testing.T) {
	text, err := EncryptPasswords("OTFLXGgnZ3xN9qV0zJNdC8YCILjGmqHE5fI9DZ7CgrnI85oGNnrhPeaobXabOTuMw6EH6RZElSlqtvvOaCKVDJD2paPmEUgmzRbFtQwbmz85r6oSz479Xuq9x1r5McxT")
	if err != nil {
		t.Error(err)
	}
	logrus.Info(text)
}
