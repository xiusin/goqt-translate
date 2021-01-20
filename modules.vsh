#!/usr/local/bin/v run

println ("同步vendor")

exec("go mod vendor") or {
	panic(err)
}

println("拷贝依赖")

cp_all("robotgo", "vendor/github.com/go-vgo/", true) or {
  panic(err)
}

cp_all("gohook", "vendor/github.com/robotn/", true) or {
  panic(err)
}

println("done!")



