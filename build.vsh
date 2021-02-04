#!/usr/local/bin/v run

import term

term.clear()

if !is_dir('packrd') {
	exec('packr2')?
}

println(term.ok_message('开始打包应用'))

system('qtdeploy build desktop')

mkdir('deploy/darwin/goqt-translate.app/Contents/MacOS/qss')?

cp_all('qss', 'deploy/darwin/goqt-translate.app/Contents/MacOS/qss', true) or {
	println('qss: ${term.fail_message(err)}')
	return
}

cp('goqt-translate.icns', 'deploy/darwin/goqt-translate.app/Contents/Resources/goqt-translate.icns') or {
	println('icns: ${term.fail_message(err)}')
	return
}

println(term.ok_message('构建完成!'))
// exec("qtdeploy run desktop")
system('./deploy/darwin/goqt-translate.app/Contents/MacOS/goqt-translate')
