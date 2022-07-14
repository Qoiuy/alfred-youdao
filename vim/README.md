供 vim 使用

go build -o alfred-youdao


vimrc 文件配置

" 翻译向下n行 并且将翻译结果放在 寄存器 a中
nnoremap <expr> yy ':<C-U>call RunYoudao(' . v:count . ')<CR>' . v:count . 'jo<Esc>'

let mapleader=";"
nmap <Leader>y "ap