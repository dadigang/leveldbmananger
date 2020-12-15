package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

func vds_database_put(key string, value string) {

	var err = db.Put([]byte(key), []byte(value), nil)
	if nil != err {
		fmt.Println(err)
	}

}

func cmdenv() {
	fmt.Println("() => Args -> ", os.Args)
	fmt.Println("() => NumCPU -> ", runtime.NumCPU())
	fmt.Println("() => NumGoroutine -> ", runtime.NumGoroutine())
	fmt.Println("() => GOOS -> ", runtime.GOOS)
	fmt.Println("() => GOARCH -> ", runtime.GOARCH)
	fmt.Println("() => Compiler -> ", runtime.Compiler)
	p, err := os.Executable()
	fmt.Println("() => os.Executable -> ", p, err)

	p, err = os.Getwd()
	if err != nil {
		panic(err)
	}
	fmt.Println("() => os.Getwd -> ", p, err)
	p, err = Home()
	fmt.Println("() => Home -> ", p, err)

}
func Home() (string, error) {
	user, err := user.Current()
	if nil == err {
		return user.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return homeWindows()
	}

	// Unix-like system, so just assume Unix
	return homeUnix()
}

func homeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
func isWindows() bool {
	return runtime.GOOS == "windows"
}
func isLinux() bool {
	return runtime.GOOS == "linux"
}
func isDarwin() bool {
	return runtime.GOOS == "darwin"
}
func isFreebsd() bool {
	return runtime.GOOS == "freebsd"
}
func vds_database_open() {
	var err error
	p, _ := os.Getwd()
	if len(os.Args) > 1 {
		p = os.Args[1]
	}
	db, err = leveldb.OpenFile(p, nil)
	if nil != err {
		fmt.Println(err)
		cmdq()
		return
	}
	fmt.Println("database opened")
}
func vds_database_get(key string) string {
	data, err := db.Get([]byte(key), nil)
	if nil != err {
		fmt.Println(err)
	}
	return string(data)
}
func vds_database_del(key string) {
	var err = db.Delete([]byte(key), nil)
	if nil != err {
		fmt.Println(err)
	}
}
func vds_database_close() {

	if nil != db {
		db.Close()
	}
}
func vds_database_is_open() bool {

	if nil != db {
		return true
	}
	return false
}
func vds_database_list(proc func(key string, value string) bool) {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		ok := proc(string(key), string(value))
		if !ok {
			break
		}
	}
	iter.Release()
	err := iter.Error()
	if nil != err {
		fmt.Println(err)
	}
}

func vds_database_count() int {
	count := 0
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		count++
	}
	iter.Release()
	err := iter.Error()
	if nil != err {
		fmt.Println(err)
	}
	return count
}

func vds_database_exist(key string) bool {
	_, err := db.Get([]byte(key), nil)
	if nil != err {
		//fmt.Println(err)
		return false
	}
	return true
}
func golang_io_readline(reader *bufio.Reader) string {
	text, _ := reader.ReadString('\n')
	text = strings.TrimSpace(text)
	// convert CRLF to LF
	text = strings.Replace(text, "\n", "", -1)
	return text
}
func readline(handler func(line string)) {
	reader := bufio.NewReader(os.Stdin)
	for {
		pwd, _ := os.Getwd()
		fmt.Print(pwd + " >> ")
		line := golang_io_readline(reader)
		if line == "" {
			continue
		}
		handler(line)

	}
}

func vdu_process_interaction_line(line string) bool {
	defer timeoutCheck("process line in ", time.Now())

	if "q" == line {
		cmdq()
		return true
	}
	if "env" == line {
		cmdenv()
		return true
	}

	if "list" == line {
		cmdlist()
		return true
	}
	if "count" == line {
		cmdcount()
		return true
	}
	if "stat" == line {
		cmdstat()
		return true
	}
	if strings.HasPrefix(line, "download ") {
		cmddownloads(strings.Fields(line)[1:])
		return true
	}
	if strings.HasPrefix(line, "rm ") {
		cmdrms(strings.Fields(line)[1:])
		return true
	}
	if strings.HasPrefix(line, "put ") {
		cmdputs(strings.Fields(line)[1:])
		return true
	}
	if strings.HasPrefix(line, "get ") {
		cmdgets(strings.Fields(line)[1:])
		return true
	}
	return false
}

func cmdlist() {
	vds_database_list(func(key string, value string) bool {

		fmt.Println(key + " <==> " + value)
		return true
	})
}

func cmdstat() {
	fmt.Println(vds_database_is_open())
}

func cmdcount() {
	count := vds_database_count()
	fmt.Println("total", count, "in database")
}

func cmddownloads(i []string) {

}
func cmdrms(i []string) {
	for _, str := range i {
		vds_database_del(str)
		fmt.Println(str,"deleted")
	}

}
func cmdputs(i []string) {
	key :=i[0]
	value :=strings.Join(i[1:]," ")
	vds_database_put(key,value)
}
func cmdgets(i []string) {
	for _, str := range i {
		value:=vds_database_get(str)
		fmt.Println(value)
	}
}
func cmdq() {
	vds_database_close()
	os.Exit(0)
}
func cmdexit() {
	vds_database_close()
	os.Exit(0)
}
func timeoutCheck(tag string, start time.Time) {
	dis := time.Since(start).Milliseconds()
	fmt.Println(tag, dis, "ms")
}
func main() {
	fmt.Println("20201213")

	cmdenv()
	vds_database_open()

	readline(func(line string) {
		fmt.Println("user input ##", line, "##")
		var p = vdu_process_interaction_line(line)
		if !p {

		}
	})

}
