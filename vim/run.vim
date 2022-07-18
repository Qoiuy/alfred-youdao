" 打开vim 执行  :call RunYoudao()
" 执行有道翻译 1. 支持参数 2. 翻译结果放在 vim寄存器里
:function RunYoudao(n_line)

	" alfred-youdao shell  环境变量在~/.zshrc里面设置的
	:let ydShell = "~/.vim/codebymyself/alfred-youdao gcd "

	" 清空数据操作
	:call setreg("a", '') " 清空寄存器

	" 取当前行向下n行 的内容
	:let s:curr_line = line(".")
	:let s:body = getline(s:curr_line, s:curr_line+a:n_line)  	"想测试的话  可以使用 :let a:n_line = 1

	:let s:cmd = printf("%s", ydShell)  " 将取到的数据拼成一个 youdao 翻译的shell
	:for s:line in s:body

	    let line = s:line  " 解决 包含(的文本无法翻译的问题
	    let line = substitute(line, '(', '（', 'g')
        let line = substitute(line, ')', '）', 'g')

		:let s:cmd = printf("%s %s ", s:cmd, line)
    :endfor
	":echomsg s:cmd

	" 使用alfred-youdao 插件翻译出来内容 有道参数从本地电脑环境变量里面读取
	" callback 用法详见 :help job-callback
	:let s:ob = job_start(["sh", "-c", s:cmd], {"callback": "HandlerYoudaoResult" })
	":echomsg s:ob

:endfunction

" HandlerYoudaoResult 有道翻译结果的处理
:function HandlerYoudaoResult (channel, msg)

    ":echomsg "n次\n" " golang代码里 print \n 一次 调用改方法一次

    " 将翻译结果写入到寄存器
    :call setreg("a", printf("%s\n", a:msg), "a")

:endfunction

nnoremap <expr> yy ':<C-U>call RunYoudao(' . v:count . ')<CR>' . v:count . 'jo<Esc>'
nmap <Leader>y "ap