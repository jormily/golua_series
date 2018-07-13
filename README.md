# golua_series
##  序列化示例

###  main.go
    func main() {
        L := lua.NewState()
        defer L.Close()

        L.PreloadModule("Luatable", series.Loader)

        if err := L.DoFile("init.lua"); err != nil {
            log.Fatal(err)
            panic(err)
        }

        if err := L.DoString(`
            local _,err,ret = xpcall(main, debug.traceback)
            if err then 
                logd(err) 
            end`); err != nil {
            log.Fatal(err)
            panic(err)
        }

    }
## init.lua

    function main()
        local Luatable = require("Luatable")

        local _test_data = { a = "hello world" ,[1] = 1,[2] = false,[3] = {a = "hellow"} }
        local userData = Luatable.encode(_test_data)
        print(type(userData))
        local t = Luatable.decode(userData)
        logd(t)

    end

## 执行效果
    userdata
    table: 0xc0422a8840 {
    [1] = 1
    [2] = false
    [3] = {
            a = "hellow"
            }
    a = "hello world"
    }
