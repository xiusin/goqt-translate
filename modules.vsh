#!/usr/local/bin/v run

defer {
    rm("modules")
}

println ("同步vendor")

exec("go mod vendor") or {
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
