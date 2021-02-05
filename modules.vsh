#!/usr/local/bin/v run

import term

term.clear()

println ("同步vendor")

exec("go mod tidy && go mod vendor") or {
	panic(err)
}

if !is_dir("robotgo") {
    println("克隆robotgo仓库")
    git_result :=  exec("git clone --depth 1 http://github.com/go-vgo/robotgo.git") or {
        panic(err)
    }
    if git_result.exit_code != 0 {
    	eprintln(git_result.output)
    	exit(1)
    }
}

if !is_dir("gohook") {
    println("克隆gohook仓库")
    git_result := exec("git clone --depth 1 http://github.com/robotn/gohook.git") or {
        panic(err)
    }
    if git_result.exit_code != 0 {
        eprintln(git_result.output)
        exit(1)
    }
}

println("拷贝依赖")

system("rm -rf vendor/github.com/go-vgo/robotgo && mkdir vendor/github.com/go-vgo/robotgo")

cp_all("robotgo/", "vendor/github.com/go-vgo/robotgo/", true) or {
  panic("拷贝robotgo失败 ${err}")
}

system("rm -rf vendor/github.com/robotn/gohook && mkdir vendor/github.com/robotn/gohook")

cp_all("gohook", "vendor/github.com/robotn/gohook/", true) or {
  panic(err)
}

println("done!")
