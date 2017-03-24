package main

import (
    "io/ioutil"
    "log"
    "os"
    "fmt"
    "path/filepath"
    "strings"
    "sort"
    "io"
    "os/exec"
    "time"
)

func init() {
    backupsSystemHost()
}

var cmdExec *exec.Cmd

func main() {
    var selected int

    hosts := loadMyHosts("D:/hosts")

    for {
        clear()

        fmt.Println("==============================选择方案===================================")
        fmt.Println()

        fileSlice := make([]string, len(hosts))
        keyList := make([]int, 0)
        for key := range hosts {
            keyList = append(keyList, key)
        }

        sort.Ints(keyList)

        for i := 0; i < len(fileSlice); i++ {
            fileSlice[i] = hosts[i+1]
        }

        fmt.Printf("\t\t\t\t%d:%v\n", 0, "还原系统默认hosts")
        for key, value := range fileSlice {
            value = strings.Split(filepath.Base(value), ".")[0]
            fmt.Printf("\t\t\t\t%d:%v\n", key + 1, value)
        }

        fmt.Println()
        fmt.Println("========================================================================")

        fmt.Scanf("%d", &selected)

        if selected == 0 {
            err := restoreSystemHost()
            checkErr(err, "还原失败")
            flushDns()
            fmt.Println("还原成功")
            time.Sleep(time.Second * 1)
            continue
        }

        switchHostsFile, exist := hosts[selected]
        if !exist {
            fmt.Println(fmt.Errorf("%s", "请选择存在的方案"))
            continue
        }

        err := switchHosts(switchHostsFile)
        checkErr(err, "切换失败")

        fmt.Println("切换成功")
        flushDns()
        time.Sleep(time.Second * 1)
    }
}

func clear() {
    cmdExec = exec.Command("cmd", "/c", "cls")
    cmdExec.Stdout = os.Stdout
    cmdExec.Run()
}

func flushDns() {
    cmdExec = exec.Command("cmd", "/c", "ipconfig/flushdns")
    cmdExec.Run()
}

/**
   获取data目录下各个方案
 */
func loadMyHosts(dir string) (map[int]string) {
    files := make(map[int]string)
    fileInfo, err := ioutil.ReadDir(dir)
    checkErr(err, "找不到方案目录")

    pth := string(os.PathSeparator)
    for key, file := range fileInfo {
        files[key + 1] = dir + pth + file.Name()
    }

    return files
}

/**
  备份原来的系统host
 */
func backupsSystemHost() {
    bakPath := "C:/Windows/System32/drivers/etc/hosts.bak"

    _, err := os.Stat(bakPath)
    if err == nil {
        return
    }

    srcFile, err := os.Open("C:/Windows/System32/drivers/etc/hosts")
    defer srcFile.Close()
    checkErr(err, "无法打开hosts文件")

    desFile, err := os.Create(bakPath)
    defer desFile.Close()
    checkErr(err, "无法创建hosts备份文件")

    _, err = io.Copy(desFile, srcFile)
    checkErr(err, "备份失败")
}

func restoreSystemHost() (error) {
    desFile, err := os.OpenFile("C:/Windows/System32/drivers/etc/hosts.bak", os.O_RDWR, 0755)
    defer desFile.Close()
    checkErr(err, "打开备份文件失败")

    srcFile, err := os.OpenFile("C:/Windows/System32/drivers/etc/hosts", os.O_TRUNC|os.O_WRONLY, 0755)
    checkErr(err, "无法打开hosts文件")

    data, err := ioutil.ReadAll(desFile);
    checkErr(err, "无法读取方案文件")

    _, err = io.WriteString(srcFile, string(data))

    if err != nil {
        return err
    }

    return nil
}

func switchHosts(selected string) (error) {
    desFile, err := os.Open(selected)
    defer desFile.Close()
    checkErr(err, fmt.Sprintf("无法打开[%s]文件", selected))

    srcFile, err := os.OpenFile("C:/Windows/System32/drivers/etc/hosts", os.O_TRUNC|os.O_WRONLY, 0755)
    checkErr(err, "无法打开hosts文件")

    data, err := ioutil.ReadAll(desFile);
    checkErr(err, "无法读取方案文件")

    _, err = io.WriteString(srcFile, string(data))

    if err != nil {
        return err
    }

    return nil
}

func checkErr(err error, msg string) {
    if err != nil {
        fmt.Println("msg:", msg)
        fmt.Print("log:")
        log.Fatal(err)
    }
}