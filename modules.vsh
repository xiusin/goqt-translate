#!/usr/local/bin/v run

defer {
    rm("modules")
}

println ("同步vendor")

exec("go mod vendor") or {
	panic(err)
}

println("拷贝依赖")

rmdir_all("vendor/github.com/go-vgo/robotgo/")
rmdir_all("vendor/github.com/robotn/gohook/")

// 挺low的cp居然不能自动创建目录
mkdir("vendor/github.com/go-vgo/robotgo/")
mkdir("vendor/github.com/robotn/gohook/")

cp_all("robotgo", "vendor/github.com/go-vgo/robotgo", true) or {
  panic(err)
}

cp_all("gohook", "vendor/github.com/robotn/gohook", true) or {
  panic(err)
}

println("done!")
