#!/usr/local/bin/v run

println ("打包")

exec("qtdeploy build desktop") or {
	panic(err)
}

mkdir("deploy/darwin/goqt-translate.app/Contents/MacOS/qss")

cp_all("qss", "deploy/darwin/goqt-translate.app/Contents/MacOS/qss", true) or {
  panic(err)
}

exec("qtdeploy run desktop") or {
	panic(err)
}

println("done!")

