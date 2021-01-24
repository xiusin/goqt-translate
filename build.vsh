#!/usr/local/bin/v run
import term

term.clear()

defer {
    rm("build")
}

println(term.ok_message('开始打包应用'))

exec("qtdeploy build desktop") or {
	println(term.fail_message(err))
	return
}

mkdir("deploy/darwin/goqt-translate.app/Contents/MacOS/qss")

cp_all("qss", "deploy/darwin/goqt-translate.app/Contents/MacOS/qss", true) or {
  println(term.fail_message(err))
  return
}

cp("goqt-translate.icns", "deploy/darwin/goqt-translate.app/Contents/Resources/goqt-translate.icns") or {
    println(term.fail_message(err))
    return
}

exec("qtdeploy run desktop")

println(term.ok_message("done!"))
